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

// KeySequences holds mappings of control sequence to a logical key code.
var KeySequences = map[ControlSequence]KeyCode{
	"\x1b": Escape,

	"\x00": Control | Space,
	"\x01": Control | A,
	"\x02": Control | B,
	"\x03": Control | C,
	"\x04": Control | D,
	"\x05": Control | E,
	"\x06": Control | F,
	"\x07": Control | G,
	"\x08": Control | H, // Backspace
	//"\x09": Control | I,
	//"\x0a": Control | J,
	"\x0b": Control | K,
	"\x0c": Control | L,
	"\x0d": Enter, // Control | M
	"\x0e": Control | N,
	"\x0f": Control | O,
	"\x10": Control | P,
	"\x11": Control | Q,
	"\x12": Control | R,
	"\x13": Control | S,
	"\x14": Control | T,
	"\x15": Control | U,
	"\x16": Control | V,
	"\x17": Control | W,
	"\x18": Control | X,
	"\x19": Control | Y,
	"\x1a": Control | Z,

	"\x1c": Control | Backslash,
	"\x1d": Control | SquareClose,
	"\x1e": Control | Circumflex,
	"\x1f": Control | Underscore,
	"\x7f": Control | Backspace,

	"\x1b[A": Up,
	"\x1b[B": Down,
	"\x1b[C": Right,
	"\x1b[D": Left,
	"\x1b[H": Home,
	"\x1b[F": End,
	"\x1b0H": Home,
	"\x1b0F": End,

	"\x0a":      Enter,
	"\x1b[3;2~": Shift | Delete,
	"\x1b[3;5~": Control | Delete,
	"\x1b[1~":   Home,
	"\x1b[2~":   Insert,
	"\x1b[3~":   Delete,
	"\x1b[4~":   End,
	"\x1b[5~":   PageUp,
	"\x1b[6~":   PageDown,
	"\x1b[7~":   Home,
	"\x1b[8~":   End,
	"\x09":      Tab,
	"\x1b[Z":    BackTab,

	"\x1bOP": F1,
	"\x1bOQ": F2,
	"\x1bOR": F3,
	"\x1bOS": F4,

	"\x1bOPA": F1, // Linux console
	"\x1b[[B": F2, // Linux console
	"\x1b[[C": F3, // Linux console
	"\x1b[[D": F4, // Linux console
	"\x1b[[E": F5, // Linux console

	"\x1b[\x11~": F1, // rxvt-unicode
	"\x1b[\x12~": F2, // rxvt-unicode
	"\x1b[\x13~": F3, // rxvt-unicode
	"\x1b[\x14~": F4, // rxvt-unicode

	"\x1b[15~":     F5,
	"\x1b[17~":     F6,
	"\x1b[18~":     F7,
	"\x1b[19~":     F8,
	"\x1b[20~":     F9,
	"\x1b[21~":     F10,
	"\x1b[22~":     F11,
	"\x1b[24~\x08": F12,
	"\x1b[\x25~":   F13,
	"\x1b[\x26~":   F14,
	"\x1b[\x28~":   F15,
	"\x1b[\x29~":   F16,
	// conflict `Home`: "\x1b[1~":      F17,
	// conflict `Insert`: "\x1b[2~":      F18,
	// conflict `Delete`: "\x1b[3~":      F19,
	// conflict `End`: "\x1b[4~":      F20,

	// Xterm
	"\x1b[1;2P": F13,
	"\x1b[1;2Q": F14,
	// &ASCIICode"\x1b[1;2\x52": F15,  // Conflicts with CPR response
	"\x1b[1;2R":    F16,
	"\x1b[\x15;2~": F17,
	"\x1b[\x17;2~": F18,
	"\x1b[\x18;2~": F19,
	"\x1b[\x19;2~": F20,
	"\x1b[\x20;2~": F21,
	"\x1b[\x21;2~": F22,
	"\x1b[\x23;2~": F23,
	"\x1b[\x24;2~": F24,

	"\x1b[1;5A": Control | Up,
	"\x1b[1;5B": Control | Down,
	"\x1b[1;5C": Control | Right,
	"\x1b[1;5D": Control | Left,

	"\x1b[1;2A": Shift | Up,
	"\x1b[1;2B": Shift | Down,
	"\x1b[1;2C": Shift | Right,
	"\x1b[1;2D": Shift | Left,

	// Tmux sends following keystrokes when control+arrow is pressed, but for
	// Emacs ansi-term sends the same sequences for normal arrow keys. Consider
	// it a normal arrow press, because that's more important.
	"\x1b0A": Up,
	"\x1b0B": Down,
	"\x1b0C": Right,
	"\x1b0D": Left,

	"\x1b[5A": Control | Up,
	"\x1b[5B": Control | Down,
	"\x1b[5C": Control | Right,
	"\x1b[5D": Control | Left,

	"\x1b[Oc": Control | Right, // rxvt
	"\x1b[Od": Control | Left,  // rxvt

	"\x1b[E": Ignore, // Xterm
	// conflict with 'End':  "\x1b[F": Ignore, // Linux console

	"\x1ba": Alt | A,
	"\x1bb": Alt | B,
	"\x1bc": Alt | C,
	"\x1bd": Alt | D,
	"\x1be": Alt | E,
	"\x1bf": Alt | F,
	"\x1bg": Alt | G,
	"\x1bh": Alt | H,
	"\x1bi": Alt | I,
	"\x1bj": Alt | J,
	"\x1bk": Alt | K,
	"\x1bl": Alt | L,
	"\x1bm": Alt | M,
	"\x1bn": Alt | N,
	"\x1bo": Alt | O,
	"\x1bp": Alt | P,
	"\x1bq": Alt | Q,
	"\x1br": Alt | R,
	"\x1bs": Alt | S,
	"\x1bt": Alt | T,
	"\x1bu": Alt | U,
	"\x1bv": Alt | V,
	"\x1bw": Alt | W,
	"\x1bx": Alt | X,
	"\x1by": Alt | Y,
	"\x1bz": Alt | Z,

	"\x1b\x08":  Alt | Backspace,
	"\x1b\x13":  Alt | Enter,
	"\x1b[3;3~": Alt | Delete,
}
