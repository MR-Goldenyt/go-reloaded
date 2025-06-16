package input

import (
	"bufio"
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
