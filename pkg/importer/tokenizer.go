package importer

import (
	"strings"
)

// Tokenize parses a raw cURL command string (potentially multiline from bash, zsh, cmd, PowerShell, Chrome, Firefox, Safari, Edge)
// into discrete argument strings, respecting quotes and line continuation characters.
func Tokenize(input string) []string {
	var tokens []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false

	runes := []rune(input)
	n := len(runes)

	for i := 0; i < n; i++ {
		ch := runes[i]

		if escaped {
			current.WriteRune(ch)
			escaped = false
			continue
		}

		if ch == '\\' && !inSingleQuote {
			if i+1 < n && (runes[i+1] == '\n' || runes[i+1] == '\r') {
				continue
			}
			if i+2 < n && runes[i+1] == '\r' && runes[i+2] == '\n' {
				i++
				continue
			}
			escaped = true
			continue
		}

		// PowerShell line continuation (` at end of line)
		if ch == '`' && !inSingleQuote && !inDoubleQuote {
			if i+1 < n && (runes[i+1] == '\n' || runes[i+1] == '\r') {
				continue
			}
		}

		// CMD line continuation (^ at end of line)
		if ch == '^' && !inSingleQuote && !inDoubleQuote {
			if i+1 < n && (runes[i+1] == '\n' || runes[i+1] == '\r') {
				continue
			}
		}

		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}

		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		if (ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r') && !inSingleQuote && !inDoubleQuote {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteRune(ch)
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
