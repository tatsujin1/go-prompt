package prompt

import (
	"strings"

	"github.com/tatsujin/go-prompt/internal/debug"
	runewidth "github.com/mattn/go-runewidth"
)

const (
	ellipsis   = "…"
	textPrefix = " "
	textSuffix = " "
	descPrefix = " "
	descSuffix = " "
)

var (
	textMargin       = Column(runewidth.StringWidth(textPrefix + textSuffix))
	descMargin       = Column(runewidth.StringWidth(descPrefix + descSuffix))
	completionMargin = textMargin + descMargin
)

// CompletionManager manages which choice is now selected.
type CompletionManager struct {
	selected          int // -1 means nothing is selected.
	choices           []Choice
	maxVisibleChoices int
	displayMode       DisplayMode
	completer         Completer

	verticalScroll int
	wordSeparator  string
	asYouType      bool
	showAtStart    bool

	formatCache *formattedChoices
}

// Choice is printed when completing.
type Choice struct {
	Text        string
	Description string
}

type formattedChoices struct {
	maxWidth   Column
	termWidth  Column
	formatted  []Choice
	totalWidth Column
	useDesc    bool
}

type DisplayMode int

const (
	SingleColumnDescription DisplayMode = iota
	SingleColumn                        // always single column
	MultiColumn                         // multi-column if possible
)

// NewCompletionManager returns initialized CompletionManager object.
func NewCompletionManager(completer Completer, maxVisible uint16) *CompletionManager {
	return &CompletionManager{
		selected:          -1,
		maxVisibleChoices: int(maxVisible),
		displayMode:       SingleColumnDescription,
		completer:         completer,
		asYouType:         true, // by default true, to be "backwards compatible"

		verticalScroll: 0,
	}
}

// Selected returns the selected choice.
func (c *CompletionManager) Selected() (s Choice, ok bool) {
	if c.selected == -1 {
		return Choice{}, false
	} else if c.selected < -1 {
		debug.Assert(false, "must not reach here")
		c.selected = -1
		return Choice{}, false
	}
	return c.choices[c.selected], true
}

// NumChoices returns the number of choices.
func (c *CompletionManager) NumChoices() int {
	return len(c.choices)
}

// Choices returns the list of choices.
func (c *CompletionManager) Choices() []Choice {
	return c.choices
}

// Reset to select nothing.
func (c *CompletionManager) Reset() {
	c.selected = -1      // nothing selected
	c.verticalScroll = 0 // scrolling at the top
	c.choices = nil
}

// Complete to generate a list of choices.
func (c *CompletionManager) FindCompletions(in Document) {
	c.choices = c.completer(in)

	for idx, choice := range c.choices {
		choice.Text = deleteBreakLineCharacters(choice.Text)
		choice.Description = deleteBreakLineCharacters(choice.Description)
		c.choices[idx] = choice
	}
}

// Next to select the next choice.
func (c *CompletionManager) Next() {
	// TODO: this scrolling handling should not be here,
	//   but in a "completion view/renderer"-type-thingie
	if c.verticalScroll+c.maxVisibleChoices-1 == c.selected {
		c.verticalScroll++
	}
	c.selected++
	c.update()
}

// Previous to select the previous choice.
func (c *CompletionManager) Previous() {
	// TODO: this scrolling handling should not be here,
	//   but in a "completion view/renderer"-type-thingie
	if c.verticalScroll == c.selected && c.selected > 0 {
		c.verticalScroll--
	}
	c.selected--
	c.update()
}

// Completing returns whether some suggestion is currently selected.
func (c *CompletionManager) Completing() bool {
	return c.selected != -1
}

func (c *CompletionManager) MaxVisibleChoices() int {
	return c.maxVisibleChoices
}

func (c *CompletionManager) update() {
	visible := c.maxVisibleChoices
	if len(c.choices) < visible {
		visible = len(c.choices)
	}

	if c.selected >= len(c.choices) {
		c.Reset()
	} else if c.selected < -1 {
		c.selected = len(c.choices) - 1
		c.verticalScroll = len(c.choices) - visible
	}
}

func (c *CompletionManager) FormatChoices(maxWidth, termWidth Column) (formatted []Choice, totalWidth Column, useDesc bool) {
	// TODO: this formatting should not be here,
	//   but in a "completion view/renderer"-type-thingie

	if c.formatCache != nil &&
		c.formatCache.maxWidth == maxWidth &&
		c.formatCache.termWidth == termWidth {
		return c.formatCache.formatted, c.formatCache.totalWidth, c.formatCache.useDesc
	}
	c.formatCache = nil

	count := len(c.choices)

	texts := make([]string, count)
	descs := make([]string, count)
	for idx, c := range c.choices {
		texts[idx] = c.Text
		descs[idx] = c.Description
	}

	useDesc = c.displayMode == SingleColumnDescription

	var textWidth Column
	var descWidth Column

	textMaxWidth := maxWidth

	if maxWidth < 20 {
		if termWidth > 50 {
			maxWidth = termWidth / 2
			textMaxWidth = maxWidth / 2
		} else {
			// nah, this terminal is too narrow!
			return
		}
	}

	//fmt.Fprintf(os.Stderr, "text column: (maxWidth: %d)\n", textMaxWidth)
	texts, textWidth = formatTexts(texts, textMaxWidth, textPrefix, textSuffix)
	if textWidth == 0 {
		return
	}

	formatted = make([]Choice, count)
	if useDesc {
		remainWidth := maxWidth - textWidth
		//fmt.Fprintf(os.Stderr, "description column: (maxWidth: %d)\n", remainWidth)
		descs, descWidth = formatTexts(descs, remainWidth, descPrefix, descSuffix)

		useDesc = descWidth > 0
	}

	for idx := range texts {
		var desc string
		if useDesc {
			desc = descs[idx]
		}
		formatted[idx] = Choice{
			Text:        texts[idx],
			Description: desc,
		}
	}
	totalWidth = textWidth + descWidth

	c.formatCache = &formattedChoices{
		maxWidth:   maxWidth,
		termWidth:  termWidth,
		formatted:  formatted,
		totalWidth: totalWidth,
		useDesc:    useDesc,
	}
	return
}

func deleteBreakLineCharacters(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

func formatTexts(texts []string, maxWidth Column, prefix, suffix string) (formatted []string, totalWidth Column) {
	wPrefix := Column(runewidth.StringWidth(prefix))
	wSuffix := Column(runewidth.StringWidth(suffix))
	wEllipsis := Column(runewidth.StringWidth(ellipsis))

	if wPrefix+wSuffix+wEllipsis >= maxWidth {
		// we don't seem to have space for anything!?
		return nil, 0
	}

	// find widest text
	var widest Column
	for _, text := range texts {
		w := Column(runewidth.StringWidth(text))
		if w > widest {
			widest = w
		}
	}
	if widest == 0 {
		return nil, 0
	}

	if wPrefix+widest+wSuffix > maxWidth {
		// need to limit thew text idth
		widest = maxWidth - wPrefix - wSuffix
	}

	widthLimit := widest
	//fmt.Fprintf(os.Stderr, "widthLimit: %d\n", widthLimit)

	formatted = make([]string, len(texts))
	for idx, text := range texts {
		w := Column(runewidth.StringWidth(text))
		if w > widthLimit {
			text = runewidth.Truncate(text, int(widthLimit), ellipsis)
			// runewidth.Truncate("您好xxx您好xxx", 11, "...") will "您好xxx..." (i.e. width 10),
			// so we need to recalculate the width (and pad it at the end if necessary)
		}
		text = runewidth.FillRight(text, int(widthLimit))

		formatted[idx] = prefix + text + suffix
		//fmt.Fprintf(os.Stderr, "-'%s' (%d)\n", formatted[idx], runewidth.StringWidth(formatted[idx]))
	}
	return formatted, widthLimit + wPrefix + wSuffix
}
