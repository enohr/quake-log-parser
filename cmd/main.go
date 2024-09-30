package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/enohr/quake-log-parser/internal/output"
	"github.com/enohr/quake-log-parser/internal/parser"
)

func main() {
	outputFilename := "output/output.json"

	parserType := flag.String("type", "", "The type of parsing (e.g., sequential, parallel)")
	inputFile := flag.String("input_file", "", "The path of file to be parsed")

	flag.Usage = func() {
		log.Printf("Usage %s -type [parallel|sequential] -input_file [file]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if !parser.ValidateParserType(*parserType) {
		log.Fatalf("Invalid parser type, check --help for instructions\n")
	}

	if *inputFile == "" {
		log.Fatalf("Invalid input file, check --help for instructions\n")
	}

	pt := parser.StringToParserType(*parserType)
	inputFilename := *inputFile

	p := parser.NewParser(pt)
	matches, err := p.Parse(inputFilename)

	if err != nil {
		log.Fatalf(err.Error())
	}

	if err := output.SaveOutput(matches, outputFilename); err != nil {
		log.Fatalf(err.Error())
	}

	slog.Info("Report saved on output/output.json file")

}
