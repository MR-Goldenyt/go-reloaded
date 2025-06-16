package utilities

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"textformatter/transform"
)

func TokenizeInput(input string) []string {
	// Pattern explanation:
	// (?i) - case-insensitive
	// (\([ \t\r]*(?:up|low|cap|bin|hex)[ \t\r]*(?:,[ \t\r]*-?\d+)?[ \t\r]*\)+) - special keywords with optional signed number and multiple trailing )
	// |(\n) - newline tokenized as "\n"
	// |([^\s]+) - any non-whitespace sequence
	pattern := `(?i)(\([ \t\r]*(?:up|low|cap|bin|hex)[ \t\r]*(?:,[ \t\r]*-?\d+)?[ \t\r]*\)+)|(\n)|([^\s]+)`

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Regex compilation error:", err)
		return nil
	}

	matches := re.FindAllString(input, -1)
	return matches
}

func ParseSpecialKeyword(s string) (string, int, bool) {
	s = strings.TrimSpace(s)

	// Must start with '(' and end with ')', and only one pair
	if !strings.HasPrefix(s, "(") || !strings.HasSuffix(s, ")") {
		return "", 1, false
	}

	// Ensure no outer text like (((cap)))
	if strings.Count(s, "(") != 1 || strings.Count(s, ")") != 1 {
		return "", 1, false
	}

	inner := strings.TrimSpace(s[1 : len(s)-1])
	parts := strings.Split(strings.ReplaceAll(inner, " ", ""), ",")

	if len(parts) < 1 {
		return "", 1, false
	}

	keyword := strings.ToLower(parts[0])

	switch keyword {
	case "up", "low", "cap":
		count := 1
		if len(parts) == 2 {
			n, ok := strconv.Atoi(parts[1])
			if ok == nil {
				count = n
			}
		}
		return keyword, count, true
	case "hex", "bin":
		if len(parts) == 1 {
			return keyword, 1, true
		}
	}
	return "", 1, false
}

func StartsWithVowel(s string) bool {
	if s == "" {
		return false
	}
	first := s[0]
	return first == 'A' || first == 'E' || first == 'I' || first == 'O' || first == 'U' ||
		first == 'a' || first == 'e' || first == 'i' || first == 'o' || first == 'u'
}

func DetectCase(words []string) []string {
	var result []string
	lineStart := 0

	for i := 0; i < len(words); i++ {
		word := words[i]
		keyword, count, valid := ParseSpecialKeyword(word)
		if valid && count > 0 {
			start := len(result) - count
			if start < lineStart {
				start = lineStart
			}
			for j := start; j < len(result); j++ {
				switch keyword {
				case "up":
					result[j] = strings.ToUpper(result[j])
				case "low":
					result[j] = strings.ToLower(result[j])
				case "cap":
					result[j] = transform.Capitalize(result[j])
				case "hex":
					trimmed := strings.TrimSpace(result[j])
					if trimmed != "" {
						if converted := transform.HexToDec(trimmed); valid {
							result[j] = converted
						}
					}
				case "bin":
					trimmed := strings.TrimSpace(result[j])
					if trimmed != "" {
						if converted := transform.BinToDec(trimmed); valid {
							result[j] = converted
						}
					}
				}
			}
			continue // skip appending keyword itself
		} else if valid && count < 1 {
			fmt.Println("Error: [(" + keyword + ", " + strconv.Itoa(count) + ")]: " + strconv.Itoa(count) + " is an invalid number\nPlease enter a number equal to or greater than 1")
			continue
		}

		// Handle "A" or "An" and fix based on next word's first letter
		if (word == "A" || word == "An" || word == "AN" || word == "a" || word == "an") && i+1 < len(words) {
			nextWord := words[i+1]
			if StartsWithVowel(nextWord) && word == "A" {
				word = "An"
			} else if !StartsWithVowel(nextWord) && word == "An" || !StartsWithVowel(nextWord) && word == "AN" {
				word = "A"
			} else if StartsWithVowel(nextWord) && word == "a" {
				word = "an"
			} else if !StartsWithVowel(nextWord) && word == "an" {
				word = "a"
			}
		}

		result = append(result, word)

		// Update lineStart index after newline token
		if word == "\n" {
			lineStart = len(result)
		}
	}
	return result
}

func JoinStrings(slice []string, sep string) string {
	result := ""
	for i := 0; i < len(slice); i++ {
		result += slice[i]

		if i < len(slice)-1 {
			curr := slice[i]
			next := slice[i+1]

			// Skip adding separator if newline or matching parens
			if curr != "\n" && next != "\n" && !(strings.HasSuffix(curr, "(") && strings.HasPrefix(next, ")")) {
				result += sep
			}
		}
	}
	return result
}
