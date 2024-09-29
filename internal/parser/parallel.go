package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/enohr/quake-log-parser/internal/model"
)

const (
	MAX_WORKERS = 5
)

type Job struct {
	MatchNumber int
	Chunk       []string
}

type Result struct {
	MatchNumber int
	Match       *model.Match
}

type Parallel struct {
}

func newParallel() *Parallel {
	return &Parallel{}
}

func (s *Parallel) Parse(file string) (map[string]*model.Match, error) {
	var wg sync.WaitGroup
	matchChunks, err := splitIntoChunks(file)

	if err != nil {
		return nil, err
	}

	jobs := make(chan Job, len(matchChunks))
	results := make(chan Result, len(matchChunks))

	for i := 0; i < MAX_WORKERS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for job := range jobs {
				results <- Result{MatchNumber: job.MatchNumber, Match: parseChunk(job.Chunk)}
			}
		}()
	}

	for matchNumber, chunk := range matchChunks {
		jobs <- Job{MatchNumber: matchNumber, Chunk: chunk}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	matches := make(map[string]*model.Match)

	for result := range results {
		matchName := fmt.Sprintf("game-%d", result.MatchNumber)
		matches[matchName] = result.Match
	}

	return matches, nil
}

func parseChunk(chunk []string) *model.Match {
	match := model.NewMatch()

	for _, line := range chunk {
		switch {
		case joinGameRegex.MatchString(line):
			m := joinGameRegex.FindStringSubmatch(line)
			player := m[1]

			playerID, err := strconv.Atoi(player)
			if err != nil {
				continue
			}
			match.AddPlayer(playerID)

		case disconnectGameRegex.MatchString(line):
			m := disconnectGameRegex.FindStringSubmatch(line)
			player := m[1]

			playerID, err := strconv.Atoi(player)
			if err != nil {
				continue
			}
			match.DisconnectPlayer(playerID)

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
	return match

}

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
