package prompt

import (
	"strings"
	"unicode/utf8"

	istrings "github.com/c-bata/go-prompt/internal/strings"
	runewidth "github.com/mattn/go-runewidth"

	"github.com/c-bata/go-prompt/internal/bisect"
)

// Document is a read-only view of the current editor content
type Document struct {
	Text string
	// This represents a index in '[]rune(Text)'.
	// So if Document is "日本(cursor)語", cursorPosition is 2.
	// But DisplayedCursorPosition returns 4 because '日' and '本' are double width characters.
	cursorPosition int

	linesCache      []string
	startIndexCache []int
}

// NewDocument return the new empty document.
func NewDocument() *Document {
	return &Document{
		Text:           "",
		cursorPosition: 0,
	}
}

// DisplayCursorPosition returns the cursor position on rendered text on terminal emulators.
// So if Document is "日本(cursor)語", DisplayedCursorPosition returns 4 because '日' and '本' are double width characters.
func (d *Document) DisplayCursorPosition() int {
	var position int
	runes := []rune(d.Text)[:d.cursorPosition]
	for i := range runes {
		position += runewidth.RuneWidth(runes[i])
	}
	return position
}

// DisplayCursorCoord returns similar to DisplayCursorPosition, but separate col & row.
func (d *Document) DisplayCursorCoord(termWidth int) Coord {
	return display_coord(termWidth, d.cursorPosition, 0, d.Text)
}

// DisplayCursorCoordWithPrefix same as DisplayCursorCoord but with a 'prefix' taken into account.
func (d *Document) DisplayCursorCoordWithPrefix(termWidth int, prefix string) Coord {
	return display_coord(termWidth, d.cursorPosition, 1, prefix, d.Text)
}

func display_coord(termWidth int, cursorPos int, useAll int, texts ...string) Coord {
	// we're doing a little extra legwork here to avoid memory allocations

	c := Coord{}

	idx := 0
	for tidx, t := range texts {
		r := strings.NewReader(t)
		all := tidx < useAll
		if !all {
			idx = 0 // now, let's start at 0
		}
		for ; all || idx < cursorPos; idx++ {
			ch, _, err := r.ReadRune()
			if err != nil {
				break
			}
			if ch == '\n' {
				c.Y++
				c.Y += c.X / termWidth
				c.X = 0
			} else {
				c.X += runewidth.RuneWidth(ch)
			}
		}
	}
	// line-wrap the last (i.e. non-terminated) line
	c.Y += c.X / termWidth
	c.X = c.X % termWidth

	return c
}

// GetCharRelativeToCursor return character relative to cursor position, or empty string
func (d *Document) GetCharRelativeToCursor(offset int) (r rune) {
	s := d.Text
	cnt := 0

	for len(s) > 0 {
		cnt++
		r, size := utf8.DecodeRuneInString(s)
		if cnt == d.cursorPosition+offset {
			return r
		}
		s = s[size:]
	}
	return 0
}

// TextBeforeCursor returns the text before the cursor.
//   includes preceeding lines if cursor is not on the 1st line.
func (d *Document) TextBeforeCursor() string {
	r := []rune(d.Text)
	return string(r[:d.cursorPosition])
}

// TextAfterCursor returns the text after the cursor.
//   includes following lines if cursor is not on the last line.
func (d *Document) TextAfterCursor() string {
	r := []rune(d.Text)
	return string(r[d.cursorPosition:])
}

// GetWordBeforeCursor returns the word before the cursor.
// If we have whitespace before the cursor this returns an empty string.
func (d *Document) GetWordBeforeCursor() string {
	x := d.TextBeforeCursor()
	return x[d.FindStartOfPreviousWord():]
}

// GetWordAfterCursor returns the word after the cursor.
// If we have whitespace after the cursor this returns an empty string.
func (d *Document) GetWordAfterCursor() string {
	x := d.TextAfterCursor()
	return x[:d.FindEndOfCurrentWord()]
}

// GetWordBeforeCursorWithSpace returns the word before the cursor.
// Unlike GetWordBeforeCursor, it returns string containing space
func (d *Document) GetWordBeforeCursorWithSpace() string {
	x := d.TextBeforeCursor()
	return x[d.FindStartOfPreviousWordWithSpace():]
}

// GetWordAfterCursorWithSpace returns the word after the cursor.
// Unlike GetWordAfterCursor, it returns string containing space
func (d *Document) GetWordAfterCursorWithSpace() string {
	x := d.TextAfterCursor()
	return x[:d.FindEndOfCurrentWordWithSpace()]
}

// GetWordBeforeCursorUntilSeparator returns the text before the cursor until next separator.
func (d *Document) GetWordBeforeCursorUntilSeparator(sep string) string {
	x := d.TextBeforeCursor()
	return x[d.FindStartOfPreviousWordUntilSeparator(sep):]
}

// GetWordAfterCursorUntilSeparator returns the text after the cursor until next separator.
func (d *Document) GetWordAfterCursorUntilSeparator(sep string) string {
	x := d.TextAfterCursor()
	return x[:d.FindEndOfCurrentWordUntilSeparator(sep)]
}

// GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor returns the word before the cursor.
// Unlike GetWordBeforeCursor, it returns string containing space
func (d *Document) GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor(sep string) string {
	x := d.TextBeforeCursor()
	return x[d.FindStartOfPreviousWordUntilSeparatorIgnoreNextToCursor(sep):]
}

// GetWordAfterCursorUntilSeparatorIgnoreNextToCursor returns the word after the cursor.
// Unlike GetWordAfterCursor, it returns string containing space
func (d *Document) GetWordAfterCursorUntilSeparatorIgnoreNextToCursor(sep string) string {
	x := d.TextAfterCursor()
	return x[:d.FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor(sep)]
}

// FindStartOfPreviousWord returns an index relative to the cursor position
// pointing to the start of the previous word. Return 0 if nothing was found.
func (d *Document) FindStartOfPreviousWord() int {
	x := d.TextBeforeCursor()
	i := strings.LastIndexByte(x, ' ')
	if i != -1 {
		return i + 1
	}
	return 0
}

// FindStartOfPreviousWordWithSpace is almost the same as FindStartOfPreviousWord.
// The only difference is to ignore contiguous spaces.
func (d *Document) FindStartOfPreviousWordWithSpace() int {
	x := d.TextBeforeCursor()
	end := istrings.LastIndexNotByte(x, ' ')
	if end == -1 {
		return 0
	}

	start := strings.LastIndexByte(x[:end], ' ')
	if start == -1 {
		return 0
	}
	return start + 1
}

// FindStartOfPreviousWordUntilSeparator is almost the same as FindStartOfPreviousWord.
// But this can specify Separator. Return 0 if nothing was found.
func (d *Document) FindStartOfPreviousWordUntilSeparator(sep string) int {
	if sep == "" {
		return d.FindStartOfPreviousWord()
	}

	x := d.TextBeforeCursor()
	i := strings.LastIndexAny(x, sep)
	if i != -1 {
		return i + 1
	}
	return 0
}

// FindStartOfPreviousWordUntilSeparatorIgnoreNextToCursor is almost the same as FindStartOfPreviousWordWithSpace.
// But this can specify Separator. Return 0 if nothing was found.
func (d *Document) FindStartOfPreviousWordUntilSeparatorIgnoreNextToCursor(sep string) int {
	if sep == "" {
		return d.FindStartOfPreviousWordWithSpace()
	}

	x := d.TextBeforeCursor()
	end := istrings.LastIndexNotAny(x, sep)
	if end == -1 {
		return 0
	}
	start := strings.LastIndexAny(x[:end], sep)
	if start == -1 {
		return 0
	}
	return start + 1
}

// FindEndOfCurrentWord returns an index relative to the cursor position.
// pointing to the end of the current word. Return 0 if nothing was found.
func (d *Document) FindEndOfCurrentWord() int {
	x := d.TextAfterCursor()
	i := strings.IndexByte(x, ' ')
	if i != -1 {
		return i
	}
	return len(x)
}

// FindEndOfCurrentWordWithSpace is almost the same as FindEndOfCurrentWord.
// The only difference is to ignore contiguous spaces.
func (d *Document) FindEndOfCurrentWordWithSpace() int {
	x := d.TextAfterCursor()

	start := istrings.IndexNotByte(x, ' ')
	if start == -1 {
		return len(x)
	}

	end := strings.IndexByte(x[start:], ' ')
	if end == -1 {
		return len(x)
	}

	return start + end
}

// FindEndOfCurrentWordUntilSeparator is almost the same as FindEndOfCurrentWord.
// But this can specify Separator. Return 0 if nothing was found.
func (d *Document) FindEndOfCurrentWordUntilSeparator(sep string) int {
	if sep == "" {
		return d.FindEndOfCurrentWord()
	}

	x := d.TextAfterCursor()
	i := strings.IndexAny(x, sep)
	if i != -1 {
		return i
	}
	return len(x)
}

// FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor is almost the same as FindEndOfCurrentWordWithSpace.
// But this can specify Separator. Return 0 if nothing was found.
func (d *Document) FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor(sep string) int {
	if sep == "" {
		return d.FindEndOfCurrentWordWithSpace()
	}

	x := d.TextAfterCursor()

	start := istrings.IndexNotAny(x, sep)
	if start == -1 {
		return len(x)
	}

	end := strings.IndexAny(x[start:], sep)
	if end == -1 {
		return len(x)
	}

	return start + end
}

// CurrentLineBeforeCursor returns the text from the start of the line until the cursor.
func (d *Document) CurrentLineBeforeCursor() string {
	before := d.TextBeforeCursor()
	lf := strings.LastIndex(before, "\n")
	if lf != -1 {
		return before[lf+1:]
	}
	return before
}

// CurrentLineAfterCursor returns the text from the cursor until the end of the line.
func (d *Document) CurrentLineAfterCursor() string {
	after := d.TextAfterCursor()
	lf := strings.Index(after, "\n")
	if lf != -1 {
		return after[:lf]
	}
	return after
}

// CurrentLine return the text on the line where the cursor is. (when the input
// consists of just one line, it equals `text`.
func (d *Document) CurrentLine() string {
	return d.CurrentLineBeforeCursor() + d.CurrentLineAfterCursor()
}

// Array pointing to the start indexes of all the lines.
func (d *Document) lineStartIndexes() []int {
	if d.startIndexCache == nil {
		lc := d.LineCount()
		lengths := make([]int, lc)
		for i, l := range d.Lines() {
			lengths[i] = len(l)
		}

		// Calculate cumulative sums.
		indexes := make([]int, lc+1)
		indexes[0] = 0 // https://github.com/jonathanslenders/python-prompt-toolkit/blob/master/prompt_toolkit/document.py#L189
		pos := 0
		for i, l := range lengths {
			pos += l + 1
			indexes[i+1] = pos
		}
		if lc > 1 {
			// Pop the last item. (This is not a new line.)
			indexes = indexes[:lc]
		}
		d.startIndexCache = indexes
	}

	return d.startIndexCache
}

// For the index of a character at a certain line, calculate the index of
// the first character on that line.
func (d *Document) findLineStartIndex(index int) (pos int, lineStartIndex int) {
	indexes := d.lineStartIndexes()
	pos = bisect.Right(indexes, index) - 1
	lineStartIndex = indexes[pos]
	return
}

// CursorPositionRow returns the current row. (0-based.)
func (d *Document) CursorPositionRow() (row int) {
	row, _ = d.findLineStartIndex(d.cursorPosition)
	return
}

// CursorPositionCol returns the current column. (0-based.)
func (d *Document) CursorPositionCol() (col int) {
	// Don't use self.text_before_cursor to calculate this. Creating substrings
	// and splitting is too expensive for getting the cursor position.
	_, index := d.findLineStartIndex(d.cursorPosition)
	col = d.cursorPosition - index
	return
}

// GetCursorLeftPosition returns the relative position for cursor left.
func (d *Document) GetCursorLeftPosition(count int) int {
	if count < 0 {
		return d.GetCursorRightPosition(-count)
	}
	if d.CursorPositionCol() > count {
		return -count
	}
	return -d.CursorPositionCol()
}

// GetCursorRightPosition returns relative position for cursor right.
func (d *Document) GetCursorRightPosition(count int) int {
	if count < 0 {
		return d.GetCursorLeftPosition(-count)
	}
	if len(d.CurrentLineAfterCursor()) > count {
		return count
	}
	return len(d.CurrentLineAfterCursor())
}

// GetCursorUpPosition return the relative cursor position (character index) where we would be
// if the user pressed the arrow-up button.
func (d *Document) GetCursorUpPosition(count int, preferredColumn int) int {
	var col int
	if preferredColumn == -1 { // -1 means nil
		col = d.CursorPositionCol()
	} else {
		col = preferredColumn
	}

	row := d.CursorPositionRow() - count
	if row < 0 {
		row = 0
	}
	return d.TranslateRowColToIndex(row, col) - d.cursorPosition
}

// GetCursorDownPosition return the relative cursor position (character index) where we would be if the
// user pressed the arrow-down button.
func (d *Document) GetCursorDownPosition(count int, preferredColumn int) int {
	var col int
	if preferredColumn == -1 { // -1 means nil
		col = d.CursorPositionCol()
	} else {
		col = preferredColumn
	}
	row := d.CursorPositionRow() + count
	return d.TranslateRowColToIndex(row, col) - d.cursorPosition
}

// Lines returns the array of all the lines.
func (d *Document) Lines() []string {
	if d.linesCache == nil {
		d.linesCache = strings.Split(d.Text, "\n")
	}
	return d.linesCache
}

// LineCount return the number of lines in this document. If the document ends
// with a trailing \n, that counts as the beginning of a new line.
func (d *Document) LineCount() int {
	if d.linesCache == nil {
		return strings.Count(d.Text, "\n") + 1
	}
	return len(d.linesCache)
}

// TranslateIndexToPosition given an index for the text, return the corresponding (row, col) tuple.
// (0-based. Returns (0, 0) for index=0.)
func (d *Document) TranslateIndexToPosition(index int) (row int, col int) {
	row, rowIndex := d.findLineStartIndex(index)
	col = index - rowIndex
	return
}

// TranslateRowColToIndex given a (row, col), return the corresponding index.
// (Row and col params are 0-based.)
func (d *Document) TranslateRowColToIndex(row int, column int) (index int) {
	indexes := d.lineStartIndexes()
	if row < 0 {
		row = 0
	} else if row > len(indexes) {
		row = len(indexes) - 1
	}
	index = indexes[row]
	line := d.Lines()[row]

	// python) result += max(0, min(col, len(line)))
	if column > 0 || len(line) > 0 {
		if column > len(line) {
			index += len(line)
		} else {
			index += column
		}
	}

	// Keep in range. (len(self.text) is included, because the cursor can be
	// right after the end of the text as well.)
	// python) result = max(0, min(result, len(self.text)))
	if index > len(d.Text) {
		index = len(d.Text)
	}
	if index < 0 {
		index = 0
	}
	return index
}

// CursorOnLastLine returns true when we are at the last line.
func (d *Document) CursorOnLastLine() bool {
	return d.CursorPositionRow() == (d.LineCount() - 1)
}

func (d *Document) CursorAtEndOfLine() bool {
	return len(d.CurrentLineAfterCursor()) == 0
}

// GetEndOfLinePosition returns relative position for the end of this line.
func (d *Document) GetEndOfLinePosition() int {
	return len([]rune(d.CurrentLineAfterCursor()))
}

func (d *Document) leadingWhitespaceInCurrentLine() (margin string) {
	trimmed := strings.TrimSpace(d.CurrentLine())
	margin = d.CurrentLine()[:len(d.CurrentLine())-len(trimmed)]
	return
}
