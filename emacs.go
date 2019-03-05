package prompt

import (
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

* [ ] Ctrl + t   Swap the last two characters before the cursor.
* [ ] Alt + t    Swap the word before the cursor with the word on/after the cursor.

* [ ] Ctrl + y   Paste (yank) the last thing to be cut.
* [ ] Ctrl + _   Undo.

* [*] Ctrl + Del Delete word after cursor
* [*] Ctrl + BS  Delete word before cursor
* [x] Alt + BS   Delete word before cursor
* [x] Alt + Del  Delete word before cursor
* [x] Alt + f    Move cursor to beginning of current/next word.
* [x] Alt + b    Move cursor to beginning of current/previous word.

* [ ] Alt + c    Capitalize word and move to next word.
* [ ] Alt + u    Uppercase word and move to next word.
* [ ] Alt + l    Lowercase word and move to next word.

*/

var emacsKeyBindings = map[KeyCode]KeyBindFunc{
	// Go to the End of the line
	Control | E: func(e *Event) { end_of_line(e.Buffer()) },
	// Go to the beginning of the line
	Control | A: func(e *Event) { beginning_of_line(e.Buffer()) },
	// Cut from cursor to the end of the Line
	Control | K: func(e *Event) {
		// TODO: kill_line()
		x := []rune(e.Buffer().Document().TextAfterCursor())
		e.Buffer().Delete(len(x))
	},
	// Cut from cursor to the beginning of the line
	Control | U: func(e *Event) {
		// TODO: backward_kill_line()
		x := []rune(e.Buffer().Document().TextBeforeCursor())
		e.Buffer().DeleteBeforeCursor(len(x))
	},
	// Delete character under the cursor
	Control | D: func(e *Event) {
		if len(e.Buffer().Text()) > 0 {
			delete_char(e.Buffer())
		} else {
			// pressing C-d on an empty edit means EOF
			e.Buffer().SetEOF()
		}
	},
	// Clear the Screen, similar to the clear command
	Control | L: func(*Event) {
		consoleWriter.EraseScreen()
		consoleWriter.CursorGoTo(0, 0)
		debug.AssertNoError(consoleWriter.Flush())
	},
	// Backspace
	Control | H: func(e *Event) { backward_delete_char(e.Buffer()) },
	// Right arrow: Forward one character
	Control | F: func(e *Event) { forward_char(e.Buffer()) },
	// Left arrow: Backward one character
	Control | B: func(e *Event) { backward_char(e.Buffer()) },
	// Cut the Word before the cursor.
	Control | W:         func(e *Event) { backward_kill_word(e.Buffer()) },
	Alt | F:             func(e *Event) { forward_word(e.Buffer()) },
	Alt | B:             func(e *Event) { backward_word(e.Buffer()) },
	Alt | Backspace:     func(e *Event) { backward_kill_word(e.Buffer()) },
	Alt | Delete:        func(e *Event) { backward_kill_word(e.Buffer()) },
	Control | Delete:    func(e *Event) { kill_word(e.Buffer()) },
	Control | Backspace: func(e *Event) { backward_kill_word(e.Buffer()) },
	Control | Left:      func(e *Event) { backward_word(e.Buffer()) },
	Control | Right:     func(e *Event) { forward_word(e.Buffer()) },
}
