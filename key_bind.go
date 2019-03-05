package prompt

type KeyBindResult error

// KeyBindFunc receives an Event and processed it.
type KeyBindFunc func(*Event) KeyBindResult

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
	End: func(e *Event) KeyBindResult {
		end_of_line(e.Buffer())
		return nil
	},
	// Go to the beginning of the line
	Home: func(e *Event) KeyBindResult {
		beginning_of_line(e.Buffer())
		return nil
	},
	// Delete character under/at the cursor
	Delete: func(e *Event) KeyBindResult {
		delete_char(e.Buffer())
		return nil
	},
	// Backspace: delete character before the cursor
	Backspace: func(e *Event) KeyBindResult {
		backward_delete_char(e.Buffer())
		return nil
	},
	// Right arrow: Forward one character
	Right: func(e *Event) KeyBindResult {
		forward_char(e.Buffer())
		return nil
	},
	// Left arrow: Backward one character
	Left: func(e *Event) KeyBindResult {
		backward_char(e.Buffer())
		return nil
	},
}
