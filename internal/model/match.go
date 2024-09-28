package model

const (
	WORLD_ID = 1022
)

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

func (m *Match) AddPlayer(playerID int) error {
	if _, ok := m.Players[playerID]; !ok {
		m.Players[playerID] = ""
	}
	return nil
}

func (m *Match) UpdateUserInfo(playerID int, playerName string) error {
	if _, ok := m.Players[playerID]; ok {
		m.Players[playerID] = playerName
	}
	return nil
}

func (m *Match) ProcessKill(killerID, victimID, meanID int) error {
	if killerID == WORLD_ID {
		m.Kills[victimID]--
	} else if killerID != victimID {
		m.Kills[killerID]++
	}
	m.TotalKills++
	return nil
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
