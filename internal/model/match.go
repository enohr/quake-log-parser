package model

import (
	"hash/fnv"
	"log/slog"
)

const (
	WORLD_ID = 1022
)

type Match struct {
	TotalKills int
	Players    map[int]*Player
}

type MatchJSON struct {
	TotalKills int
	Players    []string
	Kills      map[string]int
}

func NewMatch() *Match {
	match := &Match{
		Players: make(map[int]*Player),
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
	m.TotalKills++
	return nil
}

func (m *Match) ToMatchJSON() MatchJSON {
	matchJson := MatchJSON{
		Players: make([]string, 0),
		Kills:   make(map[string]int),
	}

	for _, player := range m.Players {
		matchJson.Players = append(matchJson.Players, player.Name)
		matchJson.Kills[player.Name] = player.Kills
	}
	matchJson.TotalKills = m.TotalKills

	return matchJson
}
