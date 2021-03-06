package prompt

import (
	"fmt"
	"os"
	rdebug "runtime/debug"
	"time"

	"github.com/tatsujin/go-prompt/internal/debug"
)

// Executor is called when user input something text.
type Executor func(string)

// Completer should return the suggest item from Document.
type Completer func(Document) []Choice

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
	debug.Log("start prompt")
	p.setUp()

	if p.completion.showAtStart && p.completion.asYouType {
		p.completion.FindCompletions(*p.buf.Document())
	}
	p.renderer.Render(p.buf, p.completion)

	bufCh := make(chan ControlSequence, 128)
	stopReadBufCh := make(chan struct{}, 1) // buffered so the deferred tear down doesn't block
	go p.readBuffer(bufCh, stopReadBufCh)

	exitCh := make(chan int)
	termSizeCh := make(chan *WinSize)
	stopHandleSignalCh := make(chan struct{}, 1) // buffered so the deferred tear down doesn't block
	go p.handleSignals(exitCh, termSizeCh, stopHandleSignalCh)

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, r)
			rdebug.PrintStack()
		}

		stopReadBufCh <- struct{}{}
		stopHandleSignalCh <- struct{}{}
		p.tearDown()
		debug.Teardown()
	}()

	for {
		select {
		case cs := <-bufCh:
			if shouldExit, exec := p.feed(cs); shouldExit {
				fmt.Fprintln(os.Stderr, "EXIT")
				p.renderer.BreakLine(p.buf, true)
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}
				return
			} else if exec != nil {
				fmt.Fprintln(os.Stderr, "EXECUTE")

				// execute entered command-line

				// Stop goroutine to run readBuffer function
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}

				// Unset raw mode
				// Reset to Blocking mode because returned EAGAIN when still set non-blocking mode.
				debug.AssertNoError(p.in.TearDown())
				p.executor(exec.input)

				if p.completion.showAtStart && p.completion.asYouType {
					p.completion.FindCompletions(*p.buf.Document())
				}
				p.renderer.Render(p.buf, p.completion)

				// Set raw mode
				debug.AssertNoError(p.in.Setup())
				go p.readBuffer(bufCh, stopReadBufCh)
				go p.handleSignals(exitCh, termSizeCh, stopHandleSignalCh)
			} else {
				if p.completion.asYouType {
					p.completion.FindCompletions(*p.buf.Document())
				}
				p.renderer.Render(p.buf, p.completion)
			}
		case w := <-termSizeCh:
			p.renderer.UpdateWinSize(w)
			p.renderer.Render(p.buf, p.completion)
		case code := <-exitCh:
			p.renderer.BreakLine(p.buf, true)
			p.tearDown()
			return code
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	return 0
}

func (p *Prompt) feed(cs ControlSequence) (shouldExit bool, exec *Exec) {
	key := FindKey(cs)

	fmt.Fprintf(os.Stderr, "--> key: %v\n", []byte(cs))

	p.buf.flags.translatedKey = Undefined

	// are we already selecting a completion suggestion?
	completing := p.completion.Completing()
	key = p.handleCompletionKeyBinding(key, completing)
	if p.completion.Completing() {
		fmt.Fprintln(os.Stderr, "\x1b[36mcompleting\x1b[m")
	}

	p.handleKeyBinding(key)

	if p.buf.flags.eof {
		shouldExit = true
		return
	} else if p.buf.flags.endEdit {
		key = KeyEnter
	} else if tkey := p.buf.flags.translatedKey; tkey != Undefined {
		if tkey == Ignore {
			return
		}
		key = tkey
	}

	switch key {
	case KeyEnter, KeyControl | KeyJ, KeyControl | KeyM:
		p.renderer.BreakLine(p.buf, false)

		exec = &Exec{input: p.buf.Text()}

		p.buf = NewBuffer()

		// TODO: should add to history after the command has been successfully parsed!
		if len(exec.input) > 0 {
			p.history.Add(exec.input)
		}
	case KeyControl | KeyC:
		p.renderer.BreakLine(p.buf, true)
		p.buf = NewBuffer()
		p.history.ClearModified()
	case KeyUp, KeyControl | KeyP:
		if !completing { // Don't use p.completion.Completing() because it takes double operation when switch to selected=-1.
			// if current edit is multi-line (and we're not on the first line), go up one line
			doc := p.buf.Document()
			if doc.CursorRow() > 0 {
				p.buf.CursorUp(1)
			} else {
				p.buf = p.history.Previous(p.buf)
			}
		}
	case KeyDown, KeyControl | KeyN:
		if !completing { // Don't use p.completion.Completing() because it takes double operation when switch to selected=-1.
			// if current edit is multi-line (and we're not on the last line), go down one line
			doc := p.buf.Document()
			if !doc.CursorOnLastLine() && doc.LineCount() > 1 {
				fmt.Fprintln(os.Stderr, "line down")
				p.buf.CursorDown(1)
			} else {
				p.buf = p.history.Next(p.buf)
			}
			return
		}
	case Undefined:
		if !p.handleControlSequenceBinding(cs) {
			p.buf.InsertText(string(cs), false, true)
		}
	}

	if p.buf.flags.eof {
		shouldExit = true
	}

	return
}

func (p *Prompt) handleCompletionKeyBinding(key KeyCode, completing bool) KeyCode {
	switch key {
	case KeyDown:
		if completing { // only if already completing
			p.completion.Next()
		}
	case KeyTab, KeyControl | KeyI: // next choice, or start completing
		if !p.completion.asYouType {
			p.completion.FindCompletions(*p.buf.Document())
		}
		p.completion.Next()
	case KeyUp:
		if completing { // only if already completing
			p.completion.Previous()
		}
	case KeyBackTab: // previous choice, or start completing
		p.completion.Previous()
	default:
		if s, ok := p.completion.Selected(); ok {
			w := p.buf.Document().GetWordBeforeCursorUntilSeparator(p.completion.wordSeparator)
			if w != "" {
				p.buf.DeleteBeforeCursor(Offset(len([]rune(w))))
			}
			p.buf.InsertText(s.Text, false, true)

			// if completion was accepted using Enter, that key shouldn't be handled when we return
			if key == KeyEnter {
				key = Ignore
			}
		}

		p.completion.Reset()
	}
	return key
}

func (p *Prompt) handleKeyBinding(key KeyCode) bool {
	ev := NewKeyEvent(p.buf, key)
	// TODO: expose an API for the handlers:
	//   the handler can then do e.g.:
	//   ev.CallFunction("delete-char-backwards", args...)

	handled := false

	// Custom key bindings
	if fn, ok := p.keyBindings[key]; ok {
		//fmt.Fprintf(os.Stderr, "executing custom key bind\n")
		fn(ev)
		handled = true
	}

	// "generic" key bindings
	if fn, ok := commonKeyBindings[key]; ok {
		///fmt.Fprintf(os.Stderr, "executing common key bind\n")
		fn(ev)
		handled = true
	}

	// mode-specific key bindings
	if p.editMode == EmacsMode {
		if fn, ok := emacsKeyBindings[key]; ok {
			//fmt.Fprintf(os.Stderr, "executing emacs key bind\n")
			fn(ev)
			handled = true
		}
	}

	if handled {
		p.postEventHandling(ev)
	}

	return handled
}

func (p *Prompt) handleControlSequenceBinding(cs ControlSequence) bool {
	if fn, ok := p.ControlSequenceBindings[cs]; ok {
		ev := NewCtrlEvent(p.buf, cs)
		fn(ev)
		p.postEventHandling(ev)
		return true
	}
	return false
}

func (p *Prompt) postEventHandling(ev *Event) {
	if ev.endEdit {
		p.buf.setEndEdit()
		ev.endEdit = false
	}
	if ev.eof {
		p.buf.setEOF()
		ev.eof = false
	}
	if ev.translatedKey != Undefined {
		p.buf.setTranslatedKey(ev.translatedKey)
		ev.translatedKey = Undefined
	}
	if ev.termTitle != nil {
		p.renderer.out.SetTitle(*ev.termTitle)
		ev.termTitle = nil
	}
}

// Input just returns user input text.
func (p *Prompt) Input() string {
	defer debug.Teardown()
	debug.Log("start prompt")
	p.setUp()

	if p.completion.showAtStart {
		p.completion.FindCompletions(*p.buf.Document())
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
				p.renderer.BreakLine(p.buf, true)
				return ""
			} else if e != nil {
				// Stop goroutine to run readBuffer function
				return e.input
			} else {
				p.completion.FindCompletions(*p.buf.Document())
				p.renderer.Render(p.buf, p.completion)
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (p *Prompt) OutputAsync(format string, a ...interface{}) {
	p.renderer.OutputAsync(p.buf, p.completion, format, a...)
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
