package prompt

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

// delete_word Delete word before the cursor
func delete_word(buf *Buffer) {
	buf.DeleteBeforeCursor(len([]rune(buf.Document().TextBeforeCursor())) - buf.Document().FindStartOfPreviousWordWithSpace())
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
	buf.CursorRight(buf.Document().FindEndOfCurrentWordWithSpace())
}

// backward_word Backward one word
func backward_word(buf *Buffer) {
	buf.CursorLeft(len([]rune(buf.Document().TextBeforeCursor())) - buf.Document().FindStartOfPreviousWordWithSpace())
}

// delete and copy word at cursor
func kill_word(buf *Buffer) {
	// TODO: if cursor is at the end of the line (and there is a following line),
	//   join with the next line and call again on the new buffer.
	doc := buf.Document()
	wend := doc.FindEndOfCurrentWordWithSpace()
	deleted := buf.Delete(wend - doc.CursorPositionCol())
	// TODO: copy 'deleted' to clipboard
	_ = deleted
}

// delete and copy word before cursor
func backward_kill_word(buf *Buffer) {
	// TODO: if cursor is at the beginning of the line (and there is a preceeding line),
	//   join with previous line and call again on new buffer.
	deleted := buf.DeleteBeforeCursor(len([]rune(buf.Document().GetWordBeforeCursorWithSpace())))
	//doc := buf.Document()
	//wstart := doc.FindStartOfPreviousWordWithSpace()
	//deleted := buf.DeleteBeforeCursor(doc.CursorPositionCol() - wstart)
	// TODO: copy 'deleted' to clipboard
	_ = deleted
}
