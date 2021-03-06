// Code generated by hand; DO NOT EDIT.
// This is a little bit stupid, but there are many public constants which is no value for writing godoc comment.

package prompt

// KeyCode is the type represents a key inserted by the user.
//go:generate stringer -type=KeyCode
type KeyCode int32

type ControlSequence string

// Key is the type contains a logical Key and its control sequence.
type KeyDefinition struct {
	Key      KeyCode
	Sequence ControlSequence
}

const (
	// Key is not defined
	Undefined KeyCode = iota

	// Key which is ignored. (The key binding for this key should not do anything.)
	Ignore

	// Matches any key.
	KeyAny

	KeyEscape

	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ

	// must prefix these with 'Key' :(
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	Key0

	KeyBacktick
	KeyCircumflex
	KeyUnderscore
	KeyMinus
	KeyEquals
	KeyBackspace
	KeySquareOpen
	KeySquareClose
	KeySingleQuote
	KeyBackslash
	KeyLessThan
	KeySpace
	KeyComma
	KeyPoint
	KeySlash

	KeyUp
	KeyDown
	KeyRight
	KeyLeft

	KeyHome
	KeyEnd
	KeyDelete
	KeyShiftDelete
	KeyControlDelete
	KeyPageUp
	KeyPageDown
	KeyBackTab
	KeyInsert

	// Aliases.
	KeyTab
	KeyEnter
	// Actually Enter equals ControlM, not ControlJ,
	// However, in prompt_toolkit, we made the mistake of translating
	// \r into \n during the input, so everyone is now handling the
	// enter key by binding ControlJ.

	// From now on, it's better to bind `ASCII_SEQUENCES.Enter` everywhere,
	// because that's future compatible, and will still work when we
	// stop replacing \r by \n.

	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24

	// Key modifiers, must be last (before the back-compat defs)
	KeyShift   = 0x10000
	KeyControl = 0x20000
	KeyAlt     = 0x40000

	// Special (not actual keys)
	ReportCursorPosition
	Vt100MouseEvent
	WindowsMouseEvent
	BracketedPaste

	// Constants for backwards compatibility:

	ControlA = KeyControl | KeyA
	ControlB = KeyControl | KeyB
	ControlC = KeyControl | KeyC
	ControlD = KeyControl | KeyD
	ControlE = KeyControl | KeyE
	ControlF = KeyControl | KeyF
	ControlG = KeyControl | KeyG
	ControlH = KeyControl | KeyH
	ControlI = KeyControl | KeyI
	ControlJ = KeyControl | KeyJ
	ControlK = KeyControl | KeyK
	ControlL = KeyControl | KeyL
	ControlM = KeyControl | KeyM
	ControlN = KeyControl | KeyN
	ControlO = KeyControl | KeyO
	ControlP = KeyControl | KeyP
	ControlQ = KeyControl | KeyQ
	ControlR = KeyControl | KeyR
	ControlS = KeyControl | KeyS
	ControlT = KeyControl | KeyT
	ControlU = KeyControl | KeyU
	ControlV = KeyControl | KeyV
	ControlW = KeyControl | KeyW
	ControlX = KeyControl | KeyX
	ControlY = KeyControl | KeyY
	ControlZ = KeyControl | KeyZ

	ControlSpace       = KeyControl | KeySpace
	ControlBackslash   = KeyControl | KeyBackslash
	ControlSquareClose = KeyControl | KeySquareClose
	ControlCircumflex  = KeyControl | KeyCircumflex
	ControlUnderscore  = KeyControl | KeyUnderscore
	ControlLeft        = KeyControl | KeyLeft
	ControlRight       = KeyControl | KeyRight
	ControlUp          = KeyControl | KeyUp
	ControlDown        = KeyControl | KeyDown

	ShiftLeft  = KeyShift | KeyLeft
	ShiftUp    = KeyShift | KeyUp
	ShiftDown  = KeyShift | KeyDown
	ShiftRight = KeyShift | KeyRight
)
