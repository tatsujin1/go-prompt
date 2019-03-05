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
	Control | E: func(e *Event) KeyBindResult { end_of_line(e.Buffer()); return nil },
	// Go to the beginning of the line
	Control | A: func(e *Event) KeyBindResult { beginning_of_line(e.Buffer()); return nil },
	// Cut from cursor to the end of the Line
	Control | K: func(e *Event) KeyBindResult {
		// TODO: kill_line()
		x := []rune(e.Buffer().Document().TextAfterCursor())
		e.Buffer().Delete(len(x))
		return nil
	},
	// Cut from cursor to the beginning of the line
	Control | U: func(e *Event) KeyBindResult {
		// TODO: backward_kill_line()
		x := []rune(e.Buffer().Document().TextBeforeCursor())
		e.Buffer().DeleteBeforeCursor(len(x))
		return nil
	},
	// Delete character under the cursor
	Control | D: func(e *Event) KeyBindResult {
		if len(e.Buffer().Text()) > 0 {
			delete_char(e.Buffer())
		} else {
			// pressing C-d on an empty edit means EOF
			e.Buffer().SetEOF()
		}
		return nil
	},
	// Clear the Screen, similar to the clear command
	Control | L: func(*Event) KeyBindResult {
		consoleWriter.EraseScreen()
		consoleWriter.CursorGoTo(0, 0)
		debug.AssertNoError(consoleWriter.Flush())
		return nil
	},
	// Backspace
	Control | H: func(e *Event) KeyBindResult { backward_delete_char(e.Buffer()); return nil },
	// Right arrow: Forward one character
	Control | F: func(e *Event) KeyBindResult { forward_char(e.Buffer()); return nil },
	// Left arrow: Backward one character
	Control | B: func(e *Event) KeyBindResult { backward_char(e.Buffer()); return nil },
	// Cut the Word before the cursor.
	Control | W:         func(e *Event) KeyBindResult { backward_kill_word(e.Buffer()); return nil },
	Alt | F:             func(e *Event) KeyBindResult { forward_word(e.Buffer()); return nil },
	Alt | B:             func(e *Event) KeyBindResult { backward_word(e.Buffer()); return nil },
	Alt | Backspace:     func(e *Event) KeyBindResult { backward_kill_word(e.Buffer()); return nil },
	Alt | Delete:        func(e *Event) KeyBindResult { backward_kill_word(e.Buffer()); return nil },
	Control | Delete:    func(e *Event) KeyBindResult { kill_word(e.Buffer()); return nil },
	Control | Backspace: func(e *Event) KeyBindResult { backward_kill_word(e.Buffer()); return nil },
}
