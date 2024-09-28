package util

import (
	"encoding/json"
	"os"

	"github.com/enohr/quake-log-parser/internal/model"
)

func SaveJSONOutput(matches map[string]model.MatchJSON, outputFilename string) error {

	content, err := json.MarshalIndent(matches, "", " ")

	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}

	_, err = outputFile.Write(content)

	if err != nil {
		return err
	}

	return nil
}
