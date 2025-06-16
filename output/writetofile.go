package output

import "os"

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
