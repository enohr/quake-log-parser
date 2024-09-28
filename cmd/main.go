package main

import (
	"fmt"
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
		fmt.Printf("Match %s\n%s", index, game)
		fmt.Printf("------------------------------\n")
	}
}
