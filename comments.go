package main

import "strings"

// Handles both inline comments and comment-only lines
func removeHashComments(input string) string {
	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		cleaned := removeCommentFromLine(line)
		// Only keep the line if it has content after removing comments
		if strings.TrimSpace(cleaned) != "" {
			result = append(result, cleaned)
		} else if cleaned != "" {
			// Keep empty lines that weren't comment-only lines
			result = append(result, cleaned)
		}
	}

	return strings.Join(result, "\n")
}

// removeCommentFromLine removes the comment portion from a single line
func removeCommentFromLine(line string) string {
	inQuotes := false
	var quoteChar rune

	for i, char := range line {
		// Handle quote tracking to avoid removing # inside strings
		if char == '"' || char == '\'' || char == '`' {
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
			}
		}

		// If we find # and we're not inside quotes, remove from here to end
		if char == '#' && !inQuotes {
			// Return the line up to the comment, with trailing whitespace removed
			return strings.TrimRight(line[:i], " \t")
		}
	}

	// No comment found, return original line
	return line
}
