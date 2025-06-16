package piscine

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func Format(InputFile string, OutputFile string) {
	// validate argument length atleast 5 or more
	if len(InputFile) < 5 {
		fmt.Println("Error: Invalid input filename. (Should be a .txt file)")
		fmt.Println("No content read. Aborting.")
		return
	}
	if len(OutputFile) < 5 {
		fmt.Println("Error: Invalid output filename. (Should be a .txt file)")
		fmt.Println("No content read. Aborting.")
		return
	}

	// validates if input file is .txt
	if InputFile[len(InputFile)-4:] != ".txt" {
		fmt.Println("Error: Invalid input filename. (Should be a .txt file)")
		fmt.Println("No content read. Aborting.")
		return
	}
	// validates if output file is .txt
	if OutputFile[len(OutputFile)-4:] != ".txt" {
		fmt.Println("Error: Invalid output filename. (Should be a .txt file)")
		fmt.Println("No content read. Aborting.")
		return
	}

	// display error if the file doesnt exist
	_, err := os.Stat(InputFile)
	if err != nil || os.IsNotExist(err) {
		fmt.Println("Error: " + InputFile + ": The system cannot find the file specified.")
		fmt.Println("No content read. Aborting.")
		return
	}

	text := Gettext(InputFile)
	// display error if file is empty
	if text == "" {
		fmt.Println("Error: Empty input file.\nNo content read. Aborting.")
		return
	}
	// exit if there are unsupported ascii characters7
	for _, v := range text {
		if (v < 32 && v != 9 && v != 10 && v != 13) || v > 126 {
			fmt.Println("Error unsupported input: illegal characters found.\nPlease enter valid ascii characters.\nNo content read. Aborting.")
			return
		}
	}

	words := DetectCase(TokenizeInput(text))
	// fmt.Println(fmt.Sprintf("%#v\n", TokenizeInput(text)))
	if words == nil {
		fmt.Println("No content read. Aborting.")
		return
	}

	str := FormatText((JoinStrings(words, " ")))

	WriteStringToFile(OutputFile, str)
}

func Gettext(input string) string {
	filename := input

	// if only filename is given then it turns it into a relative path
	if len(filename) > 0 && filename[0] != '.' && filename[0] != '/' {
		filename = "./" + filename
	}

	// read from the file
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	content := ""

	for {
		chunk, err := reader.ReadByte()
		if err != nil {
			break
		}
		content += string(chunk)
	}

	return content
}

// func SplitWhiteSpaces(s string) []string {
// 	runes := []rune(s)
// 	result := []string{}
// 	word := ""
// 	length := len(runes)

// 	for i := 0; i < length; i++ {
// 		char := runes[i]

// 		// Handle whitespace (space, tab, carriage return)
// 		if char == ' ' || char == '\t' || char == '\r' {
// 			if word != "" {
// 				result = append(result, word)
// 				word = ""
// 			}
// 			continue
// 		}

// 		// Handle newline as a separate word
// 		if char == '\n' {
// 			if word != "" {
// 				result = append(result, word)
// 				word = ""
// 			}
// 			result = append(result, "\n")
// 			continue
// 		}

// 		// Regular character accumulation
// 		word += string(char)
// 	}

// 	if word != "" {
// 		result = append(result, word)
// 	}

// 	return result
// }

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

// func isValidSurrounding(s string) bool {
// 	s = strings.TrimSpace(s)
// 	if len(s) < 4 || s[0] != '(' || s[len(s)-1] != ')' {
// 		return false
// 	}

// 	// Ensure no touching non-space chars outside the parentheses
// 	// ( e.g, "fox(cap)" should be invalid )
// 	// This is handled earlier by tokenization (splitwhitespace)

// 	inner := strings.TrimSpace(s[1 : len(s)-1])
// 	innerLower := strings.ToLower(strings.Split(inner, ",")[0])
// 	return innerLower == "up" || innerLower == "low" || innerLower == "cap"
// }

// func IsSpecialCase(s string) bool {
// 	if _, _, valid := ParseSpecialKeyword(s); !valid {
// 		return false
// 	}
// 	return true
// }

func StartsWithVowel(s string) bool {
	if s == "" {
		return false
	}
	first := s[0]
	return first == 'A' || first == 'E' || first == 'I' || first == 'O' || first == 'U' ||
		first == 'a' || first == 'e' || first == 'i' || first == 'o' || first == 'u'
}

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

// func isPunctuation(r rune) bool {
// 	// Treat common punctuation marks as quote closers
// 	punctuations := []rune{'.', ',', '!', '?', ':', ';', '-', '(', ')', '"'}
// 	for _, p := range punctuations {
// 		if r == p {
// 			return true
// 		}
// 	}
// 	return false
// }

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
					result[j] = Capitalize(result[j])
				case "hex":
					trimmed := strings.TrimSpace(result[j])
					if trimmed != "" {
						if converted := HexToDec(trimmed); valid {
							result[j] = converted
						}
					}
				case "bin":
					trimmed := strings.TrimSpace(result[j])
					if trimmed != "" {
						if converted := BinToDec(trimmed); valid {
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

func WriteStringToFile(filename string, data string) error {
	// Create or overwrite the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the string data to the file
	_, err = file.WriteString(data)
	return err
}
