package util

import (
	"os"
	"path/filepath"
)

func SaveToFile(filename string, content []byte) error {
	dir := filepath.Dir(filename)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	outputFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.Write(content)

	if err != nil {
		return err
	}

	return nil
}
