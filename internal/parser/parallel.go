package parser

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/enohr/quake-log-parser/internal/model"
)

const (
	MAX_WORKERS = 5
)

// Struct of Worker Jobs
type Job struct {
	MatchNumber int
	Chunk       []string
}

// Struct of Worker Results
type Result struct {
	MatchNumber int
	Match       *model.Match
}

type Parallel struct {
}

func newParallel() *Parallel {
	return &Parallel{}
}

// Based on a received file, parse the log in parallel and returns a
// map with information about each match found
func (s *Parallel) Parse(file string) (map[string]*model.Match, error) {
	var wg sync.WaitGroup
	// Split the file into chunks of each match
	matchChunks, err := splitIntoChunks(file)

	if err != nil {
		return nil, err
	}

	jobs := make(chan Job, len(matchChunks))
	results := make(chan Result, len(matchChunks))

	slog.Info("Starting parallel parsing")

	// Starts the worker pool
	for i := 0; i < MAX_WORKERS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for job := range jobs {
				results <- Result{MatchNumber: job.MatchNumber, Match: parseChunk(job.Chunk)}
			}
		}()
	}

	// Send each chunk to the worker pool when it has space to process
	for matchNumber, chunk := range matchChunks {
		jobs <- Job{MatchNumber: matchNumber, Chunk: chunk}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	matches := make(map[string]*model.Match)

	// Receives the worker pool results
	for result := range results {
		slog.Info("Finish processing a match.", "Match number", result.MatchNumber)
		matchName := fmt.Sprintf("game_%d", result.MatchNumber)
		matches[matchName] = result.Match
	}

	slog.Info("Parallel parsing has ended")

	return matches, nil
}

// This function will receive a chunk and process it, returning the match data
func parseChunk(chunk []string) *model.Match {
	match := model.NewMatch()

	for _, line := range chunk {
		parseLine(line, match)
	}

	return match
}
