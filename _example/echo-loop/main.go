package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	prompt "github.com/c-bata/go-prompt"
)

func completer(in prompt.Document) []prompt.Choice {
	s := []prompt.Choice{
		{"areallyunnecessarilylongword", "with an unnecessary and excessively long description text being fantastically extraneous than what would fit in a 'normal'-width terminal window which will for sure cause problems that needs to be handled gracefully"},
		{"let'smakethiscomplicated", "不必要で過度に長い説明文は、「普通の」幅の端末ウィンドウに収まるものよりも幻想的に無関係であり、確実に適切に処理される必要がある問題を引き起こすでしょう。"},
		{"second", "abc"},
		{"third", "def"},
		//{"fourth", "aaa"},
		//{"five", "bbb"},
		//{"six", "ccc"},
		//{"seven", "ddd"},
		//{"eight", "eee"},
		//{"nine", "fff"},
		//{"ten", "ggg"},
		//{"eleven", "hhh"},
		//{"twelve", "iii"},
		//{"thirteen", "j"},
		//{"fourteen", "kkk"},
		//{"fifteen", "lll"},
		//{"sixteen", "mmm"},
		//{"seventeen", "nnn"},
		//{"eighteen", "ooo"},
		//{"nineteen", "ppp"},
		//{"twenty", "qqq"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func executor(text string) {
	fmt.Printf("\x1b[34myou entered:\x1b[m \x1b[1m%s\x1b[m\n", text)
}

type MLOutdentMode int

const (
	MLOutdentOnEmpty MLOutdentMode = iota
	MLOutdentOnSecondEmpty
)

type ActonEditor struct {
	ml                  bool
	indent_str          string
	outdent_after_empty MLOutdentMode
}

func (e *ActonEditor) indent(b *prompt.Buffer) {
	b.InsertText(e.indent_str, false, true)
}

func (e *ActonEditor) outdent(b *prompt.Buffer) {
	// TODO: delete from beginning of the line?
	b.DeleteBeforeCursor(prompt.Offset(len(e.indent_str)))
}

func (e *ActonEditor) on_end_line(ev *prompt.Event) {
	// decide whether to insert a new-line or end the edit.

	buf := ev.Buffer()
	doc := buf.Document()

	add_line := func(indent bool) {
		buf.NewLine(true) // copy_margin=not in_paste_mode()
		if indent {
			e.indent(buf)
		}
		ev.SetTranslatedKey(prompt.Ignore)
	}

	e.ml = doc.LineCount() > 1

	line := doc.CurrentLineBeforeCursor()

	// if we're at the end of the line with a trailing ':'
	indent_qual := len(line) >= 2 && strings.HasSuffix(line, ":") && doc.CursorAtEndOfLine()

	if indent_qual {
		add_line(true)

	} else if e.ml {
		// qualified for outdent if the current line is empty
		outdent_qual := doc.CursorAtEndOfLine() && doc.CursorOnLastLine() && doc.CursorRow() > 0 && len(strings.TrimSpace(line)) == 0

		if outdent_qual && e.outdent_after_empty == MLOutdentOnSecondEmpty {
			// outdent if the previous line is (also) empty
			prev_line := doc.Lines()[doc.CursorRow()-1]
			outdent_qual = len(strings.TrimSpace(prev_line)) == 0
		}

		if outdent_qual {
			// outdent when applicable, else end input
			if strings.HasPrefix(line, e.indent_str) {
				add_line(false)
				e.outdent(buf)
			} else {
				ev.SetEndEdit()
			}
		} else if strings.HasPrefix(line, e.indent_str) && len(strings.TrimSpace(line)) > 0 {
			// if current non-empty line is indented, we should add a line (same indentation as current)
			add_line(false)
		} else {
			ev.SetEndEdit()
		}
	}
}

/*
// if Tab is pressed when there's no text before the cursor, just insert indentation.
func (e *ActonEditor) on_tab(ev *prompt.Event) {
	buf := ev.Buffer()
	doc := buf.Document()
	if len(strings.TrimSpace(doc.CurrentLineBeforeCursor())) == 0 {
		buf.InsertText(e.indent_str, false, true)
		ev.SetTranslatedKey(prompt.Ignore)
	}
}
*/

func main() {
	e := ActonEditor{
		indent_str: "    ",
	}

	p := prompt.New(
		executor,
		completer,
		prompt.OptionCompleteAsYouType(false),
		prompt.OptionLivePrefix(func(_ *prompt.Document, row prompt.Row) (prefix string, active bool) {
			if row == 0 {
				return "main⯈", true
			} else {
				return fmt.Sprintf("…%d⯈", row), true
			}
			return
		}),
		prompt.OptionBindKey(prompt.KeyBind{prompt.Enter, e.on_end_line}),
		prompt.OptionDescriptionBGColor(prompt.NewRGB(40, 25, 50)),
		prompt.OptionDescriptionTextColor(prompt.NewRGB(120, 120, 40)),
	)

	if false {
		go func() {
			for {
				time.Sleep(3 * time.Second)
				fmt.Fprintln(os.Stderr, "async output!")
				p.OutputAsync("asynchronously text output  %v\n", time.Now().Format(time.RFC3339))
			}
		}()
	}

	os.Exit(p.Run())
}
