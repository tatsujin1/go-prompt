package prompt

import (
	"fmt"
	"os"
)

type Event struct {
	buf *Buffer
	// Invalid if there were no key defined for the control sequence
	key KeyCode
	// empty if this is a key-defined event
	ctrlSeq ControlSequence

	// flags
	translatedKey KeyCode
	endEdit       bool
	eof           bool
	termTitle     *string // nil meaning it's not been set
}

func NewKeyEvent(b *Buffer, k KeyCode) *Event {
	return &Event{
		buf: b,
		key: k,
	}
}

func NewCtrlEvent(b *Buffer, c ControlSequence) *Event {
	return &Event{
		buf:     b,
		ctrlSeq: c,
	}
}

func (e *Event) Buffer() *Buffer {
	return e.buf
}

func (e *Event) Key() KeyCode {
	return e.key
}

func (e *Event) ControlSequence() ControlSequence {
	return e.ctrlSeq
}

func (e *Event) CallFunction(name string, args ...interface{}) {
	// TODO
	fmt.Fprintf(os.Stderr, "CallFunction: '%s' args: %v\n", name, args)
}

func (e *Event) SetEOF() {
	e.eof = true
}

func (e *Event) SetEndEdit() {
	e.endEdit = true
}

func (e *Event) SetTranslatedKey(key KeyCode) {
	e.translatedKey = key
}

func (e *Event) SetTitle(title string) {
	e.termTitle = &title
}
