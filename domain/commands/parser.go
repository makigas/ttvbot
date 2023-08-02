package commands

import "strings"

func parseCommand(cmd string) []string {
	tokens := strings.Fields(cmd)
	if strings.HasPrefix(tokens[0], "!") {
		tokens[0] = tokens[0][1:]
		return tokens
	}
	return nil
}
