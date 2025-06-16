package process

import (
	"fmt"
	"os"

	"textformatter/input"
	"textformatter/output"
	"textformatter/transform"
	"textformatter/utilities"
)

func Format(InputFile string, OutputFile string) {
	// validate argument length atleast 5 or more
	if len(InputFile) < 5 {
		fmt.Println("Error: Invalid input filename. (Should be a .txt file)")
		if len(OutputFile) < 5 {
			fmt.Println("Error: Invalid output filename. (Should be a .txt file)")
		}
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

	text := input.Gettext(InputFile)
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

	words := utilities.DetectCase(utilities.TokenizeInput(text))
	// fmt.Println(fmt.Sprintf("%#v\n", TokenizeInput(text)))
	if words == nil {
		fmt.Println("No content read. Aborting.")
		return
	}

	str := transform.FormatText((utilities.JoinStrings(words, " ")))

	output.WriteStringToFile(OutputFile, str)
}
