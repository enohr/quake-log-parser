package parser

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"

	"github.com/enohr/quake-log-parser/internal/model"
)

const (
	INIT_GAME_REGEX         = `InitGame:\s.*`
	KILLED_REGEX            = `Kill: (?P<killer_id>\d+) (?P<victim_id>\d+) (?P<mean_id>\d+)`
	JOIN_GAME_REGEX         = `ClientConnect: (?P<player_id>\d+)`
	USER_INFO_CHANGED_REGEX = `ClientUserinfoChanged: (?P<player_id>\d+) n\\(?P<player>.*?)\\`
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
	matches, err := processMatches(scanner)

	slog.Info("Log parsing has finished")

	return matches, err
}

func processMatches(scanner *bufio.Scanner) (map[string]*model.Match, error) {
	var match *model.Match
	matches := make(map[string]*model.Match)

	initGameRegex := regexp.MustCompile(INIT_GAME_REGEX)
	joinGameRegex := regexp.MustCompile(JOIN_GAME_REGEX)
	userInfoChangedRegex := regexp.MustCompile(USER_INFO_CHANGED_REGEX)
	killedRegex := regexp.MustCompile(KILLED_REGEX)

	totalGames := 0

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case initGameRegex.MatchString(line):
			if match != nil {
				name := fmt.Sprintf("game-%d", totalGames)
				matches[name] = match
				totalGames++
			}
			match = model.NewMatch()

		case joinGameRegex.MatchString(line):
			m := joinGameRegex.FindStringSubmatch(line)
			player := m[1]

			playerID, err := strconv.Atoi(player)
			if err != nil {
				continue
			}
			match.AddPlayer(playerID)

		case userInfoChangedRegex.MatchString(line):
			m := userInfoChangedRegex.FindStringSubmatch(line)
			player, playerName := m[1], m[2]

			playerID, err := strconv.Atoi(player)
			if err != nil {
				continue
			}
			match.UpdateUserInfo(playerID, playerName)

		case killedRegex.MatchString(line):
			m := killedRegex.FindStringSubmatch(line)

			killer, victim, mean := m[1], m[2], m[3]

			killerID, err := strconv.Atoi(killer)

			if err != nil {
				continue
			}

			victimID, err := strconv.Atoi(victim)

			if err != nil {
				continue
			}

			meanID, err := strconv.Atoi(mean)

			if err != nil {
				continue
			}

			match.ProcessKill(killerID, victimID, meanID)
		}
	}
	if match != nil {
		name := fmt.Sprintf("game-%d", totalGames)
		matches[name] = match
		totalGames++
	}

	return matches, nil
}
