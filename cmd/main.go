package main

import (
	"log"

	"github.com/enohr/quake-log-parser/internal/parser"
)

func main() {
	file := "input/quake.log"

	p := parser.NewParser()

	games, err := p.Parse(file)

	if err != nil {
		log.Println(err)
		return
	}

	for index, game := range games {
		log.Printf("Match %s\n%+v\n", index, game)
	}
}
