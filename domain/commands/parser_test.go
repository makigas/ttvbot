package commands

import (
	"fmt"
	"testing"
)

func TestCommandParse(t *testing.T) {
	testcases := []struct {
		given    string
		expected []string
	}{
		// No command.
		{given: "POG", expected: nil},

		// Simple commands and well-formed commands.
		{given: "!hug", expected: []string{"hug"}},
		{given: "!so @ibai", expected: []string{"so", "@ibai"}},
		{given: "!sr la cucaracha remix", expected: []string{"sr", "la", "cucaracha", "remix"}},

		// Let's fuzz this.
		{given: "!!agua", expected: []string{"!agua"}},
		{given: "   !hug", expected: []string{"hug"}},
		{given: "!so       @ibai", expected: []string{"so", "@ibai"}},
		{given: "⠀!agua ", expected: nil},
		{given: "  !sr    la    cucaracha     remix  ", expected: []string{"sr", "la", "cucaracha", "remix"}},
	}

	for _, tt := range testcases {
		actual := parseCommand(tt.given)
		if actual == nil && tt.expected != nil {
			t.Errorf("expected %s to yield %v, gave nil", tt.given, tt.expected)
		} else if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", tt.expected) {
			t.Errorf("expected %s to yield %v, gave %v", tt.given, tt.expected, actual)
		}
	}
}
