package main

import (
	"os"

	"piscine"
)

func main() {
	InputFile := os.Args[1]
	OutputFile := os.Args[2]

	piscine.Format(InputFile, OutputFile)
}
