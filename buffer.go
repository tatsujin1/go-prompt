package prompt

import (
	"strings"
	"sync"

	"github.com/tatsujin/go-prompt/internal/debug"
)

type StateFlags struct {
	endEdit       bool
	eof           bool
	translatedKey KeyCode
}

// Buffer emulates the console buffer.
type Buffer struct {
	text            string
	textLock        *sync.RWMutex
	cursor          Index // absolute character index into 'text'
	preferredColumn Index // preferred column for the next up/down movement.
	flags           StateFlags

	cacheDocument *Document
}

// NewBuffer is constructor of Buffer struct.
func NewBuffer() *Buffer {
	return &Buffer{
		textLock: &sync.RWMutex{},
	}
}

func (b *Buffer) RLock() {
	b.textLock.RLock()
}

func (b *Buffer) RUnlock() {
	b.textLock.RUnlock()
}

// Text returns string of the current line.
func (b *Buffer) Text() string {
	b.textLock.RLock()
	defer b.textLock.RUnlock()

	return b.text
}

func (b *Buffer) IsEmpty() bool {
	return len(b.text) == 0
}

// Document method to return document instance from the current text and cursor position.
func (b *Buffer) Document() *Document {
	b.textLock.RLock()
	defer b.textLock.RUnlock()

	return b.document()
}

func (b *Buffer) document() *Document {
	if b.cacheDocument == nil || b.cacheDocument.CursorIndex() != b.cursor || b.cacheDocument.Text() != b.text {
		b.cacheDocument = NewDocument(b.text, b.CursorIndex())
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

// CursorIndex returns the current cursor position byte index (0-based).
func (b *Buffer) CursorIndex() (index Index) {
	return b.cursor
}

// DisplayCursorPosition returns the cursor position on rendered text on terminal emulators.
// So if Document is "日本(cursor)語", 4 is returned because '日' and '本' are double width characters.
func (b *Buffer) CursorTextColumn() Column {
	return b.Document().CursorTextColumn()
}

// DisplayCursorCoord returns similar to DisplayCursorPosition but separate col & row.
func (b *Buffer) CursorDisplayCoord(termWidth Column) Coord {
	return b.Document().CursorDisplayCoord(termWidth)
}

// InsertText insert string from current line.
func (b *Buffer) InsertText(v string, overwrite bool, moveCursor bool) {
	b.textLock.Lock()
	defer b.textLock.Unlock()

	b.insertText(v, overwrite, moveCursor)
}

func (b *Buffer) insertText(v string, overwrite bool, moveCursor bool) {
	or := []rune(b.text)
	oc := b.cursor

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
		b.cursor += len([]rune(v))
		b.preferredColumn = b.document().CursorColumnIndex()
	}
}

// SetText method to set text and update cursor.
// (When doing this, make sure that the position is valid for this text.
// text/cursor position should be consistent at any time, otherwise set a Document instead.)
func (b *Buffer) setText(v string) {
	debug.Assert(b.cursor <= len([]rune(v)), "length of input should be shorter than cursor position")
	old := b.text
	b.text = v

	if old != v {
		// Text is changed.
		// TODO: Call callback function triggered by text changed. And also history search text should reset.
		// https://github.com/jonathanslenders/python-prompt-toolkit/blob/master/prompt_toolkit/buffer.py#L380-L384
	}
}

// Set cursor position. Return whether it changed.
func (b *Buffer) setCursorIndex(p int) {
	o := b.cursor
	if p > 0 {
		b.cursor = p
	} else {
		b.cursor = 0
	}
	if p != o {
		// Cursor position is changed.
		// TODO: Call a onCursorIndexChanged function.
	}
}

func (b *Buffer) setDocument(d *Document) {
	b.cacheDocument = d
	b.setCursorIndex(d.cursor) // Call before setText because setText check the relation between cursor and line length.
	b.setText(d.Text())
}

// CursorPrev cursor 'count' bytes backwards in the text (might cross lines).
func (b *Buffer) CursorPrev(count int) {
	if count < 0 {
		b.CursorNext(-count)
	} else if count > b.cursor {
		b.cursor = 0
	} else {
		b.cursor -= count
	}
}

// CursorNext cursor 'count' bytes forwards in the text (might cross lines).
func (b *Buffer) CursorNext(count int) {
	b.textLock.RLock()
	defer b.textLock.RUnlock()

	if count < 0 {
		b.CursorPrev(-count)
	} else if b.cursor > len(b.text) {
		b.cursor = len(b.text)
	} else {
		b.cursor += count
	}
}

// CursorLeft move to left on the current line.
func (b *Buffer) CursorLeft(count Offset) {
	b.textLock.RLock()
	defer b.textLock.RUnlock()

	b.cursor += Index(b.document().GetCursorLeftOffset(count))
	b.preferredColumn = b.document().CursorColumnIndex()
}

// CursorRight move to right on the current line.
func (b *Buffer) CursorRight(count Offset) {
	b.textLock.RLock()
	defer b.textLock.RUnlock()

	b.cursor += Index(b.document().GetCursorRightOffset(count))
	b.preferredColumn = b.document().CursorColumnIndex()
}

// CursorUp move cursor to the previous line.
// (for multi-line edit).
func (b *Buffer) CursorUp(count Row) {
	b.cursor += Index(b.Document().GetCursorUpOffset(count, b.preferredColumn))
}

// CursorDown move cursor to the next line.
// (for multi-line edit).
func (b *Buffer) CursorDown(count Row) {
	b.cursor += Index(b.Document().GetCursorDownOffset(count, b.preferredColumn))
}

// DeleteBeforeCursor delete specified number of characters before cursor and return the deleted text.
func (b *Buffer) DeleteBeforeCursor(count Offset) (deleted string) {
	b.textLock.Lock()
	defer b.textLock.Unlock()

	debug.Assert(count >= 0, "count should be positive")
	r := []rune(b.text)

	if b.cursor > 0 {
		start := b.cursor - Index(count)
		if start < 0 {
			start = 0
		}
		deleted = string(r[start:b.cursor])
		b.setDocument(NewDocument(
			string(r[:start])+string(r[b.cursor:]),
			b.cursor-len([]rune(deleted)),
		))
	}
	b.preferredColumn = b.document().CursorColumnIndex()
	return
}

// Delete specified number of characters and Return the deleted text.
func (b *Buffer) Delete(count Offset) (deleted string) {
	b.textLock.Lock()
	defer b.textLock.Unlock()

	return b.delete(count)
}

func (b *Buffer) delete(count Offset) (deleted string) {
	r := []rune(b.text)
	if b.cursor < len(r) {
		deleted = b.document().TextAfterCursor()[:count]
		b.setText(string(r[:b.cursor]) + string(r[b.cursor+len(deleted):]))
	}
	return deleted
}

// NewLine means CR.
func (b *Buffer) NewLine(copyMargin bool) {
	b.textLock.Lock()
	defer b.textLock.Unlock()

	// this must also output a '\n' to move the cursor down one line.
	// btw, Output.CursorDown(1) would not hack it (doesn't move if we're already at the bottom)
	//   we also don't have access to it.
	if copyMargin {
		b.insertText("\n"+b.document().leadingWhitespaceInCurrentLine(), false, true)
	} else {
		b.insertText("\n", false, true)
	}
}

// JoinNextLine joins the next line to the current one by deleting the line ending after the current line.
func (b *Buffer) JoinNextLine(separator string) {
	b.textLock.Lock()
	defer b.textLock.Unlock()

	if !b.document().CursorOnLastLine() {
		b.cursor += Index(b.document().GetEndOfLineOffset())
		b.delete(1)
		// Remove spaces
		b.setText(b.document().TextBeforeCursor() + separator + strings.TrimLeft(b.document().TextAfterCursor(), " "))
	}
}

// SwapCharactersBeforeCursor swaps the last two characters before the cursor.
func (b *Buffer) SwapCharactersBeforeCursor() {
	if b.cursor >= 2 {
		b.textLock.Lock()
		defer b.textLock.Unlock()

		x := b.text[b.cursor-2 : b.cursor-1]
		y := b.text[b.cursor-1 : b.cursor]
		b.setText(b.text[:b.cursor-2] + y + x + b.text[b.cursor:])
	}
}
