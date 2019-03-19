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
	KeyEnd: end_of_line,
	// Go to the beginning of the line
	KeyHome: beginning_of_line,
	// Delete character under/at the cursor
	KeyDelete: delete_char,
	// Backspace: delete character before the cursor
	KeyBackspace: backward_delete_char,
	// Right arrow: Forward one character
	KeyRight: forward_char,
	// Left arrow: Backward one character
	KeyLeft: backward_char,
}
