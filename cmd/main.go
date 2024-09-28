package main

import (
	"log"

	"github.com/enohr/quake-log-parser/internal/model"
	"github.com/enohr/quake-log-parser/internal/parser"
	"github.com/enohr/quake-log-parser/util"
)

func main() {
	inputFilename := "input/quake2.log"
	outputFilename := "output/output.json"

	p := parser.NewParser()
	matches, err := p.Parse(inputFilename)

	if err != nil {
		log.Println(err)
		return
	}

	matchesJSON := make(map[string]model.MatchJSON)
	for k, v := range matches {
		matchesJSON[k] = v.ToMatchJSON()
	}

	if err := util.SaveJSONOutput(matchesJSON, outputFilename); err != nil {
		log.Println(err)
		return
	}

}
