package prompt

import (
	"strings"

	"github.com/c-bata/go-prompt/internal/debug"
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
	textMargin       = runewidth.StringWidth(textPrefix + textSuffix)
	descMargin       = runewidth.StringWidth(descPrefix + descSuffix)
	completionMargin = textMargin + descMargin
)

// Choice is printed when completing.
type Choice struct {
	Text        string
	Description string
}

// CompletionManager manages which choice is now selected.
type CompletionManager struct {
	selected  int // -1 means nothing is selected.
	choices   []Choice
	max       uint16
	completer Completer

	verticalScroll int
	wordSeparator  string
	showAtStart    bool
}

// NewCompletionManager returns initialized CompletionManager object.
func NewCompletionManager(completer Completer, max uint16) *CompletionManager {
	return &CompletionManager{
		selected:  -1,
		max:       max,
		completer: completer,

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

// Choices returns the list of choices.
func (c *CompletionManager) Choices() []Choice {
	return c.choices
}

// Reset to select nothing.
func (c *CompletionManager) Reset() {
	c.selected = -1      // nothing selected
	c.verticalScroll = 0 // scrolling at the top
	c.Update(*NewDocument())
	return
}

// Update to update the choices.
func (c *CompletionManager) Update(in Document) {
	c.choices = c.completer(in)

	for idx, choice := range c.choices {
		choice.Text = deleteBreakLineCharacters(choice.Text)
		choice.Description = deleteBreakLineCharacters(choice.Description)
		c.choices[idx] = choice
	}
	return
}

// Next to select the next choice.
func (c *CompletionManager) Next() {
	if c.verticalScroll+int(c.max)-1 == c.selected {
		c.verticalScroll++
	}
	c.selected++
	c.update()
	return
}

// Previous to select the previous choice.
func (c *CompletionManager) Previous() {
	if c.verticalScroll == c.selected && c.selected > 0 {
		c.verticalScroll--
	}
	c.selected--
	c.update()
	return
}

// Completing returns whether some suggestion is currently selected.
func (c *CompletionManager) Completing() bool {
	return c.selected != -1
}

func (c *CompletionManager) update() {
	max := int(c.max)
	if len(c.choices) < max {
		max = len(c.choices)
	}

	if c.selected >= len(c.choices) {
		c.Reset()
	} else if c.selected < -1 {
		c.selected = len(c.choices) - 1
		c.verticalScroll = len(c.choices) - max
	}
}

func deleteBreakLineCharacters(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

func formatTexts(texts []string, maxWidth int, prefix, suffix string) (formatted []string, totalWidth int) {
	wPrefix := runewidth.StringWidth(prefix)
	wSuffix := runewidth.StringWidth(suffix)
	wEllipsis := runewidth.StringWidth(ellipsis)

	min := wPrefix + wSuffix + wEllipsis
	// find longest text
	for _, text := range texts {
		w := runewidth.StringWidth(text)
		if w > totalWidth {
			totalWidth = w
		}
	}

	if totalWidth == 0 {
		return nil, 0
	}
	if min >= maxWidth {
		return nil, 0
	}

	if wPrefix+totalWidth+wSuffix > maxWidth {
		totalWidth = maxWidth - wPrefix - wSuffix
	}

	formatted = make([]string, len(texts))
	for idx, text := range texts {
		w := runewidth.StringWidth(text)
		if w > totalWidth {
			text = runewidth.Truncate(text, totalWidth, ellipsis)
			// runewidth.Truncate("您好xxx您好xxx", 11, "...") will "您好xxx..." (i.e. width 10),
			// so we need recalculate the width (and pad it at the end)
			w = runewidth.StringWidth(text)
		}
		if w <= totalWidth {
			text = runewidth.FillRight(text, totalWidth)
		}
		formatted[idx] = prefix + text + suffix
	}
	return formatted, totalWidth
}

func formatChoices(choices []Choice, maxWidth int) (formatted []Choice, totalWidth int) {
	count := len(choices)

	texts := make([]string, count)
	descs := make([]string, count)
	for idx, c := range choices {
		texts[idx] = c.Text
		descs[idx] = c.Description
	}

	var textWidth int
	var descWidth int

	texts, textWidth = formatTexts(texts, maxWidth, textPrefix, textSuffix)
	if textWidth == 0 {
		return []Choice{}, 0
	}
	descs, descWidth = formatTexts(descs, maxWidth-textWidth, descPrefix, descSuffix)

	formatted = make([]Choice, count)
	for idx := range texts {
		formatted[idx] = Choice{Text: texts[idx], Description: descs[idx]}
	}
	totalWidth = textWidth + descWidth

	return
}
