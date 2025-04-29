package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name		string
		input		string
		expected	[]string
	} {
		{
			name: "basic: lowerCase",
			input: "hello world",
			expected: []string{"hello", "world"},
		},{
			name: "basic: randomCase",
			input: "heLlo woRld",
			expected: []string{"hello", "world"},
		},{
			name: "inter: randomCase addedWhiteSpace",
			input: "\theLlo\nwoRld\n",
			expected: []string{"hello", "world"},
		},{
			name: "edge: noInput",
			input: "",
			expected: []string{},
		},{
			name: "edge: whiteSpaceOnly",
			input: " \n\t",
			expected: []string{},
		},
	}
	for n, c := range cases {
		actual := cleanInput(c.input)
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("%v (%v) -- '%v' != '%v'", n, c.name, word, expectedWord)
			}

		}
	}
}
