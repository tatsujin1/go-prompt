package prompt

import (
	"testing"
)

func TestPosixParserGetKey(t *testing.T) {
	scenarioTable := []struct {
		name     string
		input    ControlSequence
		expected KeyCode
	}{
		{
			name:     "escape",
			input:    "\x1b",
			expected: Escape,
		},
		{
			name:     "undefined",
			input:    "a",
			expected: Undefined,
		},
	}

	for _, s := range scenarioTable {
		t.Run(s.name, func(t *testing.T) {
			key := GetKey(s.input)
			if key != s.expected {
				t.Errorf("Should be %s, but got %s", key, s.expected)
			}
		})
	}
}
