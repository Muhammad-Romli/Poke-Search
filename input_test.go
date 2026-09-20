package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "    hello Arceus     ",
			expected: []string{"hello", "arceus"},
		},
		{
			input:    "    Arceus defeat Gyaradous...    ",
			expected: []string{"arceus", "defeat", "gyaradous..."},
		},
		{
			input:    "    Groudon       VS     Kyogre    ",
			expected: []string{"groudon", "vs", "kyogre"},
		},
		{
			input:    "    ",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("This does not match: %s != %s", actual, c.expected)
			continue
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("This does not match: %s != %s", actual[i], c.expected[i])
			}
		}
	}
}
