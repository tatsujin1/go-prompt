package prompt

var clipboard string

// end_of_line Go to the End of the line
func end_of_line(e *Event) {
	buf := e.Buffer()
	doc := buf.Document()

	end_offset := doc.GetEndOfLineOffset()
	if end_offset > 0 {
		buf.CursorRight(end_offset)
	} else if len(doc.TextAfterCursor()) > 0 {
		buf.CursorDown(1) // to next line
		end_of_line(e)    // and then to the end of that line (must create new Document)
	}
}

// beginning_of_line Go to the beginning of the line
func beginning_of_line(e *Event) {
	buf := e.Buffer()
	doc := buf.Document()
	col := doc.CursorColumnIndex()
	if col > 0 {
		buf.CursorLeft(Offset(col))
	} else if doc.CursorRow() > 0 {
		buf.CursorUp(1)      // to previous line
		beginning_of_line(e) // and then to the beginning of that line (must create new Document)
	}
}

// delete_char Delete character under the cursor
func delete_char(e *Event) {
	buf := e.Buffer()
	buf.Delete(1)
}

// delete_word Delete word after the cursor
func delete_word(e *Event) {
	buf := e.Buffer()
	wend := buf.Document().FindEndOfCurrentWordWithSpace()
	buf.Delete(Offset(wend))
}

// backward_delete_char Go to Backspace
func backward_delete_char(e *Event) {
	buf := e.Buffer()
	buf.DeleteBeforeCursor(1)
}

// forward_char Forward one character
func forward_char(e *Event) {
	buf := e.Buffer()
	doc := buf.Document()

	// if cursor is at the end of the line (and there is a following line),
	//   move cursor to the beginning of the following line

	if !doc.CursorAtEndOfLine() {
		buf.CursorRight(1)
	} else if doc.CursorRow() < Row(doc.LineCount()-1) {
		buf.CursorDown(1)
		buf.CursorLeft(doc.GetBeginningOfLineOffset())
	}
}

// backward_char Backward one character
func backward_char(e *Event) {
	buf := e.Buffer()
	doc := buf.Document()

	// if cursor is at the beginning of the line (and there is a preceeding line),
	//   move cursor to the end of the preceeding line

	if len(doc.CurrentLineBeforeCursor()) > 0 {
		buf.CursorLeft(1)
	} else if doc.CursorRow() > 0 {
		buf.CursorUp(1)
		if !buf.Document().CursorAtEndOfLine() { // must create new Document
			end_of_line(e)
		}
	}
}

// forward_word Forward one word (across lines).
func forward_word(e *Event) {
	buf := e.Buffer()
	doc := buf.Document()

	wstart := doc.FindStartOfNextWord()
	if wstart == 0 { // nothing found at all -> go to end of text
		wstart = doc.GetEndOfTextOffset()
	}
	buf.CursorForward(wstart)
}

// backward_word Backward one word (across lines).
func backward_word(e *Event) {
	buf := e.Buffer()
	doc := buf.Document()

	wstart := doc.FindStartOfPreviousWord()
	buf.CursorBackward(Offset(doc.CursorIndex() - wstart))
}

// delete and copy word at cursor
func kill_word(e *Event) {
	buf := e.Buffer()
	// TODO: if cursor is at the end of the line (and there is a following line),
	//   join with the next line and call again on the new buffer.
	doc := buf.Document()

	//if ! doc.CursorOnLastLine() && doc.Cursor
	clipboard = buf.Delete(doc.FindEndOfCurrentWordWithSpace())
}

// delete and copy word before cursor
func backward_kill_word(e *Event) {
	buf := e.Buffer()
	// TODO: if cursor is at the beginning of the line (and there is a preceeding line),
	//   join with previous line and call again on new buffer.
	clipboard = buf.DeleteBeforeCursor(Offset(len([]rune(buf.Document().GetWordBeforeCursorWithSpace()))))
}

func kill_line(e *Event) {
	buf := e.Buffer()
	doc := buf.Document()
	x := []rune(doc.CurrentLineAfterCursor())
	if len(x) > 0 {
		clipboard = buf.Delete(Offset(len(x)))
	} else if !doc.CursorOnLastLine() {
		buf.DeleteBeforeCursor(1)
	}
}

func backward_kill_line(e *Event) {
	buf := e.Buffer()
	x := []rune(buf.Document().TextBeforeCursor())
	clipboard = buf.DeleteBeforeCursor(Offset(len(x)))
}

func yank(e *Event) {
	buf := e.Buffer()
	// TODO: output bracketed paste ON ("\x1b[?2004h") during rendering
	buf.InsertText(clipboard, false, true)
	// TODO: output bracketed paste OFF ("\x1b[?2004l") during rendering
}
