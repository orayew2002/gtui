package gtui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var langs = struct {
	sync.RWMutex
	table    map[string]map[string]string
	current  string
	fallback string
}{table: map[string]map[string]string{}}

func AddLanguage(code string, phrases map[string]string) {
	langs.Lock()
	defer langs.Unlock()

	existing := langs.table[code]
	if existing == nil {
		existing = make(map[string]string, len(phrases))
		langs.table[code] = existing
	}

	for key, value := range phrases {
		existing[key] = value
	}

	if langs.current == "" {
		langs.current = code
	}
	if langs.fallback == "" {
		langs.fallback = code
	}
}

func LoadLanguages(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".json" {
			continue
		}

		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}

		var tree map[string]any
		if err := json.Unmarshal(raw, &tree); err != nil {
			return fmt.Errorf("gtui: %s: %w", name, err)
		}

		phrases := map[string]string{}
		flatten("", tree, phrases)
		AddLanguage(strings.TrimSuffix(name, ".json"), phrases)
	}

	return nil
}

func flatten(prefix string, tree map[string]any, out map[string]string) {
	for key, value := range tree {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		switch v := value.(type) {
		case string:
			out[path] = v
		case map[string]any:
			flatten(path, v, out)
		default:
			out[path] = fmt.Sprint(v)
		}
	}
}

func SetLanguage(code string) {
	langs.Lock()
	langs.current = code
	langs.Unlock()
}

func Language() string {
	langs.RLock()
	defer langs.RUnlock()
	return langs.current
}

func Languages() []string {
	langs.RLock()
	defer langs.RUnlock()

	out := make([]string, 0, len(langs.table))
	for code := range langs.table {
		out = append(out, code)
	}

	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}

	return out
}

func SetFallback(code string) {
	langs.Lock()
	langs.fallback = code
	langs.Unlock()
}

func T(key string, args ...any) string {
	langs.RLock()
	phrase, ok := langs.table[langs.current][key]
	if !ok {
		phrase, ok = langs.table[langs.fallback][key]
	}
	langs.RUnlock()

	if !ok {
		phrase = key
	}
	if len(args) == 0 {
		return phrase
	}

	return fmt.Sprintf(phrase, args...)
}

func (a *App) View(build func() Component) *App {
	a.view = build
	return a.Refresh()
}

func (a *App) Refresh() *App {
	if a.view == nil {
		return a
	}

	focused := ""
	if a.focus != nil {
		focused = a.focus.ID
	}

	a.HTML(a.view())

	if focused != "" {
		if node := a.Query("#" + focused); node != nil {
			a.Focus(node)
		}
	}

	return a.Redraw()
}

func (a *App) Language(code string) *App {
	SetLanguage(code)
	return a.Refresh()
}
