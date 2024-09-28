package parser

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"slices"

	"github.com/enohr/quake-log-parser/internal/model"
)

const (
	INIT_GAME_REGEX = `InitGame:\s.*`
	KILLED_REGEX    = `\d+:\d+ Kill: \d+ \d+ \d+: (?P<killer>.*) killed (?P<victim>.*) by (?P<mean>.*)`
)

type Sequential struct {
}

func newSequential() *Sequential {
	return &Sequential{}
}

func (s *Sequential) Parse(file string) (map[string]*model.Match, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	slog.Info("Starting parsing log")

	scanner := bufio.NewScanner(f)
	matches, err := processGame(scanner)

	slog.Info("Log parsing has finished")

	return matches, err
}

func processGame(scanner *bufio.Scanner) (map[string]*model.Match, error) {
	var match *model.Match
	matches := make(map[string]*model.Match)

	initGameRegex := regexp.MustCompile(INIT_GAME_REGEX)
	killedRegex := regexp.MustCompile(KILLED_REGEX)

	totalGames := 0

	for scanner.Scan() {
		line := scanner.Text()

		if initGameRegex.MatchString(line) {
			if match != nil {
				name := fmt.Sprintf("game-%d", totalGames)
				matches[name] = match
				totalGames++
			}
			match = model.NewMatch()
		} else if m := killedRegex.FindStringSubmatch(line); m != nil {
			killer, victim, _ := m[1], m[2], m[3]

			if killer == "<world>" {
				match.Kills[victim]--
			} else {
				match.Kills[killer]++
			}

			if !slices.Contains(match.Players, killer) {
				match.Players = append(match.Players, killer)
			}

			if !slices.Contains(match.Players, victim) {
				match.Players = append(match.Players, victim)
			}

			match.TotalKills++
		}

	}
	if match != nil {
		name := fmt.Sprintf("game-%d", totalGames)
		matches[name] = match
		totalGames++
	}

	return matches, nil
}
