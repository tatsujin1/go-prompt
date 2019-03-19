package prompt

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/c-bata/go-prompt/internal/debug"
	runewidth "github.com/mattn/go-runewidth"
)

// Render to render prompt information from state of Buffer.
type Render struct {
	out ConsoleWriter
	//cursor             Cursor
	prefix             string
	prefixCallback     func(doc *Document, row Row) (prefix string, usePrefix bool)
	continuationPrefix string
	suffix             string
	suffixCallback     func(doc *Document, row Row) (prefix string, usePrefix bool)
	title              string
	termHeight         Row
	termWidth          Column

	previousCursor    Coord
	previousLineCount int

	Colors             RenderColors
	trueColorSupported bool

	outputLock *sync.Mutex
}

type RenderColors struct {
	prefixText              Color
	prefixBG                Color
	inputText               Color
	inputBG                 Color
	choiceText              Color
	choiceBG                Color
	descriptionText         Color
	descriptionBG           Color
	selectedChoiceText      Color
	selectedChoiceBG        Color
	selectedDescriptionText Color
	selectedDescriptionBG   Color
	previewChoiceText       Color
	previewChoiceBG         Color
	scrollbarThumb          Color
	scrollbarBG             Color
}

// these should only use ANSI colors
// TODO: set unspecified ones to DefaultColor (must use reflect)
//       and remove the check in SetDisplayAttributes
var defaultColors = RenderColors{
	choiceText:              Black,
	choiceBG:                Gray,
	descriptionText:         BrightBlack,
	selectedChoiceText:      White,
	selectedChoiceBG:        Blue,
	selectedDescriptionText: Gray,
	previewChoiceText:       White,
	scrollbarThumb:          BrightBlack,
}

var nilPrefix = func(*Document, Row) (string, bool) { return "", false }

func NewRender(prefix string, w ConsoleWriter) *Render {
	r := &Render{
		prefix: prefix,
		out:    w,
		//cursor: NewCursor(w),
		Colors: defaultColors,

		previousLineCount: 1,

		prefixCallback: nilPrefix,
		suffixCallback: nilPrefix,

		outputLock: &sync.Mutex{},
	}

	// https://gist.github.com/XVilka/8346728#detection
	cterm := os.Getenv("COLORTERM")
	if cterm == "truecolor" || cterm == "24bit" {
		r.trueColorSupported = true
	}

	return r
}

func (r *Render) ValidateColor(c Color) (Color, bool) {
	if r.trueColorSupported || !c.IsTrueColor() {
		return c, true
	}
	// TODO: make fallback color?
	return nil, false
}

// Setup to initialize console output.
func (r *Render) Setup() {
	if r.title != "" {
		r.out.SetTitle(r.title)
		debug.AssertNoError(r.out.Flush())
	}
}

// TearDown to clear title and erasing.
func (r *Render) TearDown() {
	r.outputLock.Lock()
	defer r.outputLock.Unlock()

	r.out.ClearTitle()
	r.out.EraseDown()
	debug.AssertNoError(r.out.Flush())
}

func (r *Render) prepareArea(lines Row) {
	for i := Row(0); i < lines; i++ {
		r.out.ScrollDown()
	}
	for i := Row(0); i < lines; i++ {
		r.out.ScrollUp()
	}
	return
}

// UpdateWinSize called when window size is changed.
func (r *Render) UpdateWinSize(ws *WinSize) {
	r.termHeight = Row(ws.Row)
	r.termWidth = Column(ws.Col)
	return
}

func dbg(m string, args ...interface{}) {
	fmt.Fprint(os.Stderr, "\x1b[33;1m")
	fmt.Fprintf(os.Stderr, m, args...)
	fmt.Fprintln(os.Stderr, "\x1b[m")
}

// Render renders to the console.
func (r *Render) Render(buffer *Buffer, compMgr *CompletionManager) {
	r.outputLock.Lock()
	buffer.RLock()

	r.render(buffer, compMgr)

	buffer.RUnlock()
	r.outputLock.Unlock()
}

func (r *Render) render(buf *Buffer, compMgr *CompletionManager) {
	// In situations where a pseudo tty is allocated (e.g. within a docker container),
	// window size via TIOCGWINSZ is not immediately available and will result in 0,0 dimensions.
	if r.termWidth == 0 {
		return
	}

	doc := buf.Document()

	// TODO: this should render into an off-screen buffer.
	//   this buffer would then be compared with the previously rendered buffer
	//   and generate actual output instructions from that diff.

	defer func() { debug.AssertNoError(r.out.Flush()) }()

	// if lines have been added to the edit, add space
	lcount := doc.LineCount()
	added := 0
	if lcount > r.previousLineCount {
		r.out.WriteRaw([]byte{'\n'})
		added = 1
	}

	// move to beginning of the current prompt
	r.promptHome(Coord{r.previousCursor.X, r.previousCursor.Y + Row(added)})

	// calculate future cursor position after prefix & line is printed
	// TODO: this requires that 'contPfx' has fixed length
	editPoint := doc.CursorDisplayCoordWithPrefix(Column(r.termWidth), r.getPrefix)

	// prepare area
	h := Row(lcount + 1 + int(compMgr.MaxVisibleChoices()))
	if h > r.termHeight || completionMargin > r.termWidth {
		r.renderWindowTooSmall()
		return
	}

	// Rendering
	r.out.HideCursor()
	defer r.out.ShowCursor()

	r.out.SaveCursor()

	// render the complete prompt; prefix and editor content
	r.out.EraseDown()
	r.renderPrompt(doc, false)

	// position the cursor at the edit point after the rendering
	r.out.RestoreCursor()
	r.move(Coord{}, editPoint)

	r.renderCompletion(buf, compMgr)

	// if a completion choice is currently selected, update the screen -- but NOT the editor content!
	if choice, ok := compMgr.Selected(); ok {
		// move to the beginning of the word being completed
		completing_word := doc.GetWordBeforeCursorUntilSeparator(compMgr.wordSeparator)
		editPoint = r.move(editPoint, Coord{Column(-runewidth.StringWidth(completing_word)), 0})

		// write the choice, using the configured preview style
		r.out.SetColor(r.Colors.previewChoiceText, r.Colors.previewChoiceBG, false)
		r.out.WriteStr(choice.Text)
		// move edit point to the end of the suggested word
		editPoint.X += Column(runewidth.StringWidth(choice.Text))
		r.out.SaveCursor()

		// write the text following the cursor (using default style)
		r.out.SetColor(DefaultColor, DefaultColor, false)
		rest := buf.Document().TextAfterCursor()
		r.out.WriteStr(rest)
		// total length of line
		eol := editPoint.X + Column(runewidth.StringWidth(rest))
		// move cursor back to the edit point
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
func (r *Render) BreakLine(buf *Buffer) {
	r.outputLock.Lock()
	defer r.outputLock.Unlock()

	// Erasing and Render
	doc := buf.Document()
	editPoint := doc.CursorDisplayCoordWithPrefix(r.termWidth, r.getPrefix)
	r.promptHome(editPoint)
	r.out.EraseDown()
	r.renderPrompt(doc, true)
	debug.AssertNoError(r.out.Flush())

	r.previousCursor = Coord{}
	r.previousLineCount = 1
}

func (r *Render) OutputAsync(buf *Buffer, compMgr *CompletionManager, format string, a ...interface{}) {
	go func() {
		r.outputLock.Lock()

		r.out.SaveCursor()

		r.promptHome(r.previousCursor)
		r.out.EraseDown()

		text := fmt.Sprintf(format, a...)
		r.out.SetColor(r.Colors.inputText, r.Colors.inputBG, false)
		r.out.WriteRawStr(text)
		// force LF
		outputLines := strings.Count(text, "\n")
		if text[len(text)-1] != '\n' {
			r.out.WriteRawStr("\n")
			outputLines++
		}
		r.out.RestoreCursor()
		r.out.CursorDown(outputLines)

		buf.RLock()
		r.render(buf, compMgr)
		buf.RUnlock()

		r.outputLock.Unlock()
	}()
}

const scrollbarWidth = 1
const safetyMargin = 1

func (r *Render) renderPrompt(doc *Document, breakLine bool) {

	// TODO: syntax highlight of ducment text
	//   porbably make something akin to the "formatted text" in prompt-toolkit

	r.out.SetColor(r.Colors.inputText, r.Colors.inputBG, false)
	for row, line := range strings.SplitAfter(doc.Text(), "\n") {
		r.out.SetColor(r.Colors.prefixText, r.Colors.prefixBG, false)
		r.out.WriteRawStr(r.getPrefix(doc, Row(row)))
		r.out.SetColor(DefaultColor, DefaultColor, false)

		r.out.WriteRawStr(line)
	}

	if breakLine {
		r.out.WriteRawStr("\n")
	}
}

func (r *Render) renderCompletion(buf *Buffer, compMgr *CompletionManager) {
	if compMgr.NumChoices() == 0 {
		return
	}

	editPoint := buf.Document().CursorDisplayCoordWithPrefix(r.termWidth, r.getPrefix)

	widthLimit := r.termWidth - editPoint.X - scrollbarWidth - safetyMargin

	formatted, width, withDesc := compMgr.FormatChoices(widthLimit, r.termWidth)
	width += scrollbarWidth

	windowHeight := Row(len(formatted))
	if windowHeight > Row(compMgr.MaxVisibleChoices()) {
		windowHeight = Row(compMgr.MaxVisibleChoices())
	}

	var cursorMoved Column

	if r.termWidth-editPoint.X < 40 || editPoint.X+width >= r.termWidth {
		cursorMoved = -editPoint.X + 10 // say, at column 10 :)
		r.move(Coord{}, Coord{cursorMoved, 0})
		// re-format the choices, we now have more space
		widthLimit = r.termWidth - (editPoint.X - cursorMoved) - scrollbarWidth - safetyMargin
		formatted, width, withDesc = compMgr.FormatChoices(widthLimit, r.termWidth)
		width += scrollbarWidth
	}

	formatted = formatted[compMgr.verticalScroll : compMgr.verticalScroll+int(windowHeight)]
	r.prepareArea(windowHeight)

	// compute scrollbar parameters
	contentHeight := compMgr.NumChoices()
	fractionVisible := float64(windowHeight) / float64(contentHeight)
	fractionAbove := float64(compMgr.verticalScroll) / float64(contentHeight)

	scrollbarHeight := int(clamp(float64(windowHeight), 1, float64(windowHeight)*fractionVisible))
	scrollbarTop := int(float64(windowHeight) * fractionAbove)

	isScrollThumb := func(row int) bool {
		return scrollbarTop <= row && row <= scrollbarTop+scrollbarHeight
	}

	selected := compMgr.selected - compMgr.verticalScroll

	for i := 0; i < int(windowHeight); i++ {
		r.out.CursorDown(1)

		// draw choice text
		if i == selected {
			r.out.SetColor(r.Colors.selectedChoiceText, r.Colors.selectedChoiceBG, true)
		} else {
			r.out.SetColor(r.Colors.choiceText, r.Colors.choiceBG, false)
		}
		r.out.WriteStr(formatted[i].Text)

		if withDesc { // might be skipped if we don't have space
			// draw choice description
			if i == selected {
				r.out.SetColor(r.Colors.selectedDescriptionText, r.Colors.selectedDescriptionBG, false)
			} else {
				r.out.SetColor(r.Colors.descriptionText, r.Colors.descriptionBG, false)
			}
			r.out.WriteStr(formatted[i].Description)
		}

		if isScrollThumb(i) {
			r.out.SetColor(DefaultColor, r.Colors.scrollbarThumb, false)
		} else {
			r.out.SetColor(DefaultColor, r.Colors.scrollbarBG, false)
		}
		r.out.SetColor(DefaultColor, DefaultColor, false)

		r.move(Coord{}, Coord{-width + 1, 0})
	}

	// move back to edit point (use RestoreCursor?)
	r.move(Coord{}, Coord{-cursorMoved, -windowHeight})

	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// getPrefix to get current prefix.
// If prefix callback is set, use that.
func (r *Render) getPrefix(doc *Document, row Row) string {
	if prefix, ok := r.prefixCallback(doc, row); ok {
		return prefix
	}
	if doc.CursorRow() == 0 {
		return r.prefix
	}
	return r.continuationPrefix
}

// getRightPrefix to get current right prefix.
// If prefix callback is set, use that.
func (r *Render) getSuffix(doc *Document, row Row) string {
	if suffix, ok := r.suffixCallback(doc, row); ok {
		return suffix
	}
	return r.suffix
}

// move moves the cursor in the 'rel' direction (right & down).
//   if 'rel' values are negative it moves in the oppositve direction
// returns 'from' + 'rel'
func (r *Render) move(from, rel Coord) Coord {
	r.out.CursorDown(int(rel.Y))
	r.out.CursorForward(int(rel.X))
	return from.Add(rel)
}

func (r *Render) promptHome(from Coord) {
	r.move(Coord{}, Coord{-from.X, -from.Y})
}

func (r *Render) lineWrap(cursor Column) bool {
	if runtime.GOOS != "windows" && cursor > 0 && cursor%r.termWidth == 0 {
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

func (r *Render) renderWindowTooSmall() {
	r.out.CursorGoTo(0, 0)
	r.out.EraseScreen()
	r.out.SetColor(Red, White, false)
	r.out.WriteStr("Your console window is too small...")
	return
}
