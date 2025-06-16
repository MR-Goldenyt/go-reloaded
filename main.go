package main

import (
	"fmt"
	"os"
	"strconv"

	process "textformatter/processor"
)

func main() {
	if len(os.Args) > 3 || len(os.Args) < 3 {
		fmt.Println("Error: Incorrect number of arguments. Expected [2], but received [" + strconv.Itoa(len(os.Args)-1) + "].\nPlease check the usage instructions.")
		return
	}
	InputFile := os.Args[1]
	OutputFile := os.Args[2]

	process.Format(InputFile, OutputFile)
}
