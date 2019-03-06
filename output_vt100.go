package prompt

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

// VT100Writer generates VT100 escape sequences.
type VT100Writer struct {
	buffer []byte
}

// WriteRaw to write raw byte array
func (w *VT100Writer) WriteRaw(data []byte) {
	w.buffer = append(w.buffer, data...)
	return
}

// Write to write safety byte array by removing control sequences.
func (w *VT100Writer) Write(data []byte) {
	w.WriteRaw(bytes.Replace(data, []byte{0x1b}, []byte{'?'}, -1))
	return
}

// WriteRawStr to write raw string
func (w *VT100Writer) WriteRawStr(data string) {
	w.WriteRaw([]byte(data))
	return
}

// WriteStr to write safety string by removing control sequences.
func (w *VT100Writer) WriteStr(data string) {
	w.Write([]byte(data))
	return
}

/* Erase */

// EraseScreen erases the screen with the background colour and moves the cursor to home (cursor doesn't move).
func (w *VT100Writer) EraseScreen() {
	w.WriteRaw([]byte{0x1b, '[', '2', 'J'})
	return
}

// EraseUp erases the screen from the current line up to the top of the screen (cursor doesn't move).
func (w *VT100Writer) EraseUp() {
	w.WriteRaw([]byte{0x1b, '[', '1', 'J'})
	return
}

// EraseDown erases the screen from the current line down to the bottom of the screen (cursor doesn't move).
func (w *VT100Writer) EraseDown() {
	w.WriteRaw([]byte{0x1b, '[', 'J'})
	return
}

// EraseStartOfLine erases from the current cursor position to the start of the current line (cursor doesn't move).
func (w *VT100Writer) EraseStartOfLine() {
	w.WriteRaw([]byte{0x1b, '[', '1', 'K'})
	return
}

// EraseEndOfLine erases from the current cursor position to the end of the current line (cursor doesn't move).
func (w *VT100Writer) EraseEndOfLine() {
	w.WriteRaw([]byte{0x1b, '[', 'K'})
	return
}

// EraseLine erases the entire current line (cursor doesn't move).
func (w *VT100Writer) EraseLine() {
	w.WriteRaw([]byte{0x1b, '[', '2', 'K'})
	return
}

/* Cursor */

// ShowCursor stops blinking cursor and show.
func (w *VT100Writer) ShowCursor() {
	w.WriteRaw([]byte{0x1b, '[', '?', '1', '2', 'l', 0x1b, '[', '?', '2', '5', 'h'})
}

// HideCursor hides cursor.
func (w *VT100Writer) HideCursor() {
	w.WriteRaw([]byte{0x1b, '[', '?', '2', '5', 'l'})
	return
}

// CursorGoTo sets the cursor position where subsequent text will begin.
func (w *VT100Writer) CursorGoTo(row, col int) {
	if row == 0 && col == 0 {
		// If no row/column parameters are provided (ie. <ESC>[H), the cursor will move to the home position.
		w.WriteRaw([]byte{0x1b, '[', 'H'})
		return
	}
	r := strconv.Itoa(row)
	c := strconv.Itoa(col)
	w.WriteRaw([]byte{0x1b, '['})
	w.WriteRaw([]byte(r))
	w.WriteRaw([]byte{';'})
	w.WriteRaw([]byte(c))
	w.WriteRaw([]byte{'H'})
	return
}

// CursorUp moves the cursor up by 'n' rows; the default count is 1.
func (w *VT100Writer) CursorUp(n int) {
	if n == 0 {
		return
	} else if n < 0 {
		w.CursorDown(-n)
		return
	}
	s := strconv.Itoa(n)
	w.WriteRaw([]byte{0x1b, '['})
	w.WriteRaw([]byte(s))
	w.WriteRaw([]byte{'A'})
	return
}

// CursorDown moves the cursor down by 'n' rows; the default count is 1.
func (w *VT100Writer) CursorDown(n int) {
	if n == 0 {
		return
	} else if n < 0 {
		w.CursorUp(-n)
		return
	}
	s := strconv.Itoa(n)
	w.WriteRaw([]byte{0x1b, '['})
	w.WriteRaw([]byte(s))
	w.WriteRaw([]byte{'B'})
	return
}

// CursorForward moves the cursor forward by 'n' columns; the default count is 1.
func (w *VT100Writer) CursorForward(n int) {
	if n == 0 {
		return
	} else if n < 0 {
		w.CursorBackward(-n)
		return
	}
	s := strconv.Itoa(n)
	w.WriteRaw([]byte{0x1b, '['})
	w.WriteRaw([]byte(s))
	w.WriteRaw([]byte{'C'})
	return
}

// CursorBackward moves the cursor backward by 'n' columns; the default count is 1.
func (w *VT100Writer) CursorBackward(n int) {
	if n == 0 {
		return
	} else if n < 0 {
		w.CursorForward(-n)
		return
	}
	s := strconv.Itoa(n)
	w.WriteRaw([]byte{0x1b, '['})
	w.WriteRaw([]byte(s))
	w.WriteRaw([]byte{'D'})
	return
}

// AskForCPR asks for a cursor position report (CPR).
func (w *VT100Writer) AskForCPR() {
	// CPR: Cursor Position Request.
	w.WriteRaw([]byte{0x1b, '[', '6', 'n'})
	return
}

// SaveCursor saves current cursor position.
func (w *VT100Writer) SaveCursor() {
	w.WriteRaw([]byte{0x1b, '[', 's'})
	return
}

// RestoreCursor restores cursor position saved by the last SaveCursor.
func (w *VT100Writer) RestoreCursor() {
	w.WriteRaw([]byte{0x1b, '[', 'u'})
	return
}

/* Scrolling */

// ScrollDown scrolls display down one line.
func (w *VT100Writer) ScrollDown() {
	w.WriteRaw([]byte{0x1b, 'D'})
	return
}

// ScrollUp scroll display up one line.
func (w *VT100Writer) ScrollUp() {
	w.WriteRaw([]byte{0x1b, 'M'})
	return
}

/* Title */

// SetTitle sets a title of terminal window.
func (w *VT100Writer) SetTitle(title string) {
	titleBytes := []byte(title)
	patterns := []struct {
		from []byte
		to   []byte
	}{
		{
			from: []byte{0x13},
			to:   []byte{},
		},
		{
			from: []byte{0x07},
			to:   []byte{},
		},
	}
	for i := range patterns {
		titleBytes = bytes.Replace(titleBytes, patterns[i].from, patterns[i].to, -1)
	}

	w.WriteRaw([]byte{0x1b, ']', '2', ';'})
	w.WriteRaw(titleBytes)
	w.WriteRaw([]byte{0x07})
	return
}

// ClearTitle clears a title of terminal window.
func (w *VT100Writer) ClearTitle() {
	w.WriteRaw([]byte{0x1b, ']', '2', ';', 0x07})
	return
}

/* Font */

// SetColor sets text and background colors. and specify whether text is bold.
func (w *VT100Writer) SetColor(fg, bg Color, bold bool) {
	if bold {
		w.SetDisplayAttributes(fg, bg, DisplayBold)
	} else {
		// If using `DisplayDefualt`, it will be broken in some environment.
		// Details are https://github.com/c-bata/go-prompt/pull/85
		w.SetDisplayAttributes(fg, bg, DisplayReset)
	}
	return
}

// SetDisplayAttributes to set VT100 display attributes.
var CSI []byte
var end []byte
var sep []byte

func init() {
	CSI = []byte{0x1b, '['} // control sequence introducer
	end = []byte{'m'}
	sep = []byte{';'}
}

func (w *VT100Writer) SetDisplayAttributes(fg, bg Color, attrs ...DisplayAttribute) {
	w.WriteRaw(CSI)


	var ok bool
	var v string

	for i := range attrs {
		if v, ok = displayAttributeParameters[attrs[i]]; !ok {
			continue
		}
		w.WriteRawStr(v)
		w.WriteRaw(sep)
	}

	// foreground color
	if v, ok = foregroundANSIColors[fg]; !ok {
		v = foregroundANSIColors[DefaultColor]
	}
	w.WriteRawStr(v)

	w.WriteRaw(sep)

	// background color
	if v, ok = backgroundANSIColors[bg]; !ok {
		v = backgroundANSIColors[DefaultColor]
	}
	w.WriteRawStr(v)

	w.WriteRaw(end)

	return
}

var displayAttributeParameters = map[DisplayAttribute]string{
	DisplayReset:        "0",
	DisplayBold:         "1",
	DisplayLowIntensity: "2",
	DisplayItalic:       "3",
	DisplayUnderline:    "4",
	DisplayBlink:        "5",
	DisplayRapidBlink:   "6",
	DisplayReverse:      "7",
	DisplayInvisible:    "8",
	DisplayCrossedOut:   "9",
	DisplayDefaultFont:  "10",
	DisplayAltFont1:     "11",
	DisplayAltFont2:     "12",
	DisplayAltFont3:     "13",
	DisplayAltFont4:     "14",
	DisplayAltFont5:     "15",
	DisplayAltFont6:     "16",
	DisplayAltFont7:     "17",
	DisplayAltFont8:     "18",
	DisplayAltFont9:     "19",
}

var foregroundANSIColors = map[Color]string{
	DefaultColor: "39",

	// Low intensity.
	Black:   "30",
	Red:     "31",
	Green:   "32",
	Yellow:  "33",
	Blue:    "34",
	Magenta: "35",
	Cyan:    "36",
	Gray:    "37",

	// High intensity.
	BrightBlack:   "90",
	BrightRed:     "91",
	BrightGreen:   "92",
	BrightYellow:  "93",
	BrightBlue:    "94",
	BrightMagenta: "95",
	BrightCyan:    "96",
	White:         "97",
}

var backgroundANSIColors = map[Color]string{
	DefaultColor: "49",

	// Low intensity.
	Black:   "40",
	Red:     "41",
	Green:   "42",
	Yellow:  "43",
	Blue:    "44",
	Magenta: "45",
	Cyan:    "46",
	Gray:    "47",

	// High intensity
	BrightBlack:   "100",
	BrightRed:     "101",
	BrightGreen:   "102",
	BrightYellow:  "103",
	BrightBlue:    "104",
	BrightMagenta: "105",
	BrightCyan:    "106",
	White:         "107",
}
