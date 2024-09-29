package parser

import (
	"bufio"
	"fmt"
	"os"
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
		matchName := fmt.Sprintf("game_%d", result.MatchNumber)
		matches[matchName] = result.Match
	}

	return matches, nil
}

func parseChunk(chunk []string) *model.Match {
	match := model.NewMatch()

	for _, line := range chunk {
		processLine(line, match)
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
