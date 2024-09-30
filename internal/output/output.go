package output

import (
	"encoding/json"
	"log/slog"

	"github.com/enohr/quake-log-parser/internal/model"
	"github.com/enohr/quake-log-parser/util"
)

func SaveOutput(matches map[string]*model.Match, filename string) error {
	slog.Info("Transforming matches to output format")
	matchesJSON, err := matchToJSON(matches)

	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(matchesJSON, "", " ")

	if err != nil {
		return err
	}

	return util.SaveToFile(filename, j)
}

func matchToJSON(matches map[string]*model.Match) (map[string]model.MatchJSON, error) {
	matchesJSON := make(map[string]model.MatchJSON)

	// Transform each match to the output format
	for k, v := range matches {
		matchesJSON[k] = v.ToMatchJSON()
	}

	return matchesJSON, nil
}
