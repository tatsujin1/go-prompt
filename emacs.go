package prompt

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt/internal/debug"
)

/*

========
PROGRESS
========

Moving the cursor
-----------------

* [x] Ctrl + a   Go to the beginning of the line (Home)
* [x] Ctrl + e   Go to the End of the line (End)
* [x] Ctrl + p   Previous command (Up arrow)
* [x] Ctrl + n   Next command (Down arrow)
* [x] Ctrl + f   Forward one character
* [x] Ctrl + b   Backward one character
* [x] Ctrl + xx  Toggle between the start of line and current cursor position

Editing
-------

* [x] Ctrl + L   Clear the Screen, similar to the clear command
* [x] Ctrl + d   Delete character under the cursor
* [x] Ctrl + h   Delete character before the cursor (Backspace)

* [x] Ctrl + w   Cut the Word before the cursor to the clipboard.
* [x] Ctrl + k   Cut the Line after the cursor to the clipboard.
* [x] Ctrl + u   Cut/delete the Line before the cursor to the clipboard.

* [ ] Ctrl + t   Swap the last two characters before the cursor (typo).
* [ ] Esc  + t   Swap the last two words before the cursor.

* [ ] ctrl + y   Paste the last thing to be cut (yank)
* [ ] ctrl + _   Undo

*/

var emacsKeyBindings = map[KeyCode]KeyBindFunc{
	// Go to the End of the line
	ControlE: func(buf *Buffer) {
		x := []rune(buf.Document().TextAfterCursor())
		buf.CursorRight(len(x))
	},
	// Go to the beginning of the line
	ControlA: func(buf *Buffer) {
		x := []rune(buf.Document().TextBeforeCursor())
		buf.CursorLeft(len(x))
	},
	// Cut the Line after the cursor
	ControlK: func(buf *Buffer) {
		x := []rune(buf.Document().TextAfterCursor())
		buf.Delete(len(x))
	},
	// Cut/delete the Line before the cursor
	ControlU: func(buf *Buffer) {
		x := []rune(buf.Document().TextBeforeCursor())
		buf.DeleteBeforeCursor(len(x))
	},
	// Delete character under the cursor
	ControlD: func(buf *Buffer) {
		if buf.Text() != "" {
			buf.Delete(1)
		}
	},
	// Backspace
	ControlH: func(buf *Buffer) {
		buf.DeleteBeforeCursor(1)
	},
	// Right allow: Forward one character
	ControlF: func(buf *Buffer) {
		buf.CursorRight(1)
	},
	// Left allow: Backward one character
	ControlB: func(buf *Buffer) {
		buf.CursorLeft(1)
	},
	// Cut the Word before the cursor.
	ControlW: func(buf *Buffer) {
		buf.DeleteBeforeCursor(len([]rune(buf.Document().GetWordBeforeCursorWithSpace())))
	},
	// Clear the Screen, similar to the clear command
	ControlL: func(buf *Buffer) {
		consoleWriter.EraseScreen()
		consoleWriter.CursorGoTo(0, 0)
		debug.AssertNoError(consoleWriter.Flush())
	},
	},
}
