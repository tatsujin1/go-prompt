package prompt

type KeyBindResult error

// KeyBindFunc receives buffer and processed it.
type KeyBindFunc func(*Buffer) KeyBindResult

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
	End: func(b *Buffer) KeyBindResult { GoLineEnd(b); return nil },
	// Go to the beginning of the line
	Home: func(b *Buffer) KeyBindResult { GoLineBeginning(b); return nil },
	// Delete character under/at the cursor
	Delete: func(b *Buffer) KeyBindResult { DeleteChar(b); return nil },
	// Backspace: delete character before the cursor
	Backspace: func(b *Buffer) KeyBindResult { DeleteBeforeChar(b); return nil },
	// Right arrow: Forward one character
	Right: func(b *Buffer) KeyBindResult { GoRightChar(b); return nil },
	// Left arrow: Backward one character
	Left: func(b *Buffer) KeyBindResult { GoLeftChar(b); return nil },
}
