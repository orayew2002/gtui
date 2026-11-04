package term

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

const (
	enterAlt = "\x1b[?1049h"
	leaveAlt = "\x1b[?1049l"
	hideCur  = "\x1b[?25l"
	showCur  = "\x1b[?25h"
	mouseOn  = "\x1b[?1000h\x1b[?1002h\x1b[?1006h"
	mouseOff = "\x1b[?1006l\x1b[?1002l\x1b[?1000l"
	reset    = "\x1b[0m"
	clearAll = "\x1b[2J\x1b[H"
	wrapOff  = "\x1b[?7l"
	wrapOn   = "\x1b[?7h"
)

type Options struct {
	Mouse      bool
	AltScreen  bool
	HideCursor bool
}

func DefaultOptions() Options {
	return Options{Mouse: true, AltScreen: true, HideCursor: true}
}

type Terminal struct {
	In   *os.File
	Out  *os.File
	opts Options

	fd    int
	state *term.State

	resize chan struct{}
	sig    chan os.Signal
	events chan Event
	done   chan struct{}
}

func Open(opts Options) (*Terminal, error) {
	t := &Terminal{
		In:     os.Stdin,
		Out:    os.Stdout,
		opts:   opts,
		fd:     int(os.Stdin.Fd()),
		resize: make(chan struct{}, 1),
		events: make(chan Event, 256),
		done:   make(chan struct{}),
	}

	state, err := term.MakeRaw(t.fd)
	if err != nil {
		return nil, err
	}
	t.state = state

	var setup string
	if opts.AltScreen {
		setup += enterAlt + clearAll
	}
	if opts.HideCursor {
		setup += hideCur
	}
	if opts.Mouse {
		setup += mouseOn
	}
	setup += wrapOff
	t.write(setup)

	t.sig = make(chan os.Signal, 1)
	signal.Notify(t.sig, syscall.SIGWINCH)
	go t.watchResize()
	go t.readInput()

	return t, nil
}

func (t *Terminal) Close() {
	select {
	case <-t.done:
		return
	default:
		close(t.done)
	}

	signal.Stop(t.sig)

	var teardown string
	teardown += wrapOn + reset
	if t.opts.Mouse {
		teardown += mouseOff
	}
	if t.opts.HideCursor {
		teardown += showCur
	}
	if t.opts.AltScreen {
		teardown += leaveAlt
	}
	t.write(teardown)

	if t.state != nil {
		term.Restore(t.fd, t.state)
		t.state = nil
	}
}

func (t *Terminal) Size() (int, int) {
	w, h, err := term.GetSize(int(t.Out.Fd()))
	if err != nil || w <= 0 || h <= 0 {
		return 80, 24
	}
	return w, h
}

func (t *Terminal) Resized() <-chan struct{} { return t.resize }

func (t *Terminal) Events() <-chan Event { return t.events }

func (t *Terminal) Done() <-chan struct{} { return t.done }

func (t *Terminal) write(s string) {
	if s != "" {
		t.Out.WriteString(s)
	}
}

func (t *Terminal) watchResize() {
	for {
		select {
		case <-t.done:
			return
		case <-t.sig:
			select {
			case t.resize <- struct{}{}:
			default:
			}
		}
	}
}

func (t *Terminal) readInput() {
	defer close(t.events)

	buf := make([]byte, 0, 256)
	chunk := make([]byte, 128)

	for {
		n, err := t.In.Read(chunk)
		if n > 0 {
			buf = append(buf, chunk[:n]...)

			for len(buf) > 0 {
				ev, used := decode(buf)
				if used == 0 {
					break
				}
				buf = buf[used:]

				if ev == nil {
					continue
				}
				select {
				case t.events <- ev:
				case <-t.done:
					return
				}
			}
		}

		if err != nil {
			return
		}

		select {
		case <-t.done:
			return
		default:
		}
	}
}
