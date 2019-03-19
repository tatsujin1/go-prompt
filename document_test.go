package prompt

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleDocument_CurrentLine() {
	d := NewDocument(
		`Hello! my name is c-bata.
This is an example of the Document component.
This component has texts displayed in terminal and cursor position.
`,
		len(`Hello! my name is c-bata.
This is an exam`),
	)
	fmt.Println("CurrentLine", d.CurrentLine())
	// Output:
	// CurrentLine This is an example of the Document component.
}

func ExampleDocument_CursorTextColumn() {
	d := NewDocument(`Hello! my name is c-bata.`, len(`Hello`))
	fmt.Println("CursorTextColumn", d.CursorTextColumn())
	// Output:
	// CursorTextColumn 5
}

func ExampleDocument_CursorTextColumn_withJapanese() {
	d := NewDocument(`こんにちは、芝田 将です。`, len([]rune("こ")))
	fmt.Println("CursorTextColumn", d.CursorTextColumn())
	// Output:
	// CursorTextColumn 2

	// (`こ` is 2 terminal columns wide)
}

func ExampleDocument_CursorRow() {
	d := NewDocument(
		`Hello! my name is c-bata.
This is an example of the Document component.
This component has texts displayed in terminal and cursor position.
`,
		len(`Hello! my name is c-bata.
This is an exam`),
	)
	fmt.Println("CursorRow", d.CursorRow())
	// Output:
	// CursorRow 1
}

func ExampleDocument_CursorColumnIndex() {
	d := NewDocument(
		`Hello! my name is c-bata.
This is an example of the Document component.
This component has texts displayed in terminal and cursor position.
`,
		len(`Hello! my name is c-bata.
This is an exam`),
	)
	fmt.Println("CursorColumnIndex", d.CursorColumnIndex())
	// Output:
	// CursorColumnIndex 15
}

func ExampleDocument_TextBeforeCursor() {
	d := NewDocument(
		`Hello! my name is c-bata.
This is an example of the Document component.
This component has texts displayed in terminal and cursor position.
`,
		len(`Hello! my name is c-bata.
This is an exam`),
	)
	fmt.Println(d.TextBeforeCursor())
	// Output:
	// Hello! my name is c-bata.
	// This is an exam
}

func ExampleDocument_TextAfterCursor() {
	d := NewDocument(
		`Hello! my name is c-bata.
This is an example of the Document component.
This component has texts displayed in terminal and cursor position.
`,
		len(`Hello! my name is c-bata.
This is an exam`),
	)
	fmt.Println(d.TextAfterCursor())
	// Output:
	// ple of the Document component.
	// This component has texts displayed in terminal and cursor position.
}

func ExampleDocument_CurrentLineBeforeCursor() {
	d := NewDocument(
		`Hello! my name is c-bata.
This is an example of the Document component.
This component has texts displayed in terminal and cursor position.
`,
		len(`Hello! my name is c-bata.
This is an exam`),
	)
	fmt.Println(d.CurrentLineBeforeCursor())
	// Output:
	// This is an exam
}

func ExampleDocument_CurrentLineAfterCursor() {
	d := NewDocument(
		`Hello! my name is c-bata.
This is an example of the Document component.
This component has texts displayed in terminal and cursor position.
`,
		len(`Hello! my name is c-bata.
This is an exam`),
	)
	fmt.Println(d.CurrentLineAfterCursor())
	// Output:
	// ple of the Document component.
}

func ExampleDocument_GetWordBeforeCursor() {
	d := NewDocument(
		`Hello! my name is c-bata.
This is an example of the Document component.
`,
		len(`Hello! my name is c-bata.
This is an exam`),
	)
	fmt.Println(d.GetWordBeforeCursor())
	// Output:
	// exam
}

func ExampleDocument_GetWordAfterCursor() {
	d := NewDocument(
		`Hello! my name is c-bata.
This is an example of the Document component.
`,
		len(`Hello! my name is c-bata.
This is an exam`),
	)
	fmt.Println(d.GetWordAfterCursor())
	// Output:
	// ple
}

func ExampleDocument_GetWordBeforeCursorWithSpace() {
	d := NewDocument(
		`Hello! my name is c-bata.
This is an example of the Document component.
`,
		len(`Hello! my name is c-bata.
This is a example `),
	)
	fmt.Println(d.GetWordBeforeCursorWithSpace())
	// Output:
	// example
}

func ExampleDocument_GetWordAfterCursorWithSpace() {
	d := NewDocument(
		`Hello! my name is c-bata.
This is an example of the Document component.
`,
		len(`Hello! my name is c-bata.
This is an`),
	)
	fmt.Println(d.GetWordAfterCursorWithSpace())
	// Output:
	//  example
}

func ExampleDocument_GetWordBeforeCursorUntilSeparator() {
	d := NewDocument(
		`hello,i am c-bata`,
		len(`hello,i am c`),
	)
	fmt.Println(d.GetWordBeforeCursorUntilSeparator(","))
	// Output:
	// i am c
}

func ExampleDocument_GetWordAfterCursorUntilSeparator() {
	d := NewDocument(
		`hello,i am c-bata,thank you for using go-prompt`,
		len(`hello,i a`),
	)
	fmt.Println(d.GetWordAfterCursorUntilSeparator(","))
	// Output:
	// m c-bata
}

func ExampleDocument_GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor() {
	d := NewDocument(
		`hello,i am c-bata,thank you for using go-prompt`,
		len(`hello,i am c-bata,`),
	)
	fmt.Println(d.GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor(","))
	// Output:
	// i am c-bata,
}

func ExampleDocument_GetWordAfterCursorUntilSeparatorIgnoreNextToCursor() {
	d := NewDocument(
		`hello,i am c-bata,thank you for using go-prompt`,
		len(`hello`),
	)
	fmt.Println(d.GetWordAfterCursorUntilSeparatorIgnoreNextToCursor(","))
	// Output:
	// ,i am c-bata
}

func TestDocument_CursorColumnIndex(t *testing.T) {
	tests := []struct {
		document *Document
		expected Index
	}{
		{
			document: NewDocument(
				"hello",
				2,
			),
			expected: 2,
		},
		{
			document: NewDocument(
				"こんにちは",
				2,
			),
			expected: 2,
		},
		{
			// If you're facing test failure on this test case and your terminal is iTerm2,
			// please check 'Profile -> Text' configuration. 'Use Unicode version 9 widths'
			// must be checked.
			// https://github.com/c-bata/go-prompt/pull/99
			document: NewDocument(
				"Добрый день",
				3,
			),
			expected: 3,
		},
	}

	for _, p := range tests {
		ac := p.document.CursorColumnIndex()
		if ac != p.expected {
			t.Errorf("Expected %#v, got %#v", p.expected, ac)
		}
	}
}

func TestDocument_CursorDisplayCoord(t *testing.T) {
	tests := []struct {
		document *Document
		expected Coord
	}{
		{
			document: NewDocument(
				"hello",
				2,
			),
			expected: Coord{2, 0},
		},
		{
			document: NewDocument(
				"こんにちは",
				2,
			),
			expected: Coord{4, 0},
		},
		{
			document: NewDocument(
				"こんに\nちは",
				5,
			),
			expected: Coord{2, 1},
		},
		{
			// If you're facing test failure on this test case and your terminal is iTerm2,
			// please check 'Profile -> Text' configuration. 'Use Unicode version 9 widths'
			// must be checked.
			// https://github.com/c-bata/go-prompt/pull/99
			document: NewDocument(
				"Добрый день",
				3,
			),
			expected: Coord{3, 0},
		},
		{
			// If you're facing test failure on this test case and your terminal is iTerm2,
			// please check 'Profile -> Text' configuration. 'Use Unicode version 9 widths'
			// must be checked.
			// https://github.com/c-bata/go-prompt/pull/99
			document: NewDocument(
				"Добр\nый день",
				7,
			),
			expected: Coord{2, 1},
		},
	}

	for idx, p := range tests {
		ac := p.document.CursorDisplayCoord(80)
		if ac != p.expected {
			t.Errorf("[%d] Expected %+v, got %+v", idx, p.expected, ac)
		}
	}
}

func TestDocument_CursorDisplayCoordWithPrefix(t *testing.T) {
	tests := []struct {
		prefix   []string
		document *Document
		expected Coord
	}{
		{
			prefix: []string{"prefix⯈"}, // 7
			document: NewDocument(
				"this does not affect the result",
				0,
			),
			expected: Coord{7 + 0, 0},
		},
		{
			prefix: []string{"prefix⯈"}, // 7
			document: NewDocument(
				"hello",
				2,
			),
			expected: Coord{7 + 2, 0},
		},
		{
			prefix: []string{"prefix⯈"}, // 7
			document: NewDocument(
				"line 1\nline 2",
				2,
			),
			expected: Coord{7 + 2, 0},
		},
		{
			prefix: []string{"prefix⯈"}, // 7
			document: NewDocument(
				"こんにちは",
				len([]rune("こ")),
			),
			expected: Coord{7 + 2, 0},
		},
		{
			prefix: []string{"this long prefix is ignored because cursor is on the next line", ""},
			document: NewDocument(
				"こんに\nちは",
				len([]rune("こんに\nち")),
			),
			expected: Coord{0 + 2, 1},
		},
		{
			prefix: []string{"prefix⯈"}, // 7
			document: NewDocument(
				"Добрый день",
				len([]rune("Добр")),
			),
			expected: Coord{7 + 4, 0},
		},
		{
			prefix: []string{"this long prefix is ignored because cursor is on the next line", ""},
			document: NewDocument(
				"Добр\nый день",
				len([]rune("Добр\nый")),
			),
			expected: Coord{0 + 2, 1},
		},
	}

	for idx, p := range tests {
		pfx_f := func(d *Document, _ Row) string {
			return p.prefix[d.CursorRow()]
		}
		ac := p.document.CursorDisplayCoordWithPrefix(80, pfx_f)
		if ac != p.expected {
			t.Errorf("[%d] Expected %+v, got %+v", idx, p.expected, ac)
		}
	}
}

func TestDocument_GetCharFromCursor(t *testing.T) {
	tests := []struct {
		document *Document
		offset   Offset
		expected rune
	}{
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("line 1\nlin"),
			),
			offset:   1,
			expected: ' ',
		},
		{
			document: NewDocument(
				"あいうえお\nかきくけこ\nさしすせそ\nたちつてと\n",
				len([]rune("あいうえお\nかき")),
			),
			offset:   -1,
			expected: 'き',
		},
		{
			document: NewDocument(
				"Добрый\nдень\nДобрый день",
				len([]rune("Добрый\nде")),
			),
			offset:   0,
			expected: 'н',
		},
	}

	for i, p := range tests {
		ac := p.document.GetCharFromCursor(p.offset)
		if ac != p.expected {
			t.Errorf("[%d] Expected %q (%d), got %q (%d)", i, string(p.expected), p.expected, string(ac), ac)
		}
	}
}

func TestDocument_TextBeforeCursor(t *testing.T) {
	tests := []struct {
		document *Document
		expected string
	}{
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("line 1\n"+"lin"),
			),
			expected: "line 1\nlin",
		},
		{
			document: NewDocument(
				"あいうえお\nかきくけこ\nさしすせそ\nたちつてと\n",
				8,
			),
			expected: "あいうえお\nかき",
		},
		{
			document: NewDocument(
				"Добрый\nдень\nДобрый день",
				9,
			),
			expected: "Добрый\nде",
		},
	}
	for i, p := range tests {
		ac := p.document.TextBeforeCursor()
		if ac != p.expected {
			t.Errorf("[%d] Expected %s, got %s", i, p.expected, ac)
		}
	}
}

func TestDocument_TextAfterCursor(t *testing.T) {
	tests := []struct {
		document *Document
		expected string
	}{
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("line 1\n"+"lin"),
			),
			expected: "e 2\nline 3\nline 4\n",
		},
		{
			document: NewDocument(
				"",
				0,
			),
			expected: "",
		},
		{
			document: NewDocument(
				"あいうえお\nかきくけこ\nさしすせそ\nたちつてと\n",
				8,
			),
			expected: "くけこ\nさしすせそ\nたちつてと\n",
		},
		{
			document: NewDocument(
				"Добрый\nдень\nДобрый день",
				9,
			),
			expected: "нь\nДобрый день",
		},
	}

	for i, p := range tests {
		ac := p.document.TextAfterCursor()
		if ac != p.expected {
			t.Errorf("[%d] Expected %#v, got %#v", i, p.expected, ac)
		}
	}
}

func TestDocument_GetWordBeforeCursor(t *testing.T) {
	tests := []struct {
		document *Document
		expected string
		sep      string
	}{
		{
			document: NewDocument(
				"apple bana",
				len("apple bana"),
			),
			expected: "bana",
		},
		{
			document: NewDocument(
				"apply -f ./file/foo.json",
				len("apply -f ./file/foo.json"),
			),
			expected: "foo.json",
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple banana orange",
				len("apple ba"),
			),
			expected: "ba",
		},
		{
			document: NewDocument(
				"apply -f ./file/foo.json",
				len("apply -f ./fi"),
			),
			expected: "fi",
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple ",
				len("apple "),
			),
			expected: "",
		},
		{
			document: NewDocument(
				"あいうえお かきくけこ さしすせそ",
				8,
			),
			expected: "かき",
		},
		{
			document: NewDocument(
				"Добрый день Добрый день",
				9,
			),
			expected: "де",
		},
	}

	for i, p := range tests {
		if p.sep == "" {
			ac := p.document.GetWordBeforeCursor()
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", i, p.expected, ac)
			}
			ac = p.document.GetWordBeforeCursorUntilSeparator("")
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", i, p.expected, ac)
			}
		} else {
			ac := p.document.GetWordBeforeCursorUntilSeparator(p.sep)
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", i, p.expected, ac)
			}
		}
	}
}

func TestDocument_GetWordBeforeCursorWithSpace(t *testing.T) {
	tests := []struct {
		document *Document
		expected string
		sep      string
	}{
		{
			document: NewDocument(
				"apple bana ",
				len("apple bana "),
			),
			expected: "bana ",
		},
		{
			document: NewDocument(
				"apply -f /path/to/file/",
				len("apply -f /path/to/file/"),
			),
			expected: "file/",
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple ",
				len("apple "),
			),
			expected: "apple ",
		},
		{
			document: NewDocument(
				"path/",
				len("path/"),
			),
			expected: "path/",
			sep:      " /",
		},
		{
			document: NewDocument(
				"あいうえお かきくけこ ",
				12,
			),
			expected: "かきくけこ ",
		},
		{
			document: NewDocument(
				"Добрый день ",
				12,
			),
			expected: "день ",
		},
	}

	for _, p := range tests {
		if p.sep == "" {
			ac := p.document.GetWordBeforeCursorWithSpace()
			if ac != p.expected {
				t.Errorf("Expected %#v, got %#v", p.expected, ac)
			}
			ac = p.document.GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor("")
			if ac != p.expected {
				t.Errorf("Expected %#v, got %#v", p.expected, ac)
			}
		} else {
			ac := p.document.GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor(p.sep)
			if ac != p.expected {
				t.Errorf("Expected %#v, got %#v", p.expected, ac)
			}
		}
	}
}

func TestDocument_FindStartOfCurrentWord(t *testing.T) {
	tests := []struct {
		document *Document
		expected Index
		sep      string
	}{
		{
			document: NewDocument(
				"apple bana",
				len("apple bana"),
			),
			expected: Index(len("apple ")),
		},
		{
			document: NewDocument(
				"apply -f ./file/foo.json",
				len("apply -f ./file/foo.json"),
			),
			expected: Index(len("apply -f ./file/")),
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple ",
				len("apple "),
			),
			expected: Index(len("apple ")),
		},
		{
			document: NewDocument(
				"apply -f ./file/foo.json",
				len("apply -f ./"),
			),
			expected: Index(len("apply -f ./")),
			sep:      " /",
		},
		{
			document: NewDocument(
				"あいうえお かきくけこ さしすせそ",
				len([]rune("あいうえお かき")),
			),
			expected: Index(len([]rune("あいうえお "))),
		},
		{
			document: NewDocument(
				"Добрый день Добрый день",
				len([]rune("Добрый д")),
			),
			expected: Index(len([]rune("Добрый "))),
		},
	}

	for idx, p := range tests {
		if p.sep == "" {
			ac := p.document.FindStartOfCurrentWord()
			if ac != p.expected {
				t.Errorf("[%d/1] Expected %#v, got %#v", idx, p.expected, ac)
			}
			ac = p.document.FindStartOfCurrentWordUntilSeparator("")
			if ac != p.expected {
				t.Errorf("[%d/2] Expected %#v, got %#v", idx, p.expected, ac)
			}
		} else {
			ac := p.document.FindStartOfCurrentWordUntilSeparator(p.sep)
			if ac != p.expected {
				t.Errorf("[%d/s] Expected %#v, got %#v", idx, p.expected, ac)
			}
		}
	}
}

func TestDocument_FindStartOfCurrentWordWithSpace(t *testing.T) {
	tests := []struct {
		document *Document
		expected Index
		sep      string
	}{
		{
			document: NewDocument(
				"apple bana ",
				len("apple bana "),
			),
			expected: Index(len("apple ")),
		},
		{
			document: NewDocument(
				"apply -f /file/foo/",
				len("apply -f /file/foo/"),
			),
			expected: Index(len("apply -f /file/")),
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple ",
				len("apple "),
			),
			expected: 0,
		},
		{
			document: NewDocument(
				"file/",
				len("file/"),
			),
			expected: 0,
			sep:      " /",
		},
		{
			document: NewDocument(
				"あいうえお かきくけこ ",
				len([]rune("あいうえお かきくけこ ")), // end of string
			),
			expected: Index(len([]rune("あいうえお "))),
		},
		{
			document: NewDocument(
				"Добрый день ",
				len([]rune("Добрый день ")), // end of string
			),
			expected: Index(len([]rune("Добрый "))),
		},
	}

	for idx, p := range tests {
		if p.sep == "" {
			ac := p.document.FindStartOfCurrentWordWithSpace()
			if ac != p.expected {
				t.Errorf("[%d/1] Expected %#v, got %#v", idx, p.expected, ac)
			}
			ac = p.document.FindStartOfCurrentWordUntilSeparatorIgnoreNextToCursor("")
			if ac != p.expected {
				t.Errorf("[%d/2] Expected %#v, got %#v", idx, p.expected, ac)
			}
		} else {
			ac := p.document.FindStartOfCurrentWordUntilSeparatorIgnoreNextToCursor(p.sep)
			if ac != p.expected {
				t.Errorf("[%d/s] Expected %#v, got %#v", idx, p.expected, ac)
			}
		}
	}
}

func TestDocument_GetWordAfterCursor(t *testing.T) {
	tests := []struct {
		document *Document
		expected string
		sep      string
	}{
		{
			document: NewDocument(
				"apple bana",
				len("apple bana"),
			),
			expected: "",
		},
		{
			document: NewDocument(
				"apply -f ./file/foo.json",
				len("apply -f ./fi"),
			),
			expected: "le",
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple bana",
				len("apple "),
			),
			expected: "bana",
		},
		{
			document: NewDocument(
				"apple bana",
				len("apple"),
			),
			expected: "",
		},
		{
			document: NewDocument(
				"apply -f ./file/foo.json",
				len("apply -f ."),
			),
			expected: "",
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple bana",
				len("ap"),
			),
			expected: "ple",
		},
		{
			document: NewDocument(
				"あいうえお かきくけこ さしすせそ",
				8,
			),
			expected: "くけこ",
		},
		{
			document: NewDocument(
				"Добрый день Добрый день",
				9,
			),
			expected: "нь",
		},
	}

	for idx, p := range tests {
		if p.sep == "" {
			ac := p.document.GetWordAfterCursor()
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
			}
			ac = p.document.GetWordAfterCursorUntilSeparator("")
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
			}
		} else {
			ac := p.document.GetWordAfterCursorUntilSeparator(p.sep)
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
			}
		}
	}
}

func TestDocument_GetWordAfterCursorWithSpace(t *testing.T) {
	tests := []struct {
		document *Document
		expected string
		sep      string
	}{
		{
			document: NewDocument(
				"apple bana",
				len("apple bana"),
			),
			expected: "",
		},
		{
			document: NewDocument(
				"apple bana",
				len("apple "),
			),
			expected: "bana",
		},
		{
			document: NewDocument(
				"/path/to",
				len("/path/"),
			),
			expected: "to",
			sep:      " /",
		},
		{
			document: NewDocument(
				"/path/to/file",
				len("/path/"),
			),
			expected: "to",
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple bana",
				len("apple"),
			),
			expected: " bana",
		},
		{
			document: NewDocument(
				"path/to",
				len("path"),
			),
			expected: "/to",
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple bana",
				len("ap"),
			),
			expected: "ple",
		},
		{
			document: NewDocument(
				"あいうえお かきくけこ さしすせそ",
				5,
			),
			expected: " かきくけこ",
		},
		{
			document: NewDocument(
				"Добрый день Добрый день",
				6,
			),
			expected: " день",
		},
	}

	for idx, p := range tests {
		if p.sep == "" {
			ac := p.document.GetWordAfterCursorWithSpace()
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
			}
			ac = p.document.GetWordAfterCursorUntilSeparatorIgnoreNextToCursor("")
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
			}
		} else {
			ac := p.document.GetWordAfterCursorUntilSeparatorIgnoreNextToCursor(p.sep)
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
			}
		}
	}
}

func TestDocument_FindEndOfCurrentWord(t *testing.T) {
	tests := []struct {
		document *Document
		expected Offset
		sep      string
	}{
		{
			document: NewDocument(
				"apple bana",
				len("apple bana"),
			),
			expected: 0,
		},
		{
			document: NewDocument(
				"apple bana",
				len("apple "),
			),
			expected: Offset(len("bana")),
		},
		{
			document: NewDocument(
				"apply -f ./file/foo.json",
				len("apply -f ./"),
			),
			expected: Offset(len("file")),
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple bana",
				len("apple"),
			),
			expected: 0,
		},
		{
			document: NewDocument(
				"apply -f ./file/foo.json",
				len("apply -f ."),
			),
			expected: 0,
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple bana",
				len("ap"),
			),
			expected: Offset(len("ple")),
		},
		{
			document: NewDocument(
				"りんご ばなな",
				len([]rune("りん")),
			),
			expected: Offset(len([]rune("ご"))),
		},
		{
			document: NewDocument(
				"りんご ばなな",
				len([]rune("りんご ばなな")),
			),
			expected: 0,
		},
		{
			document: NewDocument(
				"りんご ばなな",
				len([]rune("りんご")),
			),
			expected: 0,
		},
		{
			// Доб(cursor)рый день
			document: NewDocument(
				"Добрый день",
				len([]rune("Доб")),
			),
			expected: Offset(len([]rune("рый"))),
		},
	}

	for idx, p := range tests {
		if p.sep == "" {
			ac := p.document.FindEndOfCurrentWord()
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
			}
			ac = p.document.FindEndOfCurrentWordUntilSeparator("")
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
			}
		} else {
			ac := p.document.FindEndOfCurrentWordUntilSeparator(p.sep)
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
			}
		}
	}
}

func TestDocument_FindEndOfCurrentWordWithSpace(t *testing.T) {
	tests := []struct {
		document *Document
		expected Offset
		sep      string
	}{
		{
			document: NewDocument(
				"apple bana",
				len("apple bana"),
			),
			expected: 0,
		},
		{
			document: NewDocument(
				"apple bana",
				len("apple "),
			),
			expected: Offset(len("bana")),
		},
		{
			document: NewDocument(
				"apply -f /file/foo.json",
				len("apply -f /"),
			),
			expected: Offset(len("file")),
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple bana",
				len("apple"),
			),
			expected: Offset(len(" bana")),
		},
		{
			document: NewDocument(
				"apply -f /path/to",
				len("apply -f /path"),
			),
			expected: Offset(len("/to")),
			sep:      " /",
		},
		{
			document: NewDocument(
				"apple bana",
				len("ap"),
			),
			expected: Offset(len("ple")),
		},
		{
			document: NewDocument(
				"あいうえお かきくけこ",
				6,
			),
			expected: Offset(len([]rune("かきくけこ"))),
		},
		{
			document: NewDocument(
				"あいうえお かきくけこ",
				5,
			),
			expected: Offset(len([]rune(" かきくけこ"))),
		},
		{
			document: NewDocument(
				"Добрый день",
				6,
			),
			expected: Offset(len([]rune(" день"))),
		},
	}

	for idx, p := range tests {
		if p.sep == "" {
			ac := p.document.FindEndOfCurrentWordWithSpace()
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
			}
			ac = p.document.FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor("")
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
			}
		} else {
			ac := p.document.FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor(p.sep)
			if ac != p.expected {
				t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
			}
		}
	}
}

func TestDocument_CurrentLineBeforeCursor(t *testing.T) {
	d := NewDocument(
		"line 1\nline 2\nline 3\nline 4\n",
		len("line 1\n"+"lin"),
	)
	ac := d.CurrentLineBeforeCursor()
	ex := "lin"
	if ac != ex {
		t.Errorf("Expected %#v, got %#v", ex, ac)
	}
}

func TestDocument_CurrentLineAfterCursor(t *testing.T) {
	d := NewDocument(
		"line 1\nline 2\nline 3\nline 4\n",
		len("line 1\n"+"lin"),
	)
	ac := d.CurrentLineAfterCursor()
	ex := "e 2"
	if ac != ex {
		t.Errorf("Expected %#v, got %#v", ex, ac)
	}
}

func TestDocument_CursorOnLastLine(t *testing.T) {
	tests := []struct {
		document *Document
		expected bool
	}{
		{
			document: NewDocument(
				"",
				0,
			),
			expected: true,
		},
		{
			document: NewDocument(
				"line 1\nline 2",
				len("line 1\nli"),
			),
			expected: true,
		},
		{
			document: NewDocument(
				"line 1\nline 2\n",
				len("line 1\nli"),
			),
			expected: false,
		},
	}

	for idx, p := range tests {
		actual := p.document.CursorOnLastLine()
		if actual != p.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, actual)
		}
	}
}

func TestDocument_CursorAtEndOfLine(t *testing.T) {
	tests := []struct {
		document *Document
		expected bool
	}{
		{
			document: NewDocument(
				"",
				0,
			),
			expected: true,
		},
		{
			document: NewDocument(
				"line 1\nline 2",
				len("line 1\nli"),
			),
			expected: false,
		},
		{
			document: NewDocument(
				"line 1\nline 2\n",
				len("line 1\nli"),
			),
			expected: false,
		},
	}

	for idx, p := range tests {
		actual := p.document.CursorAtEndOfLine()
		if actual != p.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, actual)
		}
	}
}

func TestDocument_CurrentLine(t *testing.T) {
	var tests = []struct {
		document *Document
		expected string
	}{
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("line 1\nlin"),
			),
			expected: "line 2",
		},
	}
	for idx, p := range tests {
		ac := p.document.CurrentLine()
		if ac != p.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, p.expected, ac)
		}
	}
}

func TestDocument_CursorRowAndCol(t *testing.T) {
	var cursorTests = []struct {
		document    *Document
		expectedRow Row
		expectedCol Index
	}{
		{
			document: NewDocument(
				"single line",
				len("single "),
			),
			expectedRow: 0,
			expectedCol: 7,
		},
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\n",
				len("line 1\nlin"),
			),
			expectedRow: 1,
			expectedCol: 3,
		},
		{
			document: NewDocument(
				"",
				0,
			),
			expectedRow: 0,
			expectedCol: 0,
		},
		{
			document: NewDocument(
				"あいうえお かき\nくけこ",
				len([]rune("あいうえお かき")),
			),
			expectedRow: 0,
			expectedCol: 8,
		},
		{
			document: NewDocument(
				"あいうえお かき\nくけこ",
				len([]rune("あいうえお かき\nくけ")),
			),
			expectedRow: 1,
			expectedCol: 2,
		},
	}
	for idx, test := range cursorTests {
		r := test.document.CursorRow()
		c := test.document.CursorColumnIndex()
		if r != test.expectedRow || c != test.expectedCol {
			t.Errorf("[%d] Expected %d:%d, got %d:%d", idx, test.expectedRow, test.expectedCol, r, c)
		}
	}
}

func TestDocument_GetCursorLeftOffset(t *testing.T) {
	var cursorTests = []struct {
		document *Document
		offset   Offset
		expected Offset
	}{
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("line 1\nline 2\nlin"),
			),
			offset:   2,
			expected: -2,
		},
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("line 1\nline 2\nlin"),
			),
			offset:   10,
			expected: -3,
		},
	}
	for idx, test := range cursorTests {
		ac := test.document.GetCursorLeftOffset(test.offset)
		if ac != test.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, test.expected, ac)
		}
	}
}

func TestDocument_GetCursorRightOffset(t *testing.T) {
	var cursorTests = []struct {
		document *Document
		offset   Offset
		expected Offset
	}{
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("line 1\nline 2\nlin"),
			),
			offset:   2,
			expected: 2,
		},
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("line 1\nline 2\nlin"),
			),
			offset:   10,
			expected: 3,
		},
	}
	for idx, test := range cursorTests {
		ac := test.document.GetCursorRightOffset(test.offset)
		if ac != test.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, test.expected, ac)
		}
	}
}

func TestDocument_GetCursorUpOffset(t *testing.T) {
	var cursorTests = []struct {
		document  *Document
		offset    Row
		preferred Index
		expected  Offset
	}{
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("line 1\nline 2\nlin"),
			),
			offset:    2,
			preferred: -1,
			expected:  Offset(len("lin") - len("line 1\nline 2\nlin")),
		},
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("line 1\nline 2\nlin"),
			),
			offset:    100,
			preferred: -1,
			expected:  Offset(len("lin") - len("line 1\nline 2\nlin")),
		},
	}
	for idx, test := range cursorTests {
		ac := test.document.GetCursorUpOffset(test.offset, test.preferred)
		if ac != test.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, test.expected, ac)
		}
	}
}

func TestDocument_GetCursorDownOffset(t *testing.T) {
	var cursorTests = []struct {
		document  *Document
		offset    Row
		preferred Index
		expected  Offset
	}{
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("lin"),
			),
			offset:    2,
			preferred: -1,
			expected:  Offset(len("e 1\nline 2\nlin")),
		},
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("lin"),
			),
			offset:    100,
			preferred: -1,
			expected:  Offset(len("e 1\nline 2\nline 3\nline 4\n")),
		},
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				len("lin"),
			),
			offset:    1,
			preferred: len("lin"),
			expected:  Offset(len("e 1\nlin")),
		},
		{
			document: NewDocument(
				"line 1\nli",
				len("line "),
			),
			offset:    1,
			preferred: len("line "),
			expected:  Offset(len("1\nli")),
		},
	}
	for idx, test := range cursorTests {
		ac := test.document.GetCursorDownOffset(test.offset, test.preferred)
		if ac != test.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, test.expected, ac)
		}
	}
}

func TestDocument_Lines(t *testing.T) {
	d := NewDocument(
		"line 1\nline 2\nline 3\nline 4\n",
		len("line 1\nlin"),
	)
	ac := d.Lines()
	ex := []string{"line 1", "line 2", "line 3", "line 4", ""}
	if !reflect.DeepEqual(ac, ex) {
		t.Errorf("Expected %#v, got %#v", ex, ac)
	}
}

func TestDocument_LineCount(t *testing.T) {
	d := NewDocument(
		"line 1\nline 2\nline 3\nline 4\n",
		len("line 1\n"+"lin"),
	)
	ac := d.LineCount()
	ex := 5
	if ac != ex {
		t.Errorf("Expected %#v, got %#v", ex, ac)
	}
}

func TestDocument_TranslateIndexToRowCol(t *testing.T) {
	var tests = []struct {
		document    *Document
		index       int
		expectedRow Row
		expectedCol Index
	}{
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				0,
			),
			index:       len("line 1\nline 2\nlin"),
			expectedRow: 2,
			expectedCol: 3,
		},
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				0,
			),
			index:       0,
			expectedRow: 0,
			expectedCol: 0,
		},
		{
			document: NewDocument(
				"こんにちは",
				0,
			),
			index:       4,
			expectedRow: 0,
			expectedCol: 4,
		},
	}

	for idx, test := range tests {
		row, col := test.document.TranslateIndexToRowCol(test.index)
		if row != test.expectedRow {
			t.Errorf("[%d] Expected row %#v, got %#v", idx, test.expectedRow, row)
		}
		if col != test.expectedCol {
			t.Errorf("[%d] Expected col %#v, got %#v", idx, test.expectedCol, col)
		}
	}
}

func TestDocument_TranslateRowColToIndex(t *testing.T) {
	var tests = []struct {
		document *Document
		row      Row
		col      Index
		expected Index
	}{
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				0,
			),
			row:      2,
			col:      3,
			expected: len("line 1\nline 2\nlin"),
		},
		{
			document: NewDocument(
				"line 1\nline 2\nline 3\nline 4\n",
				0,
			),
			row:      0,
			col:      0,
			expected: 0,
		},
		{
			document: NewDocument(
				"あいうえお かき\nくけこ",
				0,
			),
			row:      1,
			col:      2,
			expected: 11,
		},
	}
	for idx, test := range tests {
		ac := test.document.TranslateRowColToIndex(test.row, test.col)
		if ac != test.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, test.expected, ac)
		}
	}
}

func TestDocument_GetEndOfLineOffset(t *testing.T) {
	var tests = []struct {
		document *Document
		expected Offset
	}{
		{
			document: NewDocument(
				"line 1\nline 2\nline 3",
				len("line 1\nli"),
			),
			expected: Offset(len("ne 2")),
		},
		{
			document: NewDocument(
				"あいうえお かき\nくけこ",
				len([]rune("あいうえお")),
			),
			expected: Offset(len([]rune(" かき"))),
		},
	}
	for idx, test := range tests {
		ac := test.document.GetEndOfLineOffset()
		if ac != test.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, test.expected, ac)
		}
	}
}
