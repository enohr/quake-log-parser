package model

type Match struct {
	TotalKills int
	Players    map[int]string
	Kills      map[int]int
}

type MatchJSON struct {
	TotalKills int
	Players    []string
	Kills      map[string]int
}

func NewMatch() *Match {
	match := &Match{
		Players: make(map[int]string),
		Kills:   make(map[int]int),
	}
	return match
}

func (m *Match) ToMatchJSON() MatchJSON {
	matchJson := MatchJSON{
		Players: make([]string, 0),
		Kills:   make(map[string]int),
	}

	for id, name := range m.Players {
		matchJson.Players = append(matchJson.Players, name)
		matchJson.Kills[name] = m.Kills[id]
	}
	matchJson.TotalKills = m.TotalKills

	return matchJson
}
