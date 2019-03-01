package prompt

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/c-bata/go-prompt/internal/debug"
)

// Executor is called when user input something text.
type Executor func(string)

// Completer should return the suggest item from Document.
type Completer func(Document) []Suggest

// Prompt is core struct of go-prompt.
type Prompt struct {
	in                      ConsoleParser
	buf                     *Buffer
	renderer                *Render
	executor                Executor
	history                 *History
	completion              *CompletionManager
	keyBindings             map[KeyCode]KeyBindFunc
	ControlSequenceBindings map[ControlSequence]KeyBindFunc
	editMode                EditMode
}

// Exec is the struct contains user input context.
type Exec struct {
	input string
}

// Run starts prompt.
func (p *Prompt) Run() (exitCode int) {
	defer debug.Teardown()
	debug.Log("start prompt")
	p.setUp()

	if p.completion.showAtStart {
		p.completion.Update(*p.buf.Document())
	}

	p.renderer.Render(p.buf, p.completion)

	bufCh := make(chan ControlSequence, 128)
	stopReadBufCh := make(chan struct{}, 1) // buffered so the deferred tear down doesn't block
	go p.readBuffer(bufCh, stopReadBufCh)

	exitCh := make(chan int)
	winSizeCh := make(chan *WinSize)
	stopHandleSignalCh := make(chan struct{}, 1) // buffered so the deferred tear down doesn't block
	go p.handleSignals(exitCh, winSizeCh, stopHandleSignalCh)

	defer func() {
		stopReadBufCh <- struct{}{}
		stopHandleSignalCh <- struct{}{}
		p.tearDown()
	}()

	for {
		select {
		case cs := <-bufCh:
			if shouldExit, e := p.feed(cs); shouldExit {
				p.renderer.BreakLine(p.buf)
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}
				return
			} else if e != nil {
				// Stop goroutine to run readBuffer function
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}

				// Unset raw mode
				// Reset to Blocking mode because returned EAGAIN when still set non-blocking mode.
				debug.AssertNoError(p.in.TearDown())
				p.executor(e.input)

				p.completion.Update(*p.buf.Document())
				p.renderer.Render(p.buf, p.completion)

				// Set raw mode
				debug.AssertNoError(p.in.Setup())
				go p.readBuffer(bufCh, stopReadBufCh)
				go p.handleSignals(exitCh, winSizeCh, stopHandleSignalCh)
			} else {
				p.completion.Update(*p.buf.Document())
				p.renderer.Render(p.buf, p.completion)
			}
		case w := <-winSizeCh:
			p.renderer.UpdateWinSize(w)
			p.renderer.Render(p.buf, p.completion)
		case code := <-exitCh:
			p.renderer.BreakLine(p.buf)
			p.tearDown()
			return code
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	return 0
}

func (p *Prompt) feed(cs ControlSequence) (shouldExit bool, exec *Exec) {
	key := GetKey(cs)

	fmt.Fprintf(os.Stderr, "got key: %v\n", []byte(cs))

	p.buf.flags.translated_key = Undefined

	// completion
	completing := p.completion.Completing()
	p.handleCompletionKeyBinding(key, completing)

	if bind_res, ok := p.handleKeyBinding(key); ok && bind_res != nil {
		p.presentError(bind_res)
	}

	if p.buf.flags.eof {
		shouldExit = true
		return
	}
	if p.buf.flags.translated_key != Undefined {
		if p.buf.flags.translated_key == Ignore {
			return
		}
		key = p.buf.flags.translated_key
	}

	switch key {
	case Enter, ControlJ, ControlM:
		p.renderer.BreakLine(p.buf)

		exec = &Exec{input: p.buf.Text()}
		p.buf = NewBuffer()
		if exec.input != "" {
			p.history.Add(exec.input)
		}
	case ControlC:
		p.renderer.BreakLine(p.buf)
		p.buf = NewBuffer()
		p.history.Clear()
	case Up, ControlP:
		if !completing { // Don't use p.completion.Completing() because it takes double operation when switch to selected=-1.
			// if current edit is multi-line (and we're not on the first line), go up one line
			if p.buf.Document().CursorPositionRow() > 0 {
				fmt.Fprintln(os.Stderr, "line up")
				p.buf.CursorUp(1)
			} else if newBuf, changed := p.history.Older(p.buf); changed {
				p.buf = newBuf
			}
		}
	case Down, ControlN:
		if !completing { // Don't use p.completion.Completing() because it takes double operation when switch to selected=-1.
			// if current edit is multi-line (and we're not on the last line), go down one line
			if p.buf.Document().CursorPositionRow()+1 < p.buf.Document().LineCount() {
				fmt.Fprintln(os.Stderr, "line down")
				p.buf.CursorDown(1)
			} else if newBuf, changed := p.history.Newer(p.buf); changed {
				p.buf = newBuf
			}
			return
		}
	case Undefined:
		if bind_res, ok := p.handleControlSequenceBinding(cs); ok && bind_res != nil {
			if bind_res == io.EOF {
				shouldExit = true
				return
			}
			p.presentError(bind_res)
		}
		p.buf.InsertText(string(cs), false, true)
	}

	if p.buf.flags.eof {
		shouldExit = true
	}

	return
}

func (p *Prompt) presentError(err error) {
	fmt.Fprintf(os.Stderr, "\x1b[31mError:\x1b[m [%T] %s\n", err, err)
}

func (p *Prompt) handleCompletionKeyBinding(key KeyCode, completing bool) {
	switch key {
	case Down:
		if completing { // only if already completing
			p.completion.Next()
		}
	case Tab, ControlI: // next suggestion, or start completing
		p.completion.Next()
	case Up:
		if completing { // only if already completing
			p.completion.Previous()
		}
	case BackTab: // previous suggestion, or start completing
		p.completion.Previous()
	default:
		if s, ok := p.completion.GetSelectedSuggestion(); ok {
			w := p.buf.Document().GetWordBeforeCursorUntilSeparator(p.completion.wordSeparator)
			if w != "" {
				p.buf.DeleteBeforeCursor(len([]rune(w)))
			}
			p.buf.InsertText(s.Text, false, true)
		}
		p.completion.Reset()
	}
}

func (p *Prompt) handleKeyBinding(key KeyCode) (KeyBindResult, bool) {
	// Custom key bindings
	if fn, ok := p.keyBindings[key]; ok {
		fmt.Fprintf(os.Stderr, "executing custom key bind\n")
		return fn(p.buf), true
	}

	// "generic" key bindings
	if fn, ok := commonKeyBindings[key]; ok {
		fmt.Fprintf(os.Stderr, "executing common key bind\n")
		return fn(p.buf), true
	}

	// mode-specific key bindings
	if p.editMode == EmacsMode {
		if fn, ok := emacsKeyBindings[key]; ok {
			fmt.Fprintf(os.Stderr, "executing emacs key bind\n")
			return fn(p.buf), true
		}
	}

	return nil, false
}

func (p *Prompt) handleControlSequenceBinding(cs ControlSequence) (KeyBindResult, bool) {
	if fn, ok := p.ControlSequenceBindings[cs]; ok {
		return fn(p.buf), true
	}
	return nil, false
}

// Input just returns user input text.
func (p *Prompt) Input() string {
	defer debug.Teardown()
	debug.Log("start prompt")
	p.setUp()

	if p.completion.showAtStart {
		p.completion.Update(*p.buf.Document())
	}

	p.renderer.Render(p.buf, p.completion)
	bufCh := make(chan ControlSequence, 128)
	stopReadBufCh := make(chan struct{})
	go p.readBuffer(bufCh, stopReadBufCh)

	defer func() {
		stopReadBufCh <- struct{}{}
		p.tearDown()
	}()

	for {
		select {
		case b := <-bufCh:
			if shouldExit, e := p.feed(b); shouldExit {
				p.renderer.BreakLine(p.buf)
				return ""
			} else if e != nil {
				// Stop goroutine to run readBuffer function
				return e.input
			} else {
				p.completion.Update(*p.buf.Document())
				p.renderer.Render(p.buf, p.completion)
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (p *Prompt) readBuffer(bufCh chan ControlSequence, stopCh chan struct{}) {
	debug.Log("start reading buffer")
	for {
		select {
		case <-stopCh:
			debug.Log("stop reading buffer")
			return
		default:
			if b, err := p.in.Read(); err == nil && !(len(b) == 1 && b[0] == 0) {
				bufCh <- ControlSequence(b)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (p *Prompt) setUp() {
	debug.AssertNoError(p.in.Setup())
	p.renderer.Setup()
	p.renderer.UpdateWinSize(p.in.GetWinSize())
}

func (p *Prompt) tearDown() {
	debug.AssertNoError(p.in.TearDown())
	p.renderer.TearDown()
}
