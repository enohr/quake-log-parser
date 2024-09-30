package parser

import (
	"bufio"
	"os"
	"strconv"

	"github.com/enohr/quake-log-parser/internal/model"
)

// Receive each line, parses it and update the match information
// Returns true if a new match has started
func parseLine(line string, match *model.Match) (bool, error) {
	switch {
	case initGameRegex.MatchString(line): // Checks if a new match has started
		return true, nil

	case joinGameRegex.MatchString(line): // Checks if user has joined the game
		m := joinGameRegex.FindStringSubmatch(line)
		player := m[1]

		playerID, err := strconv.Atoi(player)
		if err != nil {
			return false, err
		}
		match.AddPlayer(playerID)

	case disconnectGameRegex.MatchString(line): // Checks if user has disconnected
		m := disconnectGameRegex.FindStringSubmatch(line)
		player := m[1]

		playerID, err := strconv.Atoi(player)
		if err != nil {
			return false, err
		}
		match.DisconnectPlayer(playerID)

	case userInfoChangedRegex.MatchString(line): // Check if the user changed some info
		m := userInfoChangedRegex.FindStringSubmatch(line)
		player, playerName := m[1], m[2]

		playerID, err := strconv.Atoi(player)
		if err != nil {
			return false, err
		}
		match.UpdateUserInfo(playerID, playerName)

	case killedRegex.MatchString(line): // Check if theres a kill event
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

// Receives a file name and returns a matrix of each match found
func splitIntoChunks(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	chunk := make([]string, 0)
	chunks := make([][]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if initGameRegex.MatchString(line) {
			if len(chunk) > 0 {
				chunks = append(chunks, chunk)
			}

			chunk = make([]string, 0)
			chunk = append(chunk, line)
		} else {
			if len(chunk) > 0 {
				chunk = append(chunk, line)
			}
		}
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}
