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
			input:    "     hello world    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO  WORLD",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		result := cleanInput(c.input)

		lenResult := len(result)
		lenExpected := len(c.expected)
		if lenResult != lenExpected {
			t.Errorf("Expected %d words, got %d", lenExpected, lenResult)
		}

		for i := range result {
			word := result[i]
			expected := c.expected[i]

			if word != expected {
				t.Errorf("Expected %s, got %s", expected, word)
			}
		}
	}
}

func TestCommmands(t *testing.T) {
	config := createConfig()
	cmds := createCommands(config)

	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "help",
			expected: "Displays a help message",
		},
		{
			input:    "exit",
			expected: "Exit the Pokedex",
		},
	}

	for _, c := range cases {
		cmd := cmds[c.input]
		if cmd.Description != c.expected {
			t.Errorf("Expected %s, got %s", c.expected, cmd.Description)
		}
	}
}
