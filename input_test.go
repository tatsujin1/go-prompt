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
			key := FindKey(s.input)
			if key != s.expected {
				t.Errorf("Expected %d, but got %d", key, s.expected)
			}
		})
	}
}
