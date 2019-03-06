package prompt

import (
	"strings"

	"github.com/c-bata/go-prompt/internal/debug"
)

type StateFlags struct {
	endEdit       bool
	eof           bool
	translatedKey KeyCode
}

// Buffer emulates the console buffer.
type Buffer struct {
	text            string
	cursorPosition  int // relative to the absolute beginning
	cacheDocument   *Document
	preferredColumn int // Remember the original column for the next up/down movement.
	flags           StateFlags
}

// Text returns string of the current line.
func (b *Buffer) Text() string {
	return b.text
}

// Document method to return document instance from the current text and cursor position.
func (b *Buffer) Document() (d *Document) {
	if b.cacheDocument == nil ||
		b.cacheDocument.Text != b.Text() ||
		b.cacheDocument.cursorPosition != b.cursorPosition {
		b.cacheDocument = &Document{
			Text:           b.Text(),
			cursorPosition: b.cursorPosition,
		}
	}
	return b.cacheDocument
}

// useful for keybind functions
// TODO: there should be a formal API exposed to these functions
//   e.g. an 'event' object (as in python prompt-toolkit), where event.CurrentBuffer()
//        retrieves the buffer now supplied.
//   called named functions, in readline style e.g. "backward-delete-char"

func (b *Buffer) setEOF() {
	b.flags.eof = true
}

func (b *Buffer) setEndEdit() {
	b.flags.endEdit = true
}

func (b *Buffer) setTranslatedKey(key KeyCode) {
	b.flags.translatedKey = key
}

// DisplayCursorPosition returns the cursor position on rendered text on terminal emulators.
// So if Document is "日本(cursor)語", 4 is returned because '日' and '本' are double width characters.
func (b *Buffer) DisplayCursorPosition() int {
	return b.Document().DisplayCursorPosition()
}

// DisplayCursorCoord returns similar to DisplayCursorPosition but separate col & row.
func (b *Buffer) DisplayCursorCoord(termWidth int) Coord {
	return b.Document().DisplayCursorCoord(termWidth)
}

// InsertText insert string from current line.
func (b *Buffer) InsertText(v string, overwrite bool, moveCursor bool) {
	or := []rune(b.Text())
	oc := b.cursorPosition

	if overwrite {
		overwritten := string(or[oc : oc+len(v)])
		if strings.Contains(overwritten, "\n") {
			i := strings.IndexAny(overwritten, "\n")
			overwritten = overwritten[:i]
		}
		b.setText(string(or[:oc]) + v + string(or[oc+len(overwritten):]))
	} else {
		b.setText(string(or[:oc]) + v + string(or[oc:]))
	}

	if moveCursor {
		b.cursorPosition += len([]rune(v))
		b.preferredColumn = b.Document().CursorPositionCol()
	}
}

// SetText method to set text and update cursorPosition.
// (When doing this, make sure that the position is valid for this text.
// text/cursor position should be consistent at any time, otherwise set a Document instead.)
func (b *Buffer) setText(v string) {
	debug.Assert(b.cursorPosition <= len([]rune(v)), "length of input should be shorter than cursor position")
	old := b.text
	b.text = v

	if old != v {
		// Text is changed.
		// TODO: Call callback function triggered by text changed. And also history search text should reset.
		// https://github.com/jonathanslenders/python-prompt-toolkit/blob/master/prompt_toolkit/buffer.py#L380-L384
	}
}

// Set cursor position. Return whether it changed.
func (b *Buffer) setCursorPosition(p int) {
	o := b.cursorPosition
	if p > 0 {
		b.cursorPosition = p
	} else {
		b.cursorPosition = 0
	}
	if p != o {
		// Cursor position is changed.
		// TODO: Call a onCursorPositionChanged function.
	}
}

func (b *Buffer) setDocument(d *Document) {
	b.cacheDocument = d
	b.setCursorPosition(d.cursorPosition) // Call before setText because setText check the relation between cursorPosition and line length.
	b.setText(d.Text)
}

// CursorLeft move to left on the current line.
func (b *Buffer) CursorLeft(count int) {
	l := b.Document().GetCursorLeftPosition(count)
	b.cursorPosition += l
	b.preferredColumn = b.Document().CursorPositionCol()
	return
}

// CursorRight move to right on the current line.
func (b *Buffer) CursorRight(count int) {
	l := b.Document().GetCursorRightPosition(count)
	b.cursorPosition += l
	b.preferredColumn = b.Document().CursorPositionCol()
	return
}

// CursorUp move cursor to the previous line.
// (for multi-line edit).
func (b *Buffer) CursorUp(count int) {
	orig := b.preferredColumn
	if b.preferredColumn == -1 { // -1 means not set
		orig = b.Document().CursorPositionCol()
	}
	b.cursorPosition += b.Document().GetCursorUpPosition(count, orig)
}

// CursorDown move cursor to the next line.
// (for multi-line edit).
func (b *Buffer) CursorDown(count int) {
	orig := b.preferredColumn
	if b.preferredColumn == -1 { // -1 means not set
		orig = b.Document().CursorPositionCol()
	}
	b.cursorPosition += b.Document().GetCursorDownPosition(count, orig)
}

// DeleteBeforeCursor delete specified number of characters before cursor and return the deleted text.
func (b *Buffer) DeleteBeforeCursor(count int) (deleted string) {
	debug.Assert(count >= 0, "count should be positive")
	r := []rune(b.Text())

	if b.cursorPosition > 0 {
		start := b.cursorPosition - count
		if start < 0 {
			start = 0
		}
		deleted = string(r[start:b.cursorPosition])
		b.setDocument(&Document{
			Text:           string(r[:start]) + string(r[b.cursorPosition:]),
			cursorPosition: b.cursorPosition - len([]rune(deleted)),
		})
	}
	b.preferredColumn = b.Document().CursorPositionCol()
	return
}

// NewLine means CR.
func (b *Buffer) NewLine(copyMargin bool) {
	// this must also output a '\n' to move the cursor down one line.
	// btw, Output.CursorDown(1) would not hack it (doesn't move if we're already at the bottom)
	//   we also don't have access to it.
	if copyMargin {
		b.InsertText("\n"+b.Document().leadingWhitespaceInCurrentLine(), false, true)
	} else {
		b.InsertText("\n", false, true)
	}
}

// Delete specified number of characters and Return the deleted text.
func (b *Buffer) Delete(count int) (deleted string) {
	r := []rune(b.Text())
	if b.cursorPosition < len(r) {
		deleted = b.Document().TextAfterCursor()[:count]
		b.setText(string(r[:b.cursorPosition]) + string(r[b.cursorPosition+len(deleted):]))
	}
	return
}

// JoinNextLine joins the next line to the current one by deleting the line ending after the current line.
func (b *Buffer) JoinNextLine(separator string) {
	if !b.Document().CursorOnLastLine() {
		b.cursorPosition += b.Document().GetEndOfLinePosition()
		b.Delete(1)
		// Remove spaces
		b.setText(b.Document().TextBeforeCursor() + separator + strings.TrimLeft(b.Document().TextAfterCursor(), " "))
	}
}

// SwapCharactersBeforeCursor swaps the last two characters before the cursor.
func (b *Buffer) SwapCharactersBeforeCursor() {
	if b.cursorPosition >= 2 {
		x := b.Text()[b.cursorPosition-2 : b.cursorPosition-1]
		y := b.Text()[b.cursorPosition-1 : b.cursorPosition]
		b.setText(b.Text()[:b.cursorPosition-2] + y + x + b.Text()[b.cursorPosition:])
	}
}

// NewBuffer is constructor of Buffer struct.
func NewBuffer() (b *Buffer) {
	b = &Buffer{
		preferredColumn: -1, // current / don't care
	}
	return
}
