package parser

import "github.com/enohr/quake-log-parser/internal/model"

type Parser interface {
	Parse(file string) (map[string]*model.Match, error)
}

func NewParser() Parser {

	// TODO: Receive the type of parser and return it
	return newSequential()
}
