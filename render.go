package prompt

import (
	"fmt"
	"os"
	"runtime"

	"github.com/c-bata/go-prompt/internal/debug"
	runewidth "github.com/mattn/go-runewidth"
)

// Render to render prompt information from state of Buffer.
type Render struct {
	out                ConsoleWriter
	prefix             string
	livePrefixCallback func() (prefix string, useLivePrefix bool)
	title              string
	termHeight         uint16
	termWidth          uint16

	previousCursor    Coord
	previousLineCount int

	// colors,
	Colors RenderColors
}

type RenderColors struct {
	prefixText              Color
	prefixBG                Color
	inputText               Color
	inputBG                 Color
	previewSuggestionText   Color
	previewSuggestionBG     Color
	suggestionText          Color
	suggestionBG            Color
	selectedSuggestionText  Color
	selectedSuggestionBG    Color
	descriptionText         Color
	descriptionBG           Color
	selectedDescriptionText Color
	selectedDescriptionBG   Color
	scrollbarThumb          Color
	scrollbarBG             Color
}

var defaultColors = RenderColors{
	prefixText:              DefaultColor,
	previewSuggestionText:   White,
	previewSuggestionBG:     DefaultColor,
	suggestionText:          Black,
	suggestionBG:            Gray,
	selectedSuggestionText:  White,
	selectedSuggestionBG:    Blue,
	descriptionText:         Black,
	descriptionBG:           Yellow,
	selectedDescriptionText: Gray,
	selectedDescriptionBG:   Black,
	scrollbarThumb:          Black,
	scrollbarBG:             Gray,
}

func NewRender(prefix string, w ConsoleWriter) *Render {
	return &Render{
		prefix: prefix,
		out:    w,
		Colors: defaultColors,

		previousLineCount: 1,

		livePrefixCallback: func() (string, bool) { return "", false },
	}
}

// Setup to initialize console output.
func (r *Render) Setup() {
	if r.title != "" {
		r.out.SetTitle(r.title)
		debug.AssertNoError(r.out.Flush())
	}
}

// getCurrentPrefix to get current prefix.
// If live-prefix is enabled, return live-prefix.
func (r *Render) getCurrentPrefix() string {
	if prefix, ok := r.livePrefixCallback(); ok {
		return prefix
	}
	return r.prefix
}

func (r *Render) renderPrefix() {
	r.out.SetColor(r.Colors.prefixText, r.Colors.prefixBG, false)
	r.out.WriteStr(r.getCurrentPrefix())
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// TearDown to clear title and erasing.
func (r *Render) TearDown() {
	r.out.ClearTitle()
	r.out.EraseDown()
	debug.AssertNoError(r.out.Flush())
}

func (r *Render) prepareArea(lines int) {
	for i := 0; i < lines; i++ {
		r.out.ScrollDown()
	}
	for i := 0; i < lines; i++ {
		r.out.ScrollUp()
	}
	return
}

// UpdateWinSize called when window size is changed.
func (r *Render) UpdateWinSize(ws *WinSize) {
	r.termHeight = ws.Row
	r.termWidth = ws.Col
	return
}

func (r *Render) renderWindowTooSmall() {
	r.out.CursorGoTo(0, 0)
	r.out.EraseScreen()
	r.out.SetColor(Red, White, false)
	r.out.WriteStr("Your console window is too small...")
	return
}

func (r *Render) renderCompletion(buf *Buffer, completions *CompletionManager) {
	suggestions := completions.GetSuggestions()
	if len(completions.GetSuggestions()) == 0 {
		return
	}
	prefix := r.getCurrentPrefix()
	formatted, width := formatSuggestions(
		suggestions,
		int(r.termWidth)-runewidth.StringWidth(prefix)-1, // -1 means a width of scrollbar
	)
	// +1 means a width of scrollbar.
	width++

	windowHeight := len(formatted)
	if windowHeight > int(completions.max) {
		windowHeight = int(completions.max)
	}
	formatted = formatted[completions.verticalScroll : completions.verticalScroll+windowHeight]
	r.prepareArea(windowHeight)

	cursor := runewidth.StringWidth(prefix) + runewidth.StringWidth(buf.Document().TextBeforeCursor())
	x, _ := r.toPos(cursor)
	if x+width >= int(r.termWidth) {
		cursor = r.backward(cursor, x+width-int(r.termWidth))
	}

	contentHeight := len(completions.tmp)

	fractionVisible := float64(windowHeight) / float64(contentHeight)
	fractionAbove := float64(completions.verticalScroll) / float64(contentHeight)

	scrollbarHeight := int(clamp(float64(windowHeight), 1, float64(windowHeight)*fractionVisible))
	scrollbarTop := int(float64(windowHeight) * fractionAbove)

	isScrollThumb := func(row int) bool {
		return scrollbarTop <= row && row <= scrollbarTop+scrollbarHeight
	}

	selected := completions.selected - completions.verticalScroll
	r.out.SetColor(White, Cyan, false)
	for i := 0; i < windowHeight; i++ {
		r.out.CursorDown(1)
		if i == selected {
			r.out.SetColor(r.Colors.selectedSuggestionText, r.Colors.selectedSuggestionBG, true)
		} else {
			r.out.SetColor(r.Colors.suggestionText, r.Colors.suggestionBG, false)
		}
		r.out.WriteStr(formatted[i].Text)

		if i == selected {
			r.out.SetColor(r.Colors.selectedDescriptionText, r.Colors.selectedDescriptionBG, false)
		} else {
			r.out.SetColor(r.Colors.descriptionText, r.Colors.descriptionBG, false)
		}
		r.out.WriteStr(formatted[i].Description)

		if isScrollThumb(i) {
			r.out.SetColor(DefaultColor, r.Colors.scrollbarThumb, false)
		} else {
			r.out.SetColor(DefaultColor, r.Colors.scrollbarBG, false)
		}
		r.out.WriteStr(" ")
		r.out.SetColor(DefaultColor, DefaultColor, false)

		r.lineWrap(cursor + width)
		r.backward(cursor+width, width)
	}

	if x+width >= int(r.termWidth) {
		r.out.CursorForward(x + width - int(r.termWidth))
	}

	r.out.CursorUp(windowHeight)
	r.out.SetColor(DefaultColor, DefaultColor, false)
	return
}

func dbg(m string, args ...interface{}) {
	fmt.Fprint(os.Stderr, "\x1b[33;1m")
	fmt.Fprintf(os.Stderr, m, args...)
	fmt.Fprintln(os.Stderr, "\x1b[m")
}

// Render renders to the console.
func (r *Render) Render(buffer *Buffer, completion *CompletionManager) {
	// In situations where a pseudo tty is allocated (e.g. within a docker container),
	// window size via TIOCGWINSZ is not immediately available and will result in 0,0 dimensions.
	if r.termWidth == 0 {
		return
	}

	doc := buffer.Document()

	// TODO: this should render into an off-screen 'buffer'.
	//   this buffer would then be compared with the buffer rendered previously
	//   and generate actual output instructions from that.

	defer func() { debug.AssertNoError(r.out.Flush()) }()
	//fmt.Fprintln(os.Stderr, "------------------------- RENDER")

	// if lines have been added to the edit, add space
	lcount := doc.LineCount()
	added := 0
	if lcount > r.previousLineCount {
		r.out.WriteRaw([]byte{'\n'})
		added = 1
		dbg("added LF  (%d -> %d)", r.previousLineCount, doc.LineCount())
	}
	// move to beginning of the current prompt

	r.promptHome(Coord{r.previousCursor.X, r.previousCursor.Y + added})

	line := buffer.Text()
	prefix := r.getCurrentPrefix()
	// calculate future cursor position after prefix & line is printed
	editPoint := doc.DisplayCursorCoordWithPrefix(int(r.termWidth), prefix)
	dbg("editPoint @ %+v", editPoint)

	// prepare area
	y := lcount
	//_, y := r.toPos(cursor)

	h := y + 1 + int(completion.max)
	if h > int(r.termHeight) || completionMargin > int(r.termWidth) {
		r.renderWindowTooSmall()
		// TODO: do some better fallback  (this will just spam-loop)
		return
	}

	// Rendering
	r.out.HideCursor()
	defer r.out.ShowCursor()

	r.out.SaveCursor()

	// render the complete prompt; prefix and editor content
	r.renderPrefix()
	r.out.SetColor(r.Colors.inputText, r.Colors.inputBG, false)
	r.out.WriteStr(line)
	r.out.SetColor(DefaultColor, DefaultColor, false)
	//r.lineWrap(cursor)
	r.out.EraseDown()

	// position the cursor at the edit point after the editor rendering
	r.out.RestoreCursor()
	r.moveRel(Coord{0, 0}, editPoint)
	//cursor = r.backward(cursor, runewidth.StringWidth(line)-buffer.DisplayCursorPosition())

	r.renderCompletion(buffer, completion)

	// if a completion suggestion is currently selected update the screen -- but NOT the editor content!
	if suggest, ok := completion.GetSelectedSuggestion(); ok {
		// move to the beginning of the word being completed
		completing_word := doc.GetWordBeforeCursorUntilSeparator(completion.wordSeparator)
		editPoint = r.moveRel(editPoint, Coord{-runewidth.StringWidth(completing_word), 0})
		//cursor := r.backward(cursor, runewidth.StringWidth(buffer.Document().GetWordBeforeCursorUntilSeparator(completion.wordSeparator)))

		// write the suggestion, using the configured preview style
		r.out.SetColor(r.Colors.previewSuggestionText, r.Colors.previewSuggestionBG, false)
		r.out.WriteStr(suggest.Text)
		//cursor += runewidth.StringWidth(suggest.Text)
		// move edit point to the end of the suggested word
		editPoint.X += runewidth.StringWidth(suggest.Text)
		r.out.SaveCursor()

		// write the text following the cursor (using default style)
		r.out.SetColor(DefaultColor, DefaultColor, false)
		rest := buffer.Document().TextAfterCursor()
		r.out.WriteStr(rest)
		// total length of line
		eol := editPoint.X + runewidth.StringWidth(rest)
		// move cursor back to the edit point
		//r.backward(cursor, runewidth.StringWidth(rest))
		if r.lineWrap(eol) { // output LF if necessary
			r.out.RestoreCursor()
			r.out.CursorUp(1)
		} else {
			r.out.RestoreCursor()
		}
	}
	r.previousCursor = editPoint
	r.previousLineCount = lcount
}

// BreakLine to break line.
func (r *Render) BreakLine(buffer *Buffer) {
	// Erasing and Render
	cursor := runewidth.StringWidth(buffer.Document().TextBeforeCursor()) + runewidth.StringWidth(r.getCurrentPrefix())
	r.clear(cursor)
	r.renderPrefix()
	r.out.SetColor(r.Colors.inputText, r.Colors.inputBG, false)
	r.out.WriteStr(buffer.Document().Text + "\n")
	r.out.SetColor(DefaultColor, DefaultColor, false)
	debug.AssertNoError(r.out.Flush())

	r.previousCursor = Coord{}
	r.previousLineCount = 1
}

// clear erases the screen from a beginning of input
// even if there is line break which means input length exceeds a window's width.
func (r *Render) clear(cursor int) {
	r.move(cursor, 0)
	r.out.EraseDown()
}

// backward moves cursor to backward from a current cursor position
// regardless there is a line break.
func (r *Render) backward(from, n int) int {
	return r.move(from, from-n)
}

// move moves cursor to specified position from the beginning of input
// even if there is a line break.
func (r *Render) move(from, to int) int {
	fromX, fromY := r.toPos(from)
	toX, toY := r.toPos(to)

	dbg("move: left: %d  up: %d", fromX-toX, fromY-toY)

	r.out.CursorUp(fromY - toY)
	r.out.CursorBackward(fromX - toX)
	return to
}

// moveRel moves the cursor 'from' coord in the 'rel' direction (right & down).
//   if 'rel' values are negative it moves in the oppositve direction
// returns 'from' + 'rel'
func (r *Render) moveRel(from, rel Coord) Coord {
	dbg("moveCoord: %+v + %+v", from, rel)
	r.out.CursorDown(rel.Y)
	r.out.CursorForward(rel.X)
	return coord_add(from, rel)
}

func (r *Render) promptHome(from Coord) {
	dbg("promptHome: %+v", from)
	r.out.CursorUp(from.Y)
	r.out.CursorBackward(from.X)
}

// coord_add returns a + b
func coord_add(a, b Coord) Coord {
	return Coord{a.X + b.X, a.Y + b.Y}
}

// coord_sub returns a - b
func coord_sub(a, b Coord) Coord {
	return Coord{a.X - b.X, a.Y - b.Y}
}

// toPos returns the relative position from the beginning of the string.
func (r *Render) toPos(cursor int) (x, y int) {
	col := int(r.termWidth)
	return cursor % col, cursor / col
}

func (r *Render) lineWrap(cursor int) bool {
	if runtime.GOOS != "windows" && cursor > 0 && cursor%int(r.termWidth) == 0 {
		r.out.WriteRaw([]byte{'\n'})
		return true
	}
	return false
}

func clamp(high, low, x float64) float64 {
	switch {
	case high < x:
		return high
	case x < low:
		return low
	default:
		return x
	}
}
