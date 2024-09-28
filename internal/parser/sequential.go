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
	WORLD_ID                = 1022
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

		if initGameRegex.MatchString(line) {
			if match != nil {
				name := fmt.Sprintf("game-%d", totalGames)
				matches[name] = match
				totalGames++
			}
			match = model.NewMatch()
		} else if m := joinGameRegex.FindStringSubmatch(line); m != nil {
			player := m[1]

			playerID, err := strconv.Atoi(player)
			if err != nil {
				continue
			}

			if _, ok := match.Players[playerID]; !ok {
				match.Players[playerID] = ""
			}

		} else if m := userInfoChangedRegex.FindStringSubmatch(line); m != nil {
			player, playerName := m[1], m[2]

			playerID, err := strconv.Atoi(player)
			if err != nil {
				continue
			}

			if _, ok := match.Players[playerID]; ok {
				match.Players[playerID] = playerName
			}

		} else if m := killedRegex.FindStringSubmatch(line); m != nil {
			killer, victim, _ := m[1], m[2], m[3]

			killerID, err := strconv.Atoi(killer)

			if err != nil {
				continue
			}

			victimID, err := strconv.Atoi(victim)

			if err != nil {
				continue
			}

			if killerID == WORLD_ID {
				match.Kills[victimID]--
			} else if killerID != victimID {
				match.Kills[killerID]++
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
