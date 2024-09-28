package model

type Match struct {
	TotalKills int
	Players    []string
	Kills      map[string]int
}

func NewMatch() *Match {
	match := &Match{
		Players: make([]string, 0),
		Kills:   make(map[string]int),
	}
	return match
}
