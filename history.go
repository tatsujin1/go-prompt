package prompt

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// History stores the texts that are entered.
type History struct {
	history  []entry
	modified map[int]string
	selected int
}

type entry struct {
	time time.Time
	text string
}

// NewHistory returns new history object.
func NewHistory() *History {
	return &History{
		history: []entry{},
		//modified:      []string{""},
		modified: map[int]string{},
		selected: -1, // nothing selected
	}
}

// Save writes all history entries to 'w', returning number of entries written and/or an error.
func (h *History) Save(w io.Writer) (int, error) {
	entries := 0
	for _, entry := range h.history {
		_, err := fmt.Fprintf(w, "%d ; %s\n", entry.time.Nanosecond()/1000000, strings.Replace(entry.text, "\n", "\\n", -1))
		if err != nil {
			return entries, err
		}
		entries++
	}

	return entries, nil
}

// Load reads all history entries from 'r', returning number of entries read and/or an error.
func (h *History) Load(r io.Reader) (int, error) {
	entries := 0
	for {
		var ms int
		var text string
		_, err := fmt.Scanf("%d ; %s\n", &ms, &text)
		if err != nil {
			return entries, nil
		}
		text = strings.Replace(text, "\\n", "\n", -1)
		stamp := time.Time{}.Add(time.Duration(ms) * time.Millisecond)
		h.history = append(h.history, entry{stamp, text})

		entries++
	}
	h.ClearModified()

	return entries, nil
}

func (h *History) dump(fp *os.File) {
	for idx, item := range h.history {
		fmt.Fprintf(fp, "history[%d] \x1b[34;1m%s\x1b[m", idx, item.text)
		if mod, ok := h.modified[idx]; ok {
			fmt.Fprintf(fp, " -> \x1b[34;1m%s\x1b[m", mod)
		}
		if idx == h.selected {
			fmt.Fprintf(fp, " \x1b[1m*\x1b[m")
		}
		fmt.Fprintln(fp)
	}
}

// Add to add text in history.
func (h *History) Add(input string) {
	h.history = append(h.history, entry{time.Now(), input})
	h.ClearModified()

	h.dump(os.Stderr)
}

func (h *History) AddMany(many []string) {
	t := time.Now()
	for _, s := range many {
		h.history = append(h.history, entry{t, s})
	}
	h.ClearModified()

	h.dump(os.Stderr)
}

// ClearModified to clear the modified history entries.
func (h *History) ClearModified() {
	//fmt.Fprintln(os.Stderr, "history: clearing modified")
	h.modified = make(map[int]string, 10)
	h.selected = -1
}

// Previous saves a buffer of current line and get a buffer of previous line by up-arrow.
// The changes of line buffers are stored until a history entry is added.
func (h *History) Previous(buf *Buffer) *Buffer {
	if len(h.history) == 0 || h.selected == 0 {
		// no history, or already at the oldest entry
		h.dump(os.Stderr)
		return buf
	}

	// save the text if it was modified
	text := buf.Text()
	if h.selected == -1 || text != h.history[h.selected].text {
		h.modified[h.selected] = text
	}

	// nothing selected, select the most recent entry
	if h.selected == -1 {
		h.selected = len(h.history) - 1
	} else {
		h.selected--
	}

	h.dump(os.Stderr)

	// get text from entry
	return h.entryText()
}

// Next saves a buffer of current line and get a buffer of next line.
// The changes of line buffers are stored until a history entry is added.
func (h *History) Next(buf *Buffer) *Buffer {
	if h.selected == -1 {
		// "next" at 'nothing selected' does nothing
		h.dump(os.Stderr)
		return buf
	}

	// save the text if it was modified
	text := buf.Text()
	if h.history[h.selected].text != text {
		h.modified[h.selected] = text
	}

	if h.selected >= len(h.history)-1 {
		// already at the first entry
		h.selected = -1
	} else {
		h.selected++
	}

	h.dump(os.Stderr)

	// get text from entry
	return h.entryText()
}

func (h *History) entryText() *Buffer {
	b := NewBuffer()
	// use the modified text, if any
	if text, ok := h.modified[h.selected]; ok {
		b.InsertText(text, false, true)
	} else {
		b.InsertText(h.history[h.selected].text, false, true)
	}

	return b
}
