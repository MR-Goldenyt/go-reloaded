package transform

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func HexToDec(input string) string {
	dec, err := strconv.ParseInt(input, 16, 64)
	if err != nil {
		fmt.Printf("Error: Invalid hexadecimal number: %v\n\n", input)
		return ""
	}
	return fmt.Sprintf("%d", dec)
}

func BinToDec(input string) string {
	dec, err := strconv.ParseInt(input, 2, 64)
	if err != nil {
		fmt.Printf("Error: Invalid binary number: %v\n\n", input)
		return ""
	}
	return fmt.Sprintf("%d", dec)
}

func Capitalize(s string) string {
	runes := []rune(s)
	inWord := false

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			if !inWord {
				if r >= 'a' && r <= 'z' {
					runes[i] = r - 32
				}
				inWord = true
			} else {
				if r >= 'A' && r <= 'Z' {
					runes[i] = r + 32
				}
			}
		} else {
			inWord = false
		}
	}
	return string(runes)
}

func Handlequotes(s string) string {
	var result strings.Builder

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		str := line
		length := len(str)
		prevWasQuoteCluster := false
		inQuote := false
		var quoteString strings.Builder

		for j := 0; j < length; j++ {
			// contraction: single quote with non-space before and after
			if j > 0 && j < length-1 && str[j] == '\'' &&
				!unicode.IsSpace(rune(str[j-1])) && !unicode.IsSpace(rune(str[j+1])) {
				if inQuote {
					quoteString.WriteByte('\'')
				} else {
					result.WriteByte('\'')
				}
				prevWasQuoteCluster = false
				continue
			}

			if str[j] == '\'' {
				if !inQuote {
					// count consecutive quotes
					start := j
					for j < length && str[j] == '\'' {
						j++
					}
					count := j - start
					pairs := count / 2
					dangling := count % 2

					// space between clusters
					if prevWasQuoteCluster {
						result.WriteByte(' ')
					}

					// write pairs
					for p := 0; p < pairs; p++ {
						if p > 0 {
							result.WriteByte(' ')
						}
						result.WriteString("''")
					}

					// dangling starts non-empty or stays as dangling
					if dangling == 1 {
						if pairs > 0 {
							result.WriteByte(' ')
						}
						// look ahead: if not at end and next char is not quote -> it's a non-empty start
						if j < length && str[j] != '\'' {
							inQuote = true
							quoteString.Reset()
						} else {
							result.WriteByte('\'')
						}
					}

					prevWasQuoteCluster = true
					j-- // adjust
				} else {
					// closing non-empty
					trimmed := strings.TrimSpace(quoteString.String())
					result.WriteByte('\'')
					result.WriteString(trimmed)
					result.WriteByte('\'')
					inQuote = false
					quoteString.Reset()
					prevWasQuoteCluster = true
				}
			} else {
				// normal char
				if inQuote {
					quoteString.WriteByte(str[j])
				} else {
					result.WriteByte(str[j])
					prevWasQuoteCluster = false
				}
			}
		}

		// unclosed quote
		if inQuote {
			trimmed := strings.TrimSpace(quoteString.String())
			result.WriteByte('\'')
			result.WriteString(trimmed)
		}

		if i < len(lines)-1 {
			result.WriteByte('\n')
		}
	}

	return result.String()
}

func FormatText(input string) string {
	runes := []rune(input)
	punctuation := map[rune]bool{
		'.': true, ',': true, '!': true, '?': true, ':': true, ';': true,
	}

	var result []rune
	i := 0

	for i < len(runes) {
		// Handle newline correctly
		if runes[i] == '\n' {
			result = append(result, '\n')
			i++
			continue
		}
		// Skip multiple spaces
		if runes[i] == ' ' || runes[i] == '\t' || runes[i] == '\r' {
			result = append(result, ' ')
			for i < len(runes) && (runes[i] == ' ' || runes[i] == '\t' || runes[i] == '\r') {
				i++
			}
			continue
		}

		// Handle punctuation group
		if punctuation[runes[i]] {
			// Remove trailing space in result if present
			if len(result) > 0 && result[len(result)-1] == ' ' {
				result = result[:len(result)-1]
			}

			// Append full punctuation group (e.g. ..., !?, etc.)
			for i < len(runes) && punctuation[runes[i]] {
				result = append(result, runes[i])
				i++
			}

			// Add space if next is not punctuation or space
			if i < len(runes) && !(runes[i] == ' ' || runes[i] == '\t' || runes[i] == '\n' || runes[i] == '\r') && !punctuation[runes[i]] {
				result = append(result, ' ')
			}
			continue
		}

		// Regular character
		result = append(result, runes[i])
		i++
	}

	// Remove trailing space if any
	if len(result) > 0 && result[len(result)-1] == ' ' {
		result = result[:len(result)-1]
	}

	// Handle punctuations
	result = []rune(Handlequotes(string(result)))

	return string(result)
}
