package main

import (
	"log"

	"github.com/enohr/quake-log-parser/internal/output"
	"github.com/enohr/quake-log-parser/internal/parser"
)

func main() {
	inputFilename := "input/quake.log"
	outputFilename := "output/output.json"

	p := parser.NewParser()
	matches, err := p.Parse(inputFilename)

	if err != nil {
		log.Println(err)
		return
	}

	if err := output.SaveOutput(matches, outputFilename); err != nil {
		log.Println(err)
		return
	}

}
