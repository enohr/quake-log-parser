package util

import (
	"hash/fnv"
	"log/slog"
	"os"
	"path/filepath"
)

// Save an array of bytes to the given filename
func SaveToFile(filename string, content []byte) error {
	dir := filepath.Dir(filename)

	slog.Info("Creating output directory", "directory", dir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	slog.Info("Creating output file", "file", filename)
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

// Generates a hash of given username
func GenerateNameFNVHash(name string) int {
	h := fnv.New64a()
	h.Write([]byte(name))
	return int(h.Sum64())
}
