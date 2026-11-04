// Package lsp is a language server client: enough of the protocol for an
// editor to open a file, keep the server up to date as it is typed into, and
// get diagnostics, completions and hover text back.
//
// It speaks to any server that follows the specification over stdin and
// stdout, which is how gopls, rust-analyzer, pylsp and the rest all run:
//
//	client, err := lsp.Start("/path/to/project", "gopls")
//	client.Open("/path/to/project/main.go", "go", text)
//	client.OnDiagnostics(func(path string, items []lsp.Diagnostic) { ... })
//	client.Change(path, text)          // after every keystroke is fine
//	items, _ := client.Complete(path, line, col)
package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf16"
)

// Diagnostic is one problem the server found, in editor coordinates: lines and
// columns counting from one, the way a person reads them.
type Diagnostic struct {
	Line     int
	Col      int
	EndLine  int
	EndCol   int
	Severity string // error, warning, info or hint
	Message  string
	Source   string
}

// Completion is one suggestion.
type Completion struct {
	Label  string
	Detail string
	Insert string
	Kind   string
}

// Client is a running language server.
type Client struct {
	cmd  *exec.Cmd
	in   io.WriteCloser
	out  *bufio.Reader
	root string

	mu      sync.Mutex
	next    int
	waiting map[int]chan json.RawMessage
	open    map[string]int // version by path
	closed  bool

	onDiagnostics func(path string, items []Diagnostic)
	onLog         func(text string)

	ready chan struct{}
	stop  chan struct{}
}

// Start launches a language server for a project and finishes the handshake, so
// the client it returns is ready to be used.
func Start(root string, command ...string) (*Client, error) {
	if len(command) == 0 {
		return nil, errors.New("lsp: no command to run")
	}

	root, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("lsp: root: %w", err)
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = root

	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("lsp: stdin: %w", err)
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("lsp: stdout: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("lsp: start %s: %w", command[0], err)
	}

	c := &Client{
		cmd:     cmd,
		in:      in,
		out:     bufio.NewReaderSize(out, 1<<16),
		root:    root,
		waiting: map[int]chan json.RawMessage{},
		open:    map[string]int{},
		ready:   make(chan struct{}),
		stop:    make(chan struct{}),
	}

	go c.listen()

	if err := c.handshake(); err != nil {
		c.Stop()
		return nil, err
	}

	return c, nil
}

// OnDiagnostics registers what to do when the server reports problems. It is
// called from the client's own goroutine, so an app should hand the work back
// to its loop, with app.Do for instance.
func (c *Client) OnDiagnostics(fn func(path string, items []Diagnostic)) {
	c.mu.Lock()
	c.onDiagnostics = fn
	c.mu.Unlock()
}

// OnLog registers what to do with the server's own messages, which is where a
// server explains why it is unhappy.
func (c *Client) OnLog(fn func(text string)) {
	c.mu.Lock()
	c.onLog = fn
	c.mu.Unlock()
}

// Open tells the server about a file and its current text.
func (c *Client) Open(path, language, text string) error {
	c.mu.Lock()
	c.open[path] = 1
	c.mu.Unlock()

	return c.notify("textDocument/didOpen", map[string]any{
		"textDocument": map[string]any{
			"uri":        fileURI(path),
			"languageId": language,
			"version":    1,
			"text":       text,
		},
	})
}

// Change sends the whole text again, which every server accepts and which costs
// nothing an editor can feel.
func (c *Client) Change(path, text string) error {
	c.mu.Lock()
	version := c.open[path] + 1
	c.open[path] = version
	c.mu.Unlock()

	return c.notify("textDocument/didChange", map[string]any{
		"textDocument": map[string]any{
			"uri":     fileURI(path),
			"version": version,
		},
		"contentChanges": []any{map[string]any{"text": text}},
	})
}

// Save tells the server the file is on disk now.
func (c *Client) Save(path string) error {
	return c.notify("textDocument/didSave", map[string]any{
		"textDocument": map[string]any{"uri": fileURI(path)},
	})
}

// Complete asks what could go at a position. Line and col count from one.
func (c *Client) Complete(path string, line, col int, text string) ([]Completion, error) {
	raw, err := c.request("textDocument/completion", map[string]any{
		"textDocument": map[string]any{"uri": fileURI(path)},
		"position":     position(text, line, col),
	}, 4*time.Second)
	if err != nil {
		return nil, err
	}

	var list struct {
		Items []completionItem `json:"items"`
	}
	if err := json.Unmarshal(raw, &list); err == nil && len(list.Items) > 0 {
		return completions(list.Items), nil
	}

	var items []completionItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, nil
	}
	return completions(items), nil
}

// Hover asks what is at a position, as plain text.
func (c *Client) Hover(path string, line, col int, text string) (string, error) {
	raw, err := c.request("textDocument/hover", map[string]any{
		"textDocument": map[string]any{"uri": fileURI(path)},
		"position":     position(text, line, col),
	}, 4*time.Second)
	if err != nil {
		return "", err
	}

	var hover struct {
		Contents json.RawMessage `json:"contents"`
	}
	if err := json.Unmarshal(raw, &hover); err != nil {
		return "", nil
	}

	return plainText(hover.Contents), nil
}

// Stop closes the server down.
func (c *Client) Stop() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	c.mu.Unlock()

	close(c.stop)

	_ = c.notify("exit", nil)
	_ = c.in.Close()

	done := make(chan error, 1)
	go func() { done <- c.cmd.Wait() }()

	select {
	case <-done:
	case <-time.After(time.Second):
		_ = c.cmd.Process.Kill()
	}

	return nil
}

func (c *Client) handshake() error {
	_, err := c.request("initialize", map[string]any{
		"processId": nil,
		"rootUri":   fileURI(c.root),
		"capabilities": map[string]any{
			"textDocument": map[string]any{
				"synchronization": map[string]any{"didSave": true},
				"publishDiagnostics": map[string]any{
					"relatedInformation": false,
				},
				"completion": map[string]any{
					"completionItem": map[string]any{"snippetSupport": false},
				},
				"hover": map[string]any{"contentFormat": []string{"plaintext", "markdown"}},
			},
		},
	}, 15*time.Second)
	if err != nil {
		return fmt.Errorf("lsp: initialize: %w", err)
	}

	if err := c.notify("initialized", map[string]any{}); err != nil {
		return err
	}

	close(c.ready)
	return nil
}

// listen reads the server's side of the conversation for as long as it lasts.
func (c *Client) listen() {
	for {
		body, err := readMessage(c.out)
		if err != nil {
			c.failWaiting(err)
			return
		}

		var msg struct {
			ID     *int            `json:"id"`
			Method string          `json:"method"`
			Result json.RawMessage `json:"result"`
			Error  *failure        `json:"error"`
			Params json.RawMessage `json:"params"`
		}
		if json.Unmarshal(body, &msg) != nil {
			continue
		}

		switch {
		case msg.ID != nil && msg.Method == "":
			c.answer(*msg.ID, msg.Result, msg.Error)

		case msg.Method == "textDocument/publishDiagnostics":
			c.publish(msg.Params)

		case msg.Method == "window/logMessage", msg.Method == "window/showMessage":
			c.log(msg.Params)

		case msg.ID != nil:
			// a request from the server; answer nothing rather than hang it
			c.reply(*msg.ID)
		}
	}
}

func (c *Client) answer(id int, result json.RawMessage, bad *failure) {
	c.mu.Lock()
	ch := c.waiting[id]
	delete(c.waiting, id)
	c.mu.Unlock()

	if ch == nil {
		return
	}

	if bad != nil {
		close(ch)
		return
	}

	ch <- result
	close(ch)
}

func (c *Client) publish(params json.RawMessage) {
	var report struct {
		URI         string       `json:"uri"`
		Diagnostics []diagnostic `json:"diagnostics"`
	}
	if json.Unmarshal(params, &report) != nil {
		return
	}

	c.mu.Lock()
	fn := c.onDiagnostics
	c.mu.Unlock()

	if fn == nil {
		return
	}

	out := make([]Diagnostic, 0, len(report.Diagnostics))
	for _, item := range report.Diagnostics {
		out = append(out, Diagnostic{
			Line:     item.Range.Start.Line + 1,
			Col:      item.Range.Start.Character + 1,
			EndLine:  item.Range.End.Line + 1,
			EndCol:   item.Range.End.Character + 1,
			Severity: severity(item.Severity),
			Message:  item.Message,
			Source:   item.Source,
		})
	}

	fn(uriPath(report.URI), out)
}

func (c *Client) log(params json.RawMessage) {
	var message struct {
		Message string `json:"message"`
	}
	if json.Unmarshal(params, &message) != nil {
		return
	}

	c.mu.Lock()
	fn := c.onLog
	c.mu.Unlock()

	if fn != nil {
		fn(message.Message)
	}
}

// reply answers a request from the server with nothing, which is better than
// leaving it waiting.
func (c *Client) reply(id int) {
	_ = c.send(map[string]any{"jsonrpc": "2.0", "id": id, "result": nil})
}

func (c *Client) failWaiting(err error) {
	c.mu.Lock()
	for id, ch := range c.waiting {
		close(ch)
		delete(c.waiting, id)
	}
	c.mu.Unlock()
}

func (c *Client) request(method string, params any, wait time.Duration) (json.RawMessage, error) {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil, errors.New("lsp: client is closed")
	}
	c.next++
	id := c.next
	ch := make(chan json.RawMessage, 1)
	c.waiting[id] = ch
	c.mu.Unlock()

	body := map[string]any{"jsonrpc": "2.0", "id": id, "method": method}
	if params != nil {
		body["params"] = params
	}

	if err := c.send(body); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	select {
	case result, ok := <-ch:
		if !ok {
			return nil, fmt.Errorf("lsp: %s failed", method)
		}
		return result, nil
	case <-ctx.Done():
		c.mu.Lock()
		delete(c.waiting, id)
		c.mu.Unlock()
		return nil, fmt.Errorf("lsp: %s timed out", method)
	}
}

func (c *Client) notify(method string, params any) error {
	body := map[string]any{"jsonrpc": "2.0", "method": method}
	if params != nil {
		body["params"] = params
	}
	return c.send(body)
}

func (c *Client) send(body map[string]any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := fmt.Fprintf(c.in, "Content-Length: %d\r\n\r\n", len(payload)); err != nil {
		return err
	}
	_, err = c.in.Write(payload)
	return err
}

// readMessage reads one framed message: headers, a blank line, then that many
// bytes of JSON.
func readMessage(r *bufio.Reader) ([]byte, error) {
	size := -1

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}

		name, value, ok := strings.Cut(line, ":")
		if !ok || !strings.EqualFold(strings.TrimSpace(name), "content-length") {
			continue
		}

		size, err = strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return nil, fmt.Errorf("lsp: bad content length %q", value)
		}
	}

	if size < 0 {
		return nil, errors.New("lsp: message with no content length")
	}

	body := make([]byte, size)
	if _, err := io.ReadFull(r, body); err != nil {
		return nil, err
	}

	return body, nil
}

// failure is the error shape a response can carry instead of a result.
type failure struct {
	Message string `json:"message"`
}

type diagnostic struct {
	Range struct {
		Start struct {
			Line      int `json:"line"`
			Character int `json:"character"`
		} `json:"start"`
		End struct {
			Line      int `json:"line"`
			Character int `json:"character"`
		} `json:"end"`
	} `json:"range"`
	Severity int    `json:"severity"`
	Message  string `json:"message"`
	Source   string `json:"source"`
}

type completionItem struct {
	Label      string `json:"label"`
	Detail     string `json:"detail"`
	InsertText string `json:"insertText"`
	Kind       int    `json:"kind"`
}

func completions(items []completionItem) []Completion {
	out := make([]Completion, 0, len(items))

	for _, item := range items {
		insert := item.InsertText
		if insert == "" {
			insert = item.Label
		}

		out = append(out, Completion{
			Label:  item.Label,
			Detail: item.Detail,
			Insert: insert,
			Kind:   kind(item.Kind),
		})
	}

	return out
}

func severity(v int) string {
	switch v {
	case 1:
		return "error"
	case 2:
		return "warning"
	case 3:
		return "info"
	case 4:
		return "hint"
	}
	return "error"
}

func kind(v int) string {
	names := map[int]string{
		2: "method", 3: "function", 5: "field", 6: "variable", 7: "class",
		8: "interface", 9: "module", 10: "property", 14: "keyword",
		21: "constant", 22: "struct", 25: "parameter",
	}
	if name, ok := names[v]; ok {
		return name
	}
	return "text"
}

// position turns an editor's line and rune column into the protocol's zero
// based line and UTF-16 character, which is what the specification counts.
func position(text string, line, col int) map[string]any {
	character := col - 1

	if text != "" {
		lines := strings.Split(text, "\n")
		if line >= 1 && line <= len(lines) {
			runes := []rune(lines[line-1])
			if col-1 <= len(runes) {
				character = len(utf16.Encode(runes[:col-1]))
			}
		}
	}

	return map[string]any{"line": line - 1, "character": character}
}

func fileURI(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return "file://" + (&url.URL{Path: abs}).EscapedPath()
}

func uriPath(uri string) string {
	parsed, err := url.Parse(uri)
	if err != nil {
		return strings.TrimPrefix(uri, "file://")
	}
	return parsed.Path
}

// plainText flattens whatever shape a server chose for hover text.
func plainText(raw json.RawMessage) string {
	var text string
	if json.Unmarshal(raw, &text) == nil {
		return text
	}

	var one struct {
		Value string `json:"value"`
	}
	if json.Unmarshal(raw, &one) == nil && one.Value != "" {
		return one.Value
	}

	var many []struct {
		Value string `json:"value"`
	}
	if json.Unmarshal(raw, &many) == nil {
		parts := make([]string, 0, len(many))
		for _, item := range many {
			if item.Value != "" {
				parts = append(parts, item.Value)
			}
		}
		return strings.Join(parts, "\n")
	}

	return ""
}
