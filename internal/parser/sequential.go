package parser

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/enohr/quake-log-parser/internal/model"
)

type Sequential struct {
}

func newSequential() *Sequential {
	return &Sequential{}
}

// Based on a received file, parse the log in serial and returns a
// map with information about each match found
func (s *Sequential) Parse(file string) (map[string]*model.Match, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	slog.Info("Starting sequential parsing")

	scanner := bufio.NewScanner(f)
	matches, err := parseMatches(scanner)

	slog.Info("Sequential parsing has ended")

	return matches, err
}

func parseMatches(scanner *bufio.Scanner) (map[string]*model.Match, error) {
	var match *model.Match
	matches := make(map[string]*model.Match)

	matchNumber := 0

	for scanner.Scan() {
		line := scanner.Text()

		newMatch, err := parseLine(line, match)
		if err != nil {
			return nil, err
		}

		// If has a current match, add it to the result map
		if newMatch {
			if match != nil {
				slog.Info("Finish processing a match.", "Match number", matchNumber)
				name := fmt.Sprintf("game_%d", matchNumber)
				matches[name] = match
				matchNumber++

			}
			match = model.NewMatch()
		}

	}

	// Adds the last match to the result map
	if match != nil {
		slog.Info("Finish processing a match.", "Match number", matchNumber)
		name := fmt.Sprintf("game_%d", matchNumber)
		matches[name] = match
		matchNumber++
	}

	return matches, nil
}
