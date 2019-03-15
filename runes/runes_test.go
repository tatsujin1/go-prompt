package runes

import (
	"testing"
)

func Test_HasPrefix(t *testing.T) {

	tests := []struct {
		s        []rune
		prefix   []rune
		expected bool
	}{
		{
			s:        []rune("something nothing"),
			prefix:   []rune("some"),
			expected: true,
		},
		{
			s:        []rune("something nothing"),
			prefix:   []rune("s"),
			expected: true,
		},
		{
			s:        []rune("something nothing"),
			prefix:   []rune("thing"),
			expected: false,
		},
		{
			s:        []rune("something nothing"),
			prefix:   []rune(""),
			expected: true,
		},
		{
			s:        []rune("こんにちは"),
			prefix:   []rune("こん"),
			expected: true,
		},
		{
			s:        []rune("こんにちは"),
			prefix:   []rune("こ"),
			expected: true,
		},
		{
			s:        []rune("こんにちは"),
			prefix:   []rune("んに"),
			expected: false,
		},
	}
	for idx, tt := range tests {
		ac := HasPrefix(tt.s, tt.prefix)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}

func Test_HasSuffix(t *testing.T) {
	tests := []struct {
		s        []rune
		suffix   []rune
		expected bool
	}{
		{
			s:        []rune("something nothing"),
			suffix:   []rune("thing"),
			expected: true,
		},
		{
			s:        []rune("something nothing"),
			suffix:   []rune("not"),
			expected: false,
		},
		{
			s:        []rune("something nothing"),
			suffix:   []rune("g"),
			expected: true,
		},
		{
			s:        []rune("something nothing"),
			suffix:   []rune("some"),
			expected: false,
		},
		{
			s:        []rune("something nothing"),
			suffix:   []rune(""),
			expected: true,
		},
		{
			s:        []rune("こんにちは"),
			suffix:   []rune("ちは"),
			expected: true,
		},
		{
			s:        []rune("こんにちは"),
			suffix:   []rune("は"),
			expected: true,
		},
		{
			s:        []rune("こんにちは"),
			suffix:   []rune("こん"),
			expected: false,
		},
	}
	for idx, tt := range tests {
		ac := HasSuffix(tt.s, tt.suffix)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}

func Test_SplitRune(t *testing.T) {
	tests := []struct {
		s        []rune
		sep      rune
		expected [][]rune
	}{
		{
			s:   []rune("something nothing"),
			sep: ' ',
			expected: [][]rune{
				[]rune("something"),
				[]rune("nothing"),
			},
		},
		{
			s:        []rune(""),
			sep:      ' ',
			expected: [][]rune{[]rune{}},
		},
		{
			s:   []rune("something nothing"),
			sep: 't',
			expected: [][]rune{
				[]rune("some"),
				[]rune("hing no"),
				[]rune("hing"),
			},
		},
		{
			s:   []rune("something nothing"),
			sep: 'g',
			expected: [][]rune{
				[]rune("somethin"),
				[]rune(" nothin"),
				[]rune(""),
			},
		},
		{
			s:   []rune("こんにちは"),
			sep: 'に',
			expected: [][]rune{
				[]rune("こん"),
				[]rune("ちは"),
			},
		},
	}

	arr_eq := func(a, b [][]rune) bool {
		if len(a) != len(b) {
			return false
		}
		for idx := 0; idx < len(a); idx++ {
			if Compare(a[idx], b[idx]) != 0 {
				return false
			}
		}
		return true
	}

	for idx, tt := range tests {
		ac := SplitRune(tt.s, tt.sep)
		if !arr_eq(ac, tt.expected) {
			t.Errorf("[%d] Expected %+q, got %+q", idx, tt.expected, ac)
		}
	}
}

func Test_Index(t *testing.T) {
	tests := []struct {
		s        []rune
		needle   []rune
		expected int
	}{
		{
			s:        []rune("something nothing"),
			needle:   []rune("thing"),
			expected: len([]rune("some")),
		},
		{
			s:        []rune("something nothing"),
			needle:   []rune("e"),
			expected: len([]rune("som")),
		},
		{
			s:        []rune("こんにちは"),
			needle:   []rune("にち"),
			expected: len([]rune("こん")),
		},
		{
			s:        []rune("こんにちは"),
			needle:   []rune("asdf"),
			expected: -1,
		},
		{
			s:        []rune("こんにちは"),
			needle:   []rune(""),
			expected: 0,
		},
	}
	for idx, tt := range tests {
		ac := Index(tt.s, tt.needle)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}

func Test_IndexAny(t *testing.T) {
	tests := []struct {
		s        []rune
		any      []rune
		expected int
	}{
		{
			s:        []rune("something nothing"),
			any:      []rune("pqzxi"),
			expected: len([]rune("someth")),
		},
		{
			s:        []rune("something nothing"),
			any:      []rune("x"),
			expected: -1,
		},
		{
			s:        []rune("こんにちは"),
			any:      []rune("칹に"),
			expected: len([]rune("こん")),
		},
		{
			s:        []rune("こんにちは"),
			any:      []rune("asdf"),
			expected: -1,
		},
	}
	for idx, tt := range tests {
		ac := IndexAny(tt.s, tt.any)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}

func Test_IndexNotAny(t *testing.T) {
	tests := []struct {
		s        []rune
		any      []rune
		expected int
	}{
		{
			s:        []rune("something nothing"),
			any:      []rune("gnihto s"),
			expected: len([]rune("so")),
		},
		{
			s:        []rune("something nothing"),
			any:      []rune(""),
			expected: -1,
		},
		{
			s:        []rune("こんにちは"),
			any:      []rune("칹にこん"),
			expected: len([]rune("こんに")),
		},
		{
			s:        []rune("こんにちは"),
			any:      []rune("asdf"),
			expected: 0,
		},
	}
	for idx, tt := range tests {
		ac := IndexNotAny(tt.s, tt.any)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}

func Test_IndexRune(t *testing.T) {
	tests := []struct {
		s        []rune
		r        rune
		expected int
	}{
		{
			s:        []rune("something nothing"),
			r:        'i',
			expected: len([]rune("someth")),
		},
		{
			s:        []rune("something nothing"),
			r:        'n',
			expected: len([]rune("somethi")),
		},
		{
			s:        []rune("こんにちは"),
			r:        'ち',
			expected: len([]rune("こんに")),
		},
		{
			s:        []rune("こんにちは"),
			r:        '칹',
			expected: -1,
		},
	}
	for idx, tt := range tests {
		ac := IndexRune(tt.s, tt.r)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}

func Test_IndexNotRune(t *testing.T) {
	tests := []struct {
		s        []rune
		r        rune
		expected int
	}{
		{
			s:        []rune("ssssssssXsssXss"),
			r:        's',
			expected: len([]rune("ssssssss")),
		},
		{
			s:        []rune("ssssssssss"),
			r:        's',
			expected: -1,
		},
		{
			s:        []rune(""),
			r:        'X',
			expected: -1,
		},
		{
			s:        []rune("こここん"),
			r:        'こ',
			expected: len([]rune("こここ")),
		},
	}
	for idx, tt := range tests {
		ac := IndexNotRune(tt.s, tt.r)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}

func Test_LastIndex(t *testing.T) {
	tests := []struct {
		s        []rune
		needle   []rune
		expected int
	}{
		{
			s:        []rune("something nothing"),
			needle:   []rune("thing"),
			expected: len([]rune("something no")),
		},
		{
			s:        []rune("something nothing"),
			needle:   []rune("e"),
			expected: len([]rune("som")),
		},
		{
			s:        []rune("こんにちは"),
			needle:   []rune("にち"),
			expected: len([]rune("こん")),
		},
		{
			s:        []rune("こんにちは"),
			needle:   []rune("asdf"),
			expected: -1,
		},
		{
			s:        []rune("こんにちは"),
			needle:   []rune(""),
			expected: len([]rune("こんにちは")),
		},
	}
	for idx, tt := range tests {
		ac := LastIndex(tt.s, tt.needle)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}

func Test_Compare(t *testing.T) {
	tests := []struct {
		s1       []rune
		s2       []rune
		expected int
	}{
		{
			s1:       []rune("something nothing"),
			s2:       []rune("something nothing"),
			expected: 0,
		},
		{
			s1:       []rune("something nothingX"),
			s2:       []rune("something nothing"),
			expected: 1,
		},
		{
			s1:       []rune("something nothing"),
			s2:       []rune("something nothingX"),
			expected: -1,
		},
		{
			s1:       []rune("somethinG nothing"),
			s2:       []rune("something nothing"),
			expected: -1,
		},
		{
			s1:       []rune("something nothing"),
			s2:       []rune("somethinG nothing"),
			expected: 1,
		},
	}
	for idx, tt := range tests {
		ac := Compare(tt.s1, tt.s2)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}

func Test_LastIndexRune(t *testing.T) {
	tests := []struct {
		s        []rune
		r        rune
		expected int
	}{
		{
			s:        []rune("something nothing"),
			r:        'i',
			expected: len([]rune("something noth")),
		},
		{
			s:        []rune("something nothing"),
			r:        'n',
			expected: len([]rune("something nothi")),
		},
		{
			s:        []rune("こんにちは"),
			r:        'ち',
			expected: len([]rune("こんに")),
		},
		{
			s:        []rune("こんにちは"),
			r:        '칹',
			expected: -1,
		},
	}
	for idx, tt := range tests {
		ac := LastIndexRune(tt.s, tt.r)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}

func Test_LastIndexNotRune(t *testing.T) {
	tests := []struct {
		s        []rune
		r        rune
		expected int
	}{
		{
			s:        []rune("ssssssssXssssXss"),
			r:        's',
			expected: len([]rune("ssssssssXssss")),
		},
		{
			s:        []rune("ssssssssss"),
			r:        's',
			expected: -1,
		},
		{
			s:        []rune(""),
			r:        'X',
			expected: -1,
		},
		{
			s:        []rune("こここんこんこ"),
			r:        'こ',
			expected: len([]rune("こここんこ")),
		},
	}
	for idx, tt := range tests {
		ac := LastIndexNotRune(tt.s, tt.r)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}

func Test_LastIndexAny(t *testing.T) {
	tests := []struct {
		s        []rune
		any      []rune
		expected int
	}{
		{
			s:        []rune("something nothing"),
			any:      []rune("pqzxi"),
			expected: len([]rune("something noth")),
		},
		{
			s:        []rune("something nothing"),
			any:      []rune("x"),
			expected: -1,
		},
		{
			s:        []rune("こんにちはに"),
			any:      []rune("칹に"),
			expected: len([]rune("こんにちは")),
		},
		{
			s:        []rune("こんにちは"),
			any:      []rune("asdf"),
			expected: -1,
		},
	}
	for idx, tt := range tests {
		ac := LastIndexAny(tt.s, tt.any)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}

func Test_LastIndexNotAny(t *testing.T) {
	tests := []struct {
		s        []rune
		any      []rune
		expected int
	}{
		{
			s:        []rune("something nothing"),
			any:      []rune("gniht s"),
			expected: len([]rune("something n")),
		},
		{
			s:        []rune("something nothing"),
			any:      []rune(""),
			expected: -1,
		},
		{
			s:        []rune("こんにちは"),
			any:      []rune("칹にこん"),
			expected: len([]rune("こんにち")),
		},
		{
			s:        []rune("こんにちは"),
			any:      []rune("こんにちは"),
			expected: -1,
		},
	}
	for idx, tt := range tests {
		ac := LastIndexNotAny(tt.s, tt.any)
		if ac != tt.expected {
			t.Errorf("[%d] Expected %#v, got %#v", idx, tt.expected, ac)
		}
	}
}
