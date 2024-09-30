package parser

import (
	"github.com/enohr/quake-log-parser/internal/model"
)

type Parser interface {
	Parse(file string) (map[string]*model.Match, error)
}

type ParserType int

// Creates a parser based on ParserType
func NewParser(pt ParserType) Parser {
	switch pt {
	case ParallelParser:
		return newParallel()
	case SequentialParser:
		return newSequential()

	}
	return nil
}

const (
	UnknownParser ParserType = iota
	ParallelParser
	SequentialParser
)

var parserTypeStrings = map[string]ParserType{
	"":           UnknownParser,
	"parallel":   ParallelParser,
	"sequential": SequentialParser,
}

// Returns the ParserType based on received String
func StringToParserType(parserType string) ParserType {
	if pt, ok := parserTypeStrings[parserType]; ok {
		return pt
	}
	return UnknownParser
}

// Validate received parser type
func ValidateParserType(parserType string) bool {
	pt := StringToParserType(parserType)
	return pt != UnknownParser
}
