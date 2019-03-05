package prompt

// KeyBindFunc receives an Event and processed it.
type KeyBindFunc func(*Event)

// KeyBind represents which key should do what operation.
type KeyBind struct {
	Key KeyCode
	Fn  KeyBindFunc
}

// ControlBind binds a specific control sequence to an operation.
type ControlSequenceBind struct {
	Sequence ControlSequence
	Fn       KeyBindFunc
}

// KeyBindMode to switch a key binding flexibly.
type EditMode string

const (
	// SimpleMode is a mode without any keyboard shortcuts
	SimpleMode EditMode = "common"
	// EmacsKeyBind is a mode to use emacs-like keyboard shortcuts
	EmacsMode EditMode = "emacs"
	// TODO: vi?
)

var commonKeyBindings = map[KeyCode]KeyBindFunc{
	// Go to the End of the line
	End: func(e *Event) { end_of_line(e.Buffer()) },
	// Go to the beginning of the line
	Home: func(e *Event) { beginning_of_line(e.Buffer()) },
	// Delete character under/at the cursor
	Delete: func(e *Event) { delete_char(e.Buffer()) },
	// Backspace: delete character before the cursor
	Backspace: func(e *Event) { backward_delete_char(e.Buffer()) },
	// Right arrow: Forward one character
	Right: func(e *Event) { forward_char(e.Buffer()) },
	// Left arrow: Backward one character
	Left: func(e *Event) { backward_char(e.Buffer()) },
}
