package prompt

// KeyBindFunc receives buffer and processed it.
type KeyBindFunc func(*Buffer)

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
type KeyBindMode string

const (
	// CommonKeyBind is a mode without any keyboard shortcut
	CommonKeyBind KeyBindMode = "common"
	// EmacsKeyBind is a mode to use emacs-like keyboard shortcut
	EmacsKeyBind KeyBindMode = "emacs"
)

var commonKeyBindings = map[KeyCode]KeyBindFunc{
	// Go to the End of the line
	End: GoLineEnd,
	// Go to the beginning of the line
	Home: GoLineBeginning,
	// Delete character under/at the cursor
	Delete: DeleteChar,
	// Backspace: delete character before the cursor
	Backspace: DeleteBeforeChar,
	// Right arrow: Forward one character
	Right: GoRightChar,
	// Left arrow: Backward one character
	Left: GoLeftChar,
}
