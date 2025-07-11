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
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  Go is awesome  ",
			expected: []string{"go", "is", "awesome"},
		},
		{
			input:    "  Leading and trailing spaces   ",
			expected: []string{"leading", "and", "trailing", "spaces"},
		},
		{
			input:    "  Mixed CASES  ",
			expected: []string{"mixed", "cases"},
		},
		{
			input:    "  Multiple   spaces  ",
			expected: []string{"multiple", "spaces"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range c.expected {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("cleanInput(%q) = %q, expected %q", c.input, actual, c.expected)
			}
		}
	}

}
