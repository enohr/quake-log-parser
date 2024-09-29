package model

import (
	"hash/fnv"
	"log/slog"
	"slices"
)

const (
	WORLD_ID = 1022
)

type Match struct {
	TotalKills   int
	Players      map[int]*Player
	MeansOfDeath map[MeanOfDeath]int
}

type MatchJSON struct {
	TotalKills   int            `json:"total_kills"`
	Players      []string       `json:"players"`
	Kills        map[string]int `json:"kills"`
	MeansOfDeath map[string]int `json:"kills_by_means"`
}

func NewMatch() *Match {
	match := &Match{
		Players:      make(map[int]*Player),
		MeansOfDeath: make(map[MeanOfDeath]int),
	}
	return match
}

func (m *Match) AddPlayer(playerID int) error {
	if _, ok := m.Players[playerID]; !ok {
		m.Players[playerID] = &Player{}
	}
	return nil
}

func (m *Match) UpdateUserInfo(playerID int, playerName string) error {
	h := fnv.New64a()
	h.Write([]byte(playerName))
	newID := int(h.Sum64())

	if player, ok := m.Players[newID]; ok {
		slog.Info("User has logged again. Restoring his kills", "id", playerID, "name", playerName, "kills", player.Kills)
		m.Players[playerID] = player
		delete(m.Players, newID)
		return nil
	}

	if _, ok := m.Players[playerID]; ok {
		m.Players[playerID].Name = playerName
	}
	return nil
}

func (m *Match) DisconnectPlayer(playerID int) error {
	if player, ok := m.Players[playerID]; ok {
		name := player.Name

		// TODO: Move it to utils
		h := fnv.New64a()
		h.Write([]byte(name))
		newID := h.Sum64()

		m.Players[int(newID)] = player
		delete(m.Players, playerID)
	}

	return nil
}

func (m *Match) ProcessKill(killerID, victimID, meanID int) error {
	if killerID == WORLD_ID {
		m.Players[victimID].Kills--
	} else if killerID != victimID {
		m.Players[killerID].Kills++
	}

	m.MeansOfDeath[MeanOfDeath(meanID)]++
	m.TotalKills++
	return nil
}

func (m *Match) ToMatchJSON() MatchJSON {
	matchJson := MatchJSON{
		Players:      make([]string, 0),
		Kills:        make(map[string]int),
		MeansOfDeath: make(map[string]int),
	}

	for _, player := range m.Players {
		matchJson.Players = append(matchJson.Players, player.Name)
		matchJson.Kills[player.Name] = player.Kills
	}
	for mean, kills := range m.MeansOfDeath {
		matchJson.MeansOfDeath[mean.String()] = kills
	}

	slices.Sort(matchJson.Players)
	matchJson.TotalKills = m.TotalKills

	return matchJson
}
