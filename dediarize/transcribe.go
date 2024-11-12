package dediarize

import (
	"fmt"
	"os"
)

func TranscribeJsonFile(inputFile, outputFile string) error {
	jsonFile, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", inputFile, err)
	}
	defer jsonFile.Close()

	fo, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", outputFile, err)
	}
	defer fo.Close()

	w, err := NewFileWriter(fo)
	if err != nil {
		return fmt.Errorf("failed to create file writer: %w", err)
	}

	p, err := NewParser(w)
	if err != nil {
		return fmt.Errorf("failed to create parser: %w", err)
	}

	if err := p.Parse(jsonFile); err != nil {
		return fmt.Errorf("failed to parse json: %w", err)
	}
	return nil
}
