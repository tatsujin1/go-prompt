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

* [x] Ctrl + y   Paste (yank) the last thing to be cut.
* [ ] Ctrl + _   Undo.

* [x] Ctrl + Del Delete word after cursor
* [x] Ctrl + BS  Delete word before cursor
* [x] Alt + BS   Delete word before cursor
* [x] Alt + Del  Delete word before cursor
* [x] Alt + f    Move cursor to beginning of current/next word.
* [x] Alt + b    Move cursor to beginning of current/previous word.
* [x] Alt + d    Delete word before cursor

* [ ] Alt + c    Capitalize word and move to next word.
* [ ] Alt + u    Uppercase word and move to next word.
* [ ] Alt + l    Lowercase word and move to next word.

*/

var emacsKeyBindings = map[KeyCode]KeyBindFunc{
	// Go to the End of the line
	KeyControl | KeyE: end_of_line,
	// Go to the beginning of the line
	KeyControl | KeyA: beginning_of_line,
	// Cut from cursor to the end of the Line
	KeyControl | KeyK: kill_line,
	// Cut from cursor to the beginning of the line
	KeyControl | KeyU: backward_kill_line,
	// Delete character under the cursor
	KeyControl | KeyD: func(e *Event) {
		if e.Buffer().IsEmpty() {
			// pressing C-d in an empty edit means EOF
			e.SetEOF()
		} else {
			delete_char(e)
		}
	},
	// Clear the Screen, similar to the clear command
	KeyControl | KeyL: func(*Event) {
		// TODO: clear_screen(e.Buffer())
		consoleWriter.EraseScreen()
		consoleWriter.CursorGoTo(0, 0)
		debug.AssertNoError(consoleWriter.Flush())
	},
	KeyControl | KeyH:     backward_delete_char,
	KeyControl | KeyF:     forward_char,
	KeyControl | KeyB:     backward_char,
	KeyControl | KeyW:     backward_kill_word,
	KeyControl | KeyY:     yank,
	KeyAlt | KeyF:         forward_word,
	KeyAlt | KeyB:         backward_word,
	KeyAlt | KeyBackspace: backward_kill_word,
	KeyAlt | KeyDelete:    backward_kill_word,
	KeyAlt | KeyD:         kill_word,
	//KeyAlt | KeyW:             kill_ring_save,
	KeyControl | KeyDelete:    kill_word,
	KeyControl | KeyBackspace: backward_kill_word,
	KeyControl | KeyLeft:      backward_word,
	KeyControl | KeyRight:     forward_word,
	/*KeyAlt | KeyC: func(e *Event) {
		fmt.Fprintln(os.Stderr, "M-c")
	},
	KeyAlt | KeyU: func(e *Event) {
		fmt.Fprintln(os.Stderr, "M-u")
	},
	KeyAlt | KeyL: func(e *Event) {
		fmt.Fprintln(os.Stderr, "M-l")
	},*/
}
