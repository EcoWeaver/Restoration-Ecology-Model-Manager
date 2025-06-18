// Preprocessing Backup SQL File

package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please also enter filename as input.")
	}
	inputFile := os.Args[1]
	content, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error when reading the file: %v", err)
	}

	data := string(content)
	lines := strings.Split(data, "\n")
	var outputLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "INSERT INTO public.") {
			cleanLine := strings.Replace(trimmed, "INSERT INTO public.", "INSERT INTO ", 1)
			outputLines = append(outputLines, cleanLine)
		}
	}
	outputData := strings.Join(outputLines, "\n") + "\n"

	tmpFile := inputFile + ".tmp"
	err = os.WriteFile(tmpFile, []byte(outputData), 0644)
	if err != nil {
		log.Fatalf("Error when writing into temporary file: %v", err)
	}

	err = os.Remove(inputFile)
	if err != nil {
		log.Fatalf("Error when deleting the input file:  %v", err)
	}

	err = os.Rename(tmpFile, inputFile)
	if err != nil {
		log.Fatalf("Error when renaming the temporary file %v", err)
	}
}
