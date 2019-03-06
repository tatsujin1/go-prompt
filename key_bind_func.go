package prompt

import (
	"fmt"
	"os"
)

var clipboard string

// end_of_line Go to the End of the line
func end_of_line(buf *Buffer) {
	doc := buf.Document()

	rest_len := len([]rune(doc.CurrentLineAfterCursor()))
	if rest_len > 0 {
		buf.CursorRight(rest_len)
	} else if len(doc.TextAfterCursor()) > 0 {
		buf.cursorPosition++
		end_of_line(buf)
	}
}

// beginning_of_line Go to the beginning of the line
func beginning_of_line(buf *Buffer) {
	doc := buf.Document()
	lead_len := len([]rune(doc.CurrentLineBeforeCursor()))
	if lead_len > 0 {
		buf.CursorLeft(lead_len)
	} else if len(doc.TextBeforeCursor()) > 0 {
		buf.cursorPosition-- // move before '\n'
		beginning_of_line(buf)
	}
}

// delete_char Delete character under the cursor
func delete_char(buf *Buffer) {
	buf.Delete(1)
}

// delete_word Delete word after the cursor
func delete_word(buf *Buffer) {
	wend := buf.Document().FindEndOfCurrentWordWithSpace()
	buf.Delete(wend)
}

// backward_delete_char Go to Backspace
func backward_delete_char(buf *Buffer) {
	buf.DeleteBeforeCursor(1)
}

// forward_char Forward one character
func forward_char(buf *Buffer) {
	buf.CursorRight(1)
}

// backward_char Backward one character
func backward_char(buf *Buffer) {
	buf.CursorLeft(1)
}

// forward_word Forward one word
func forward_word(buf *Buffer) {
	// TODO: if cursor is at the end of the line (and there is a following line),
	//   move cursor to the beginning of the following line and call again
	buf.CursorRight(buf.Document().FindEndOfCurrentWordWithSpace())
}

// backward_word Backward one word
func backward_word(buf *Buffer) {
	// TODO: if cursor is at the beginning of the line (and there is a preceeding line),
	//   move cursor to the end of the preceeding line and call again
	doc := buf.Document()
	wstart := doc.FindStartOfPreviousWordWithSpace()
	buf.CursorLeft(len([]rune(doc.TextBeforeCursor())) - wstart)
}

// delete and copy word at cursor
func kill_word(buf *Buffer) {
	// TODO: if cursor is at the end of the line (and there is a following line),
	//   join with the next line and call again on the new buffer.
	doc := buf.Document()

	//if ! doc.CursorOnLastLine() && doc.Cursor
	clipboard = buf.Delete(doc.FindEndOfCurrentWordWithSpace())
}

// delete and copy word before cursor
func backward_kill_word(buf *Buffer) {
	// TODO: if cursor is at the beginning of the line (and there is a preceeding line),
	//   join with previous line and call again on new buffer.
	clipboard = buf.DeleteBeforeCursor(len([]rune(buf.Document().GetWordBeforeCursorWithSpace())))
}

func kill_line(buf *Buffer) {
	doc := buf.Document()
	x := []rune(doc.CurrentLineAfterCursor())
	if len(x) > 0 {
		clipboard = buf.Delete(len(x))
	} else if !doc.CursorOnLastLine() {
		buf.DeleteBeforeCursor(1)
	}
}

func backward_kill_line(buf *Buffer) {
	x := []rune(buf.Document().TextBeforeCursor())
	clipboard = buf.DeleteBeforeCursor(len(x))
}

func yank(buf *Buffer) {
	// TODO: output bracketed paste ON ("\x1b[?2004h") during rendering
	buf.InsertText(clipboard, false, true)
	// TODO: output bracketed paste OFF ("\x1b[?2004l") during rendering
}
