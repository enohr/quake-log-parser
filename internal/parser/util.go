package parser

import (
	"strconv"

	"github.com/enohr/quake-log-parser/internal/model"
)

func processLine(line string, match *model.Match) (bool, error) {
	switch {
	case initGameRegex.MatchString(line):
		return true, nil

	case joinGameRegex.MatchString(line):
		m := joinGameRegex.FindStringSubmatch(line)
		player := m[1]

		playerID, err := strconv.Atoi(player)
		if err != nil {
			return false, err
		}
		match.AddPlayer(playerID)

	case disconnectGameRegex.MatchString(line):
		m := disconnectGameRegex.FindStringSubmatch(line)
		player := m[1]

		playerID, err := strconv.Atoi(player)
		if err != nil {
			return false, err
		}
		match.DisconnectPlayer(playerID)

	case userInfoChangedRegex.MatchString(line):
		m := userInfoChangedRegex.FindStringSubmatch(line)
		player, playerName := m[1], m[2]

		playerID, err := strconv.Atoi(player)
		if err != nil {
			return false, err
		}
		match.UpdateUserInfo(playerID, playerName)

	case killedRegex.MatchString(line):
		m := killedRegex.FindStringSubmatch(line)

		killer, victim, mean := m[1], m[2], m[3]

		killerID, err := strconv.Atoi(killer)

		if err != nil {
			return false, err
		}

		victimID, err := strconv.Atoi(victim)

		if err != nil {
			return false, err
		}

		meanID, err := strconv.Atoi(mean)

		if err != nil {
			return false, err
		}

		match.ProcessKill(killerID, victimID, meanID)
	}
	return false, nil
}
