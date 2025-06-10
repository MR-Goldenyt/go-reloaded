package piscine

import (
	"bufio"
	"fmt"
	"os"
)

func Gettext(input string) string {
	filename := input

	// if only filename is given then it turns it into a relative path
	if len(filename) > 0 && filename[0] != '.' && filename[0] != '/' {
		filename = "./" + filename
	}

	// read from the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
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

func SplitWhiteSpaces(s string) []string {
	runes := []rune(s)
	result := []string{}
	word := ""
	length := len(runes)

	for i := 0; i < length; i++ {
		char := runes[i]

		// Handle whitespace (space, tab, carriage return)
		if char == ' ' || char == '\t' || char == '\r' {
			if word != "" {
				result = append(result, word)
				word = ""
			}
			continue
		}

		// Handle newline as a separate word
		if char == '\n' {
			if word != "" {
				result = append(result, word)
				word = ""
			}
			result = append(result, "\n")
			continue
		}

		// Handle special cases like (cap), (low), (cap, 3), etc.
		if char == '(' {
			start := i
			j := i
			for j < length && runes[j] != ')' {
				j++
			}
			if j < length && runes[j] == ')' {
				group := string(runes[start : j+1])
				if IsSpecialCase(group) {
					if word != "" {
						result = append(result, word)
						word = ""
					}
					result = append(result, group)
					i = j // Move i to closing parenthesis so next loop iteration starts after it
					continue
				}
			}
		}

		// Regular character accumulation
		word += string(char)
	}

	if word != "" {
		result = append(result, word)
	}

	return result
}

func IsSpecialCase(s string) bool {
	// Check exact special cases without number
	if s == "(low)" || s == "(up)" || s == "(cap)" {
		return true
	}

	// Check for special cases with numbers
	if len(s) >= 7 && (s[:5] == "(low," || s[:5] == "(cap," || s[:4] == "(up,") {
		// Ensure last character is ')'
		if s[len(s)-1] == ')' {
			// Look for space after comma
			for i := 0; i < len(s)-1; i++ {
				if s[i] == ' ' {
					digits := s[i+1 : len(s)-1]
					// Check all digits in substring
					for _, ch := range digits {
						if ch < '0' || ch > '9' {
							return false
						}
					}
					return true
				}
			}
		}
	}

	return false
}

func parseSpecialCase(s string) (string, int) {
	// No number
	if s == "(low)" {
		return "low", 1
	} else if s == "(up)" {
		return "up", 1
	} else if s == "(cap)" {
		return "cap", 1
	}

	// With number
	var caseType string
	var start int

	if s[:5] == "(low," {
		caseType = "low"
		start = 5
	} else if s[:5] == "(cap," {
		caseType = "cap"
		start = 5
	} else if s[:4] == "(up," {
		caseType = "up"
		start = 4
	}

	// Find space
	space := -1
	for i := start; i < len(s)-1; i++ {
		if s[i] == ' ' {
			space = i
			break
		}
	}

	if space == -1 {
		return caseType, 1
	}

	// Extract number
	num := 0
	for i := space + 1; i < len(s)-1; i++ {
		ch := s[i]
		if ch >= '0' && ch <= '9' {
			num = num*10 + int(ch-'0')
		} else {
			return caseType, 1
		}
	}

	if num <= 0 {
		num = 1
	}

	return caseType, num
}

func StartsWithVowel(s string) bool {
	if s == "" {
		return false
	}
	first := s[0]
	return first == 'A' || first == 'E' || first == 'I' || first == 'O' || first == 'U' ||
		first == 'a' || first == 'e' || first == 'i' || first == 'o' || first == 'u'
}

func IntToString(n int) string {
	if n == 0 {
		return "0"
	}

	digits := []byte{}

	// Extract digits in reverse order
	for n > 0 {
		digit := n % 10
		digits = append(digits, byte(digit)+'0')
		n /= 10
	}

	// Reverse the slice to get correct order
	i := 0
	j := len(digits) - 1

	for i < j {
		digits[i], digits[j] = digits[j], digits[i]
		i++
		j--
	}

	return string(digits)
}

func HexToDecimalString(s string) string {
	result := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		var val int
		if c >= '0' && c <= '9' {
			val = int(c - '0')
		} else if c >= 'a' && c <= 'f' {
			val = int(c - 'a' + 10)
		} else if c >= 'A' && c <= 'F' {
			val = int(c - 'A' + 10)
		} else {
			break // invalid char
		}
		result = result*16 + val
	}
	return IntToString(result)
}

func BinToDecimalString(s string) string {
	result := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '0' {
			result = result * 2
		} else if c == '1' {
			result = result*2 + 1
		} else {
			break // invalid char
		}
	}
	return IntToString(result)
}

func ToLower(s string) string {
	str := ""
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			str += string(v + 32)
		} else {
			str += string(v)
		}
	}
	return str
}

func ToUpper(s string) string {
	str := ""
	for _, v := range s {
		if v >= 97 && v <= 122 {
			str += string(v - 32)
		} else {
			str += string(v)
		}
	}
	return str
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

func FormatText(input string) string {
	runes := []rune(input)
	punctuation := map[rune]bool{
		'.': true, ',': true, '!': true, '?': true, ':': true, ';': true,
	}

	var result []rune
	i := 0
	length := len(runes)

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

		// Handle single-quoted expressions
		if runes[i] == '\'' {
			i++
			contentStart := i
			for i < length && runes[i] != '\'' {
				i++
			}
			contentEnd := i
			if i < length && runes[i] == '\'' {
				// Trim leading spaces
				for contentStart < contentEnd && runes[contentStart] == ' ' {
					contentStart++
				}
				// Trim trailing spaces
				for contentEnd > contentStart && runes[contentEnd-1] == ' ' {
					contentEnd--
				}
				result = append(result, '\'')
				for j := contentStart; j < contentEnd; j++ {
					result = append(result, runes[j])
				}
				result = append(result, '\'')
				i++ // move past closing quote
				continue
			} else {
				// No closing quote found, treat as normal character
				result = append(result, '\'')
				continue
			}
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

	return string(result)
}

func DetectCase(words []string) []string {
	result := []string{}
	length := len(words)

	for i := 0; i < length; i++ {
		word := words[i]

		// Handle special cases
		if IsSpecialCase(word) {
			caseType, count := parseSpecialCase(word)

			// Apply transformation to previous `count` words
			n := len(result)
			start := n - count
			if start < 0 {
				start = 0
			}

			for j := start; j < n; j++ {
				switch caseType {
				case "low":
					result[j] = ToLower(result[j])
				case "up":
					result[j] = ToUpper(result[j])
				case "cap":
					result[j] = Capitalize(result[j])
				}
			}
			// Do not add the special case to the result
			continue
		}
		// Handle Hex to decimal conversion
		if word == "(hex)" {
			if len(result) > 0 {
				result[len(result)-1] = HexToDecimalString(result[len(result)-1])
			}
			continue
		}
		// Handle binary to decimal conversion
		if word == "(bin)" {
			if len(result) > 0 {
				result[len(result)-1] = BinToDecimalString(result[len(result)-1])
			}
			continue
		}
		// Handle "A" or "An" and fix based on next word's first letter
		if (word == "A" || word == "An" || word == "AN" || word == "a" || word == "an") && i+1 < length {
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
	}

	return result
}

func JoinStrings(slice []string, sep string) string {
	result := ""
	for i := 0; i < len(slice); i++ {
		result += slice[i]

		if i < len(slice)-1 && slice[i] != "\n" && slice[i+1] != "\n" {
			result += sep
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

func Format(InputFile string, OutputFile string) {
	text := Gettext(InputFile)
	if text == "" {
		fmt.Println("No content read. Aborting.")
		return
	}

	words := DetectCase(SplitWhiteSpaces(Gettext(InputFile)))
	str := FormatText(JoinStrings(words, " "))
	WriteStringToFile(OutputFile, str)
}
