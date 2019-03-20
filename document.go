package prompt

import (
	"strings"
	"unicode"

	runewidth "github.com/mattn/go-runewidth"

	"github.com/tatsujin/go-prompt/internal/bisect"
	"github.com/tatsujin/go-prompt/internal/runes"
)

type Column int  // terminal column when rendered
type Row int     // index of a terminal or multi-line text row
type Index = int // absolute character index (into []rune) (a type alias because it's used with bisect)
type Offset int  // relative character offset between two "Index" values

// Document is a read-only view of the current editor content
type Document struct {
	//text  string
	text []rune // the text as a rune slice

	// 'cursor' is an index into 'text';
	// if 'text' is "日本(cursor)語", 'cursor' is 2.
	// But DisplayedCursorPosition returns 4 because '日' and '本' are double width characters.
	cursor Index

	linesCache      [][]rune
	startIndexCache []Index
}

// NewDocument return the new empty document.
func NewDocument(text string, cpos Index) *Document {
	return &Document{
		//text:   text,
		text:   []rune(text),
		cursor: cpos,
	}
}

func (d *Document) Text() string {
	return string(d.text)
}

func IsWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// CursorTextColumn returns the column at which the cursor would be if the text was rendered in a terminal (with infinite width).
// So if Document is "日本(cursor)語", CursorColumn returns 4 because '日' and '本' are double width characters.
// This does not handle multi-line text.
func (d *Document) CursorTextColumn() Column {
	return Column(runewidth.StringWidth(d.TextBeforeCursor()))
}

// CursorDisplayCoord returns similar to CursorColumn, but both col & row are returned.
// It is assumed that the text starts at column 0.
func (d *Document) CursorDisplayCoord(termWidth Column) Coord {
	return d.CursorDisplayCoordWithPrefix(termWidth, nil)
}

// CursorDisplayCoordWithPrefix same as CursorCoord but with prefix(es) taken into account.
// It is assumed that the text starts at column 0.
func (d *Document) CursorDisplayCoordWithPrefix(termWidth Column, prefix func(doc *Document, row Row) string) Coord {

	// TODO: this assumes that strings returned by prefix() does not contain '\n'

	var x Column
	var y Row
	var cpos Index
	var end int
	done := false
	for row, rtext := range d.lines() {
		if end = Index(len(rtext)); cpos+end >= d.cursor {
			end = d.cursor - cpos
			done = true // this is the last line we'll process
		}
		var w Column
		if prefix != nil {
			w = Column(runewidth.StringWidth(prefix(d, Row(row))))
		}
		w += Column(runewidth.StringWidth(string(rtext[:end])))

		x = w % termWidth
		y += 1 + Row(w/termWidth)

		advance := Index(Offset(len(rtext)) + LFsize)
		cpos += advance

		if done {
			break
		}
	}
	// -1: we want 'index', not 'number-of-lines'
	return Coord{x, y - 1}
}

// GetCharRelativeToCursor return character relative to cursor position (0 = at cursor), or empty string
func (d *Document) GetCharFromCursor(offset Offset) (r rune) {
	if d.cursor+Index(offset) >= len(d.text) {
		return 0
	}
	return d.text[d.cursor+Index(offset)]
}

// TextBeforeCursor returns the text before the cursor.
//   includes all preceeding lines.
func (d *Document) TextBeforeCursor() string {
	return string(d.textBeforeCursor())
}

func (d *Document) textBeforeCursor() []rune {
	return d.text[:d.cursor]
}

// TextAfterCursor returns the text after the cursor.
//   includes all following lines.
func (d *Document) TextAfterCursor() string {
	return string(d.textAfterCursor())
}
func (d *Document) textAfterCursor() []rune {
	return d.text[d.cursor:]
}

// GetWordBeforeCursor returns the word(part) before the cursor.
// If we have whitespace before the cursor this returns an empty string.
func (d *Document) GetWordBeforeCursor() string {
	start := d.FindStartOfCurrentWord()
	return string(d.textBeforeCursor()[start:])
}

// GetWordAfterCursor returns the word(part) after the cursor.
// If we have whitespace after the cursor this returns an empty string.
func (d *Document) GetWordAfterCursor() string {
	end := d.FindEndOfCurrentWord()
	return string(d.textAfterCursor()[:end])
}

// GetWordBeforeCursorWithSpace returns the word(part) before the cursor.
// Unlike GetWordBeforeCursor, it returns string containing space
func (d *Document) GetWordBeforeCursorWithSpace() string {
	start := d.FindStartOfCurrentWordWithSpace()
	return string(d.textBeforeCursor()[start:])
}

// GetWordAfterCursorWithSpace returns the word(part) after the cursor.
// Unlike GetWordAfterCursor, it returns string containing space
func (d *Document) GetWordAfterCursorWithSpace() string {
	end := d.FindEndOfCurrentWordWithSpace()
	return string(d.textAfterCursor()[:end])
}

// GetWordBeforeCursorUntilSeparator returns the text before the cursor until next separator.
func (d *Document) GetWordBeforeCursorUntilSeparator(sep string) string {
	start := d.FindStartOfCurrentWordUntilSeparator(sep)
	return string(d.textBeforeCursor()[start:])
}

// GetWordAfterCursorUntilSeparator returns the text after the cursor until next separator.
func (d *Document) GetWordAfterCursorUntilSeparator(sep string) string {
	end := d.FindEndOfCurrentWordUntilSeparator(sep)
	return string(d.textAfterCursor()[:end])
}

// GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor returns the word before the cursor.
// Unlike GetWordBeforeCursor, it returns string containing space
func (d *Document) GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor(sep string) string {
	start := d.FindStartOfCurrentWordUntilSeparatorIgnoreNextToCursor(sep)
	return string(d.textBeforeCursor()[start:])
}

// GetWordAfterCursorUntilSeparatorIgnoreNextToCursor returns the word after the cursor.
// Unlike GetWordAfterCursor, it returns string containing space
func (d *Document) GetWordAfterCursorUntilSeparatorIgnoreNextToCursor(sep string) string {
	end := d.FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor(sep)
	return string(d.textAfterCursor()[:end])
}

// FindStartOfCurrentWord returns an index to the start of the word
// the cursor is currently at. Return 0 if nothing was found.
func (d *Document) FindStartOfCurrentWord() Index {
	before := d.textBeforeCursor()

	if idx := runes.LastIndexRune(before, ' '); idx == -1 {
		//fmt.Fprintf(os.Stderr, "%q start: %d\n", string(before), idx)
		return 0
	} else {
		//fmt.Fprintf(os.Stderr, "%q start: %d\n", string(before), idx)
		return Index(idx + 1)
	}
}

// FindStartOfCurrentWordWithSpace is similar to FindStartOfCurrentWord,
// but it ignores contiguous spaces.
func (d *Document) FindStartOfCurrentWordWithSpace() Index {
	before := d.textBeforeCursor()

	if end := runes.LastIndexNotRune(before, ' '); end == -1 {
		return 0
	} else if start := runes.LastIndexRune(before[:end], ' '); start == -1 {
		return 0
	} else {
		return Index(start + 1)
	}
}

// FindStartOfCurrentWordUntilSeparator is similar to FindStartOfCurrentWord.
// But this can specify Separator. Return 0 if nothing was found.
func (d *Document) FindStartOfCurrentWordUntilSeparator(sep string) Index {
	if sep == "" {
		return d.FindStartOfCurrentWord()
	}

	before := d.textBeforeCursor()

	if idx := runes.LastIndexAny(before, []rune(sep)); idx != -1 {
		return Index(idx + 1)
	}
	return 0
}

// FindStartOfCurrentWordUntilSeparatorIgnoreNextToCursor is similar to FindStartOfCurrentWordWithSpace.
// But this can specify Separator. Return 0 if nothing was found.
func (d *Document) FindStartOfCurrentWordUntilSeparatorIgnoreNextToCursor(sep string) Index {
	if sep == "" {
		return d.FindStartOfCurrentWordWithSpace()
	}

	before := d.textBeforeCursor()
	rsep := []rune(sep)

	if end := runes.LastIndexNotAny(before, rsep); end == -1 {
		return 0
	} else if start := runes.LastIndexAny(before[:end], rsep); start == -1 {
		return 0
	} else {
		return Index(start + 1)
	}
}

// FindEndOfCurrentWord returns am offset from the cursor to the end of the current word.
// Return 0 if nothing was found.
func (d *Document) FindEndOfCurrentWord() Offset {
	after := d.textAfterCursor()

	if idx := runes.IndexRune(after, ' '); idx == -1 {
		return Offset(len(after))
	} else {
		return Offset(idx)
	}
}

// FindEndOfCurrentWordWithSpace is similar to FindEndOfCurrentWord.
// The only difference is to ignore contiguous spaces.
func (d *Document) FindEndOfCurrentWordWithSpace() Offset {
	after := d.textAfterCursor()

	if start := runes.IndexNotRune(after, ' '); start == -1 {
		return Offset(len(after))
	} else if end := runes.IndexRune(after[start:], ' '); end == -1 {
		return Offset(len(after))
	} else {
		return Offset(start + end)
	}
}

// FindEndOfCurrentWordUntilSeparator is similar to FindEndOfCurrentWord.
// But this can specify Separator. Return 0 if nothing was found.
func (d *Document) FindEndOfCurrentWordUntilSeparator(sep string) Offset {
	if sep == "" {
		return d.FindEndOfCurrentWord()
	}

	after := d.textAfterCursor()

	if idx := runes.IndexAny(after, []rune(sep)); idx == -1 {
		return Offset(len(after))
	} else {
		return Offset(idx)
	}
}

// FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor is similar to FindEndOfCurrentWordWithSpace.
// But this can specify Separator. Return 0 if nothing was found.
func (d *Document) FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor(sep string) Offset {
	if sep == "" {
		return d.FindEndOfCurrentWordWithSpace()
	}

	after := d.textAfterCursor()

	rsep := []rune(sep)

	if start := runes.IndexNotAny(after, rsep); start == -1 {
		return Offset(len(after))
	} else if end := runes.IndexAny(after[start:], rsep); end == -1 {
		return Offset(len(after))
	} else {
		return Offset(start + end)
	}
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
	after := d.textAfterCursor()

	if lf := runes.IndexRune(after, '\n'); lf == -1 {
		return string(after)
	} else {
		return string(after[:lf])
	}
}

// CurrentLine return the text of the line the cursor is on. (when the input
// consists of just one line, it equals `text`.
func (d *Document) CurrentLine() string {
	//return d.CurrentLineBeforeCursor() + d.CurrentLineAfterCursor()
	// TODO: is this faster?
	return string(d.lines()[d.CursorRow()])
}

var LFsize = Offset(len([]rune("\n")))

// Array with byte indexes to the start of all the lines.
func (d *Document) lineStartIndexes() []Index {
	if d.startIndexCache == nil {
		lcount := d.LineCount()
		lengths := make([]int, lcount)
		for i, line := range d.lines() {
			lengths[i] = len(line)
		}

		// Calculate cumulative sums.
		indexes := make([]Index, lcount+1)
		indexes[0] = 0 // https://github.com/jonathanslenders/python-prompt-toolkit/blob/master/prompt_toolkit/document.py#L189
		var pos Offset
		for i, ln := range lengths {
			pos += Offset(ln) + LFsize
			indexes[i+1] = Index(pos)
		}
		if lcount > 1 {
			// Pop the last item. (This is not a new line.)
			indexes = indexes[:lcount]
		}
		d.startIndexCache = indexes
	}

	return d.startIndexCache
}

// For the index of a character, return row index and the index of the first character on that line.
func (d *Document) findLineStartIndex(index Index) (row Row, lineStartIndex Index) {
	indexes := d.lineStartIndexes()
	row = Row(bisect.Right(indexes, int(index)) - 1)
	lineStartIndex = indexes[row]
	return
}

// CursorIndex returns the absolute character index of the cursor's position (including line feeds, etc).
func (d *Document) CursorIndex() (index Index) {
	return d.cursor
}

// CursorRow returns the row index to which the cursor is at (multi-line editing).
func (d *Document) CursorRow() (row Row) {
	r, _ := d.findLineStartIndex(d.cursor)
	return r
}

// CursorColumnIndex returns the row-relative character index to the character at the cursor in the row the cursor is on.
func (d *Document) CursorColumnIndex() (col Index) {
	// Don't use self.text_before_cursor to calculate this. Creating substrings
	// and splitting is too expensive for getting the cursor position.

	// NOTE: hm, but 'lineStartIndexes' does string splitting and creates substrings?

	_, index := d.findLineStartIndex(d.cursor)
	return d.cursor - index
}

// GetCursorLeftOffset returns the relative position for moving the cursor left.
func (d *Document) GetCursorLeftOffset(off Offset) Offset {
	if off < 0 {
		return d.GetCursorRightOffset(-off)
	}
	cursorCol := Offset(d.CursorColumnIndex())
	if cursorCol > off {
		return -off
	}
	return -cursorCol
}

// GetCursorRightOffset returns relative position for moving the cursor right.
func (d *Document) GetCursorRightOffset(off Offset) Offset {
	if off < 0 {
		return d.GetCursorLeftOffset(-off)
	}
	after := []rune(d.CurrentLineAfterCursor())
	if Offset(len(after)) > off {
		return off
	}
	return Offset(len(after))
}

// GetCursorUpOffset return the relative cursor offset needed to move the cursor
// if the user pressed the up-arrow key.
func (d *Document) GetCursorUpOffset(off Row, preferred Index) Offset {
	if off < 0 {
		return d.GetCursorDownOffset(-off, preferred)
	}
	if preferred == -1 { // use current
		preferred = d.CursorColumnIndex()
	}
	row := d.CursorRow() - off
	return Offset(d.TranslateRowColToIndex(row, preferred) - d.cursor)
}

// GetCursorDownOffset return the relative cursor offset needed to move the cursor
// if the user pressed the down-arrow key.
func (d *Document) GetCursorDownOffset(off Row, preferred Index) Offset {
	if off < 0 {
		return d.GetCursorUpOffset(-off, preferred)
	}
	if preferred == -1 { // use current
		preferred = d.CursorColumnIndex()
	}
	row := d.CursorRow() + off
	return Offset(d.TranslateRowColToIndex(row, preferred) - d.cursor)
}

// Lines returns the array of all the lines.
func (d *Document) Lines() []string {
	// TODO: separate cache for string-typed lines ?
	return strings.Split(string(d.text), "\n")
}

func (d *Document) lines() [][]rune {
	if d.linesCache == nil {
		d.linesCache = runes.SplitRune(d.text, '\n')
	}
	return d.linesCache
}

// LineCount return the number of lines in this document. If the document ends
// with a trailing \n, that counts as the beginning of a new line.
func (d *Document) LineCount() int {
	return len(d.lines())
}

// TranslateIndexToOffset given an index for the text, return the corresponding (row, col) tuple.
// Returns (0, 0) for index=0.
func (d *Document) TranslateIndexToRowCol(index Index) (row Row, col Index) {
	row, rowIndex := d.findLineStartIndex(index)
	col = index - rowIndex
	return
}

// TranslateRowColToIndex given a (row, col), return the corresponding absolute index.
func (d *Document) TranslateRowColToIndex(row Row, column Index) (index Index) {
	indexes := d.lineStartIndexes()
	if row < 0 {
		row = 0
	} else if row > Row(len(indexes)) {
		row = Row(len(indexes)) - 1
	}
	index = indexes[row]
	line := d.lines()[row]

	// python: index += max(0, min(col, len(line)))
	if column > 0 || len(line) > 0 {
		if column > len(line) {
			index += len(line)
		} else {
			index += column
		}
	}

	// Keep in range. (len(self.text) is included, because the cursor can be
	// right after the end of the text as well.)
	// python: result = max(0, min(result, len(self.text)))
	if index > len(d.text) {
		index = len(d.text)
	} else if index < 0 {
		index = 0
	}
	return index
}

// CursorOnLastLine returns true when we are at the last line.
func (d *Document) CursorOnLastLine() bool {
	return d.CursorRow() == Row(d.LineCount()-1)
}

func (d *Document) CursorAtEndOfLine() bool {
	return len(d.CurrentLineAfterCursor()) == 0
}

// GetEndOfLineColumn returns relative character offset to the end of the current line.
func (d *Document) GetEndOfLineOffset() Offset {
	return Offset(len([]rune(d.CurrentLineAfterCursor())))
}

func (d *Document) leadingWhitespaceInCurrentLine() (margin string) {
	trimmed := strings.TrimSpace(d.CurrentLine())
	return d.CurrentLine()[:len(d.CurrentLine())-len(trimmed)]
}
