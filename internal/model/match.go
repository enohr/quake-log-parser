package model

import (
	"fmt"
	"strings"
)

type Match struct {
	TotalKills int
	Players    map[int]string
	Kills      map[int]int
}

func NewMatch() *Match {
	match := &Match{
		Players: make(map[int]string),
		Kills:   make(map[int]int),
	}
	return match
}

func (m *Match) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Total Kills: %d\n", m.TotalKills)
	fmt.Fprintf(&b, "Players (%d):\n", len(m.Players))

	for _, player := range m.Players {
		fmt.Fprintf(&b, "  %s\n", player)
	}

	fmt.Fprintf(&b, "Kills:\n")

	for player, kills := range m.Kills {
		fmt.Fprintf(&b, "  %s - %d kill(s)\n", player, kills)
	}

	return b.String()

}
