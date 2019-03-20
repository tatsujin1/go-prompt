package prompt

// WinSize represents the width and height of terminal.
type WinSize struct {
	Row uint16
	Col uint16
}

// ConsoleParser is an interface to abstract input layer.
type ConsoleParser interface {
	// Setup should be called before starting input
	Setup() error
	// TearDown should be called after stopping input
	TearDown() error
	// GetWinSize returns WinSize object to represent width and height of terminal.
	GetWinSize() *WinSize
	// Read returns byte array.
	Read() ([]byte, error)
}

// FindKey returns Key correspond to input byte codes, or Undefined if no key is defined.
func FindKey(cs ControlSequence) KeyCode {
	if key, ok := KeySequences[cs]; ok {
		return key
	}
	return Undefined
}

// HasControlModifier returns whether the given key code has the Control modifier.
func HasControlModifier(k KeyCode) bool {
	return k&KeyControl > 0
}

// HasAltModifier returns whether the given key code has the Alt modifier.
func HasAltModifier(k KeyCode) bool {
	return k&KeyAlt > 0
}

// KeySequences holds mappings of control sequence to a logical key code.
var KeySequences = map[ControlSequence]KeyCode{
	"\x1b": KeyEscape,

	"\x00": KeyControl | KeySpace,
	"\x01": KeyControl | KeyA,
	"\x02": KeyControl | KeyB,
	"\x03": KeyControl | KeyC,
	"\x04": KeyControl | KeyD,
	"\x05": KeyControl | KeyE,
	"\x06": KeyControl | KeyF,
	"\x07": KeyControl | KeyG,
	"\x08": KeyControl | KeyH, // KeyBackspace
	//"\x09": KeyControl | KeyI,
	//"\x0a": KeyControl | KeyJ,
	"\x0b": KeyControl | KeyK,
	"\x0c": KeyControl | KeyL,
	"\x0d": KeyEnter, // KeyControl | KeyM
	"\x0e": KeyControl | KeyN,
	"\x0f": KeyControl | KeyO,
	"\x10": KeyControl | KeyP,
	"\x11": KeyControl | KeyQ,
	"\x12": KeyControl | KeyR,
	"\x13": KeyControl | KeyS,
	"\x14": KeyControl | KeyT,
	"\x15": KeyControl | KeyU,
	"\x16": KeyControl | KeyV,
	"\x17": KeyControl | KeyW,
	"\x18": KeyControl | KeyX,
	"\x19": KeyControl | KeyY,
	"\x1a": KeyControl | KeyZ,

	"\x1c": KeyControl | KeyBackslash,
	"\x1d": KeyControl | KeySquareClose,
	"\x1e": KeyControl | KeyCircumflex,
	"\x1f": KeyControl | KeyUnderscore,
	"\x7f": KeyControl | KeyBackspace,

	"\x1b[A": KeyUp,
	"\x1b[B": KeyDown,
	"\x1b[C": KeyRight,
	"\x1b[D": KeyLeft,
	"\x1b[H": KeyHome,
	"\x1b[F": KeyEnd,
	"\x1b0H": KeyHome,
	"\x1b0F": KeyEnd,

	"\x0a":      KeyEnter,
	"\x1b[3;2~": KeyShift | KeyDelete,
	"\x1b[3;5~": KeyControl | KeyDelete,
	"\x1b[1~":   KeyHome,
	"\x1b[2~":   KeyInsert,
	"\x1b[3~":   KeyDelete,
	"\x1b[4~":   KeyEnd,
	"\x1b[5~":   KeyPageUp,
	"\x1b[6~":   KeyPageDown,
	"\x1b[7~":   KeyHome,
	"\x1b[8~":   KeyEnd,
	"\x09":      KeyTab,
	"\x1b[Z":    KeyBackTab,

	"\x1bOP": KeyF1,
	"\x1bOQ": KeyF2,
	"\x1bOR": KeyF3,
	"\x1bOS": KeyF4,

	"\x1bOPA": KeyF1, // Linux console
	"\x1b[[B": KeyF2, // Linux console
	"\x1b[[C": KeyF3, // Linux console
	"\x1b[[D": KeyF4, // Linux console
	"\x1b[[E": KeyF5, // Linux console

	"\x1b[\x11~": KeyF1, // rxvt-unicode
	"\x1b[\x12~": KeyF2, // rxvt-unicode
	"\x1b[\x13~": KeyF3, // rxvt-unicode
	"\x1b[\x14~": KeyF4, // rxvt-unicode

	"\x1b[15~":     KeyF5,
	"\x1b[17~":     KeyF6,
	"\x1b[18~":     KeyF7,
	"\x1b[19~":     KeyF8,
	"\x1b[20~":     KeyF9,
	"\x1b[21~":     KeyF10,
	"\x1b[22~":     KeyF11,
	"\x1b[24~\x08": KeyF12,
	"\x1b[\x25~":   KeyF13,
	"\x1b[\x26~":   KeyF14,
	"\x1b[\x28~":   KeyF15,
	"\x1b[\x29~":   KeyF16,
	// conflict `Home`: "\x1b[1~": Key     F17,
	// conflict `Insert`: "\x1b[2~": Key     F18,
	// conflict `Delete`: "\x1b[3~": Key     F19,
	// conflict `End`: "\x1b[4~": Key     F20,

	// Xterm
	"\x1b[1;2P": KeyF13,
	"\x1b[1;2Q": KeyF14,
	// &ASCIICode"\x1b[1;2\x52": KeyF15,  // Conflicts with CPR response
	"\x1b[1;2R":    KeyF16,
	"\x1b[\x15;2~": KeyF17,
	"\x1b[\x17;2~": KeyF18,
	"\x1b[\x18;2~": KeyF19,
	"\x1b[\x19;2~": KeyF20,
	"\x1b[\x20;2~": KeyF21,
	"\x1b[\x21;2~": KeyF22,
	"\x1b[\x23;2~": KeyF23,
	"\x1b[\x24;2~": KeyF24,

	"\x1b[1;5A": KeyControl | KeyUp,
	"\x1b[1;5B": KeyControl | KeyDown,
	"\x1b[1;5C": KeyControl | KeyRight,
	"\x1b[1;5D": KeyControl | KeyLeft,

	"\x1b[1;2A": KeyShift | KeyUp,
	"\x1b[1;2B": KeyShift | KeyDown,
	"\x1b[1;2C": KeyShift | KeyRight,
	"\x1b[1;2D": KeyShift | KeyLeft,

	// Tmux sends following keystrokes when control+arrow is pressed, but for
	// Emacs ansi-term sends the same sequences for normal arrow keys. Consider
	// it a normal arrow press, because that's more important.
	"\x1b0A": KeyUp,
	"\x1b0B": KeyDown,
	"\x1b0C": KeyRight,
	"\x1b0D": KeyLeft,

	"\x1b[5A": KeyControl | KeyUp,
	"\x1b[5B": KeyControl | KeyDown,
	"\x1b[5C": KeyControl | KeyRight,
	"\x1b[5D": KeyControl | KeyLeft,

	"\x1b[Oc": KeyControl | KeyRight, // rxvt
	"\x1b[Od": KeyControl | KeyLeft,  // rxvt

	"\x1b[E": Ignore, // Xterm
	// conflict with 'End':  "\x1b[F": KeyIgnore, // Linux console

	"\x1ba": KeyAlt | KeyA,
	"\x1bb": KeyAlt | KeyB,
	"\x1bc": KeyAlt | KeyC,
	"\x1bd": KeyAlt | KeyD,
	"\x1be": KeyAlt | KeyE,
	"\x1bf": KeyAlt | KeyF,
	"\x1bg": KeyAlt | KeyG,
	"\x1bh": KeyAlt | KeyH,
	"\x1bi": KeyAlt | KeyI,
	"\x1bj": KeyAlt | KeyJ,
	"\x1bk": KeyAlt | KeyK,
	"\x1bl": KeyAlt | KeyL,
	"\x1bm": KeyAlt | KeyM,
	"\x1bn": KeyAlt | KeyN,
	"\x1bo": KeyAlt | KeyO,
	"\x1bp": KeyAlt | KeyP,
	"\x1bq": KeyAlt | KeyQ,
	"\x1br": KeyAlt | KeyR,
	"\x1bs": KeyAlt | KeyS,
	"\x1bt": KeyAlt | KeyT,
	"\x1bu": KeyAlt | KeyU,
	"\x1bv": KeyAlt | KeyV,
	"\x1bw": KeyAlt | KeyW,
	"\x1bx": KeyAlt | KeyX,
	"\x1by": KeyAlt | KeyY,
	"\x1bz": KeyAlt | KeyZ,

	"\x1b\x08":  KeyAlt | KeyBackspace,
	"\x1b\x13":  KeyAlt | KeyEnter,
	"\x1b[3;3~": KeyAlt | KeyDelete,
}
