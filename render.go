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
	out                ConsoleWriter
	prefix             string
	livePrefixCallback func() (prefix string, useLivePrefix bool)
	title              string
	termHeight         int
	termWidth          int

	previousCursor      Coord
	previousLineCount   int
	previousRenderLines int

	Colors             RenderColors
	TrueColorSupported bool

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

func NewRender(prefix string, w ConsoleWriter) *Render {
	r := &Render{
		prefix: prefix,
		out:    w,
		Colors: defaultColors,

		previousLineCount:   1,

		livePrefixCallback: func() (string, bool) { return "", false },

		outputLock: &sync.Mutex{},
	}

	// https://gist.github.com/XVilka/8346728#detection
	cterm := os.Getenv("COLORTERM")
	if cterm == "truecolor" || cterm == "24bit" {
		r.TrueColorSupported = true
	}

	return r
}

func (r *Render) ValidateColor(c Color) (Color, bool) {
	if r.TrueColorSupported || !c.IsTrueColor() {
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
	r.termHeight = int(ws.Row)
	r.termWidth = int(ws.Col)
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
	//dbg("------------------------- RENDER")

	// if lines have been added to the edit, add space
	lcount := doc.LineCount()
	added := 0
	if lcount > r.previousLineCount {
		r.out.WriteRaw([]byte{'\n'})
		added = 1
		//dbg("added LF  (%d -> %d)", r.previousLineCount, doc.LineCount())
	}

	// move to beginning of the current prompt
	r.promptHome(Coord{r.previousCursor.X, r.previousCursor.Y + added})

	prefix := r.getCurrentPrefix()
	// calculate future cursor position after prefix & line is printed
	editPoint := doc.DisplayCursorCoordWithPrefix(r.termWidth, prefix)
	//dbg("editPoint @ %+v", editPoint)

	// prepare area
	h := lcount + 1 + int(compMgr.MaxVisibleChoices())
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
		editPoint = r.move(editPoint, Coord{-runewidth.StringWidth(completing_word), 0})

		// write the choice, using the configured preview style
		r.out.SetColor(r.Colors.previewChoiceText, r.Colors.previewChoiceBG, false)
		r.out.WriteStr(choice.Text)
		// move edit point to the end of the suggested word
		editPoint.X += runewidth.StringWidth(choice.Text)
		r.out.SaveCursor()

		// write the text following the cursor (using default style)
		r.out.SetColor(DefaultColor, DefaultColor, false)
		rest := buf.Document().TextAfterCursor()
		r.out.WriteStr(rest)
		// total length of line
		eol := editPoint.X + runewidth.StringWidth(rest)
		// move cursor back to the edit point
		if r.lineWrap(eol) { // output LF if necessary
			dbg("choice wrapped!\n")
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
	editPoint := doc.DisplayCursorCoordWithPrefix(r.termWidth, r.getCurrentPrefix())
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
	r.renderPrefix()
	r.out.SetColor(r.Colors.inputText, r.Colors.inputBG, false)
	r.out.WriteStr(doc.Text)
	if breakLine {
		r.out.WriteRawStr("\n")
	}
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

func (r *Render) renderCompletion(buf *Buffer, compMgr *CompletionManager) int {
	if compMgr.NumChoices() == 0 {
		return 0
	}
	editPoint := buf.Document().DisplayCursorCoordWithPrefix(r.termWidth, r.getCurrentPrefix())

	widthLimit := r.termWidth - editPoint.X - scrollbarWidth - safetyMargin

	formatted, width, withDesc := compMgr.FormatChoices(widthLimit, r.termWidth)
	width += scrollbarWidth

	windowHeight := len(formatted)
	if windowHeight > int(compMgr.MaxVisibleChoices()) {
		windowHeight = int(compMgr.MaxVisibleChoices())
	}

	cursorMoved := 0

	if r.termWidth-editPoint.X < 40 || editPoint.X+width >= r.termWidth {
		cursorMoved = -editPoint.X + 10 // say, at column 10 :)
		r.move(Coord{}, Coord{cursorMoved, 0})
		// re-format the choices, we now have more space
		widthLimit = r.termWidth - (editPoint.X - cursorMoved) - scrollbarWidth - safetyMargin
		formatted, width, withDesc = compMgr.FormatChoices(widthLimit, r.termWidth)
		width += scrollbarWidth
	}

	formatted = formatted[compMgr.verticalScroll : compMgr.verticalScroll+windowHeight]
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

	for i := 0; i < windowHeight; i++ {
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

	return windowHeight
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
	// TODO: should the prefix include style information?
	r.out.SetColor(r.Colors.prefixText, r.Colors.prefixBG, false)
	r.out.WriteRawStr(r.getCurrentPrefix())
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// move moves the cursor in the 'rel' direction (right & down).
//   if 'rel' values are negative it moves in the oppositve direction
// returns 'from' + 'rel'
func (r *Render) move(from, rel Coord) Coord {
	//dbg("move: %+v", rel)
	r.out.CursorDown(rel.Y)
	r.out.CursorForward(rel.X)
	return from.Add(rel)
}

func (r *Render) promptHome(from Coord) {
	//dbg("promptHome: %+v", from)
	r.move(Coord{}, Coord{-from.X, -from.Y})
}

// toCoord returns the relative position from the beginning of the string.
func (r *Render) toCoord(cursor int) Coord {
	col := int(r.termWidth)
	return Coord{cursor % col, cursor / col}
}

func (r *Render) lineWrap(cursor int) bool {
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
