package model

import (
	"slices"
	"sort"

	"github.com/enohr/quake-log-parser/util"
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
	Leaderboard  []string       `json:"leaderboard"`
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
	oldHashID := util.GenerateNameFNVHash(playerName)

	if player, ok := m.Players[oldHashID]; ok {
		m.Players[playerID] = player
		delete(m.Players, oldHashID)
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
		oldHashID := util.GenerateNameFNVHash(name)

		m.Players[int(oldHashID)] = player
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
		Leaderboard:  make([]string, 0),
		MeansOfDeath: make(map[string]int),
	}

	for _, player := range m.Players {
		matchJson.Players = append(matchJson.Players, player.Name)
		matchJson.Leaderboard = append(matchJson.Leaderboard, player.Name)
		matchJson.Kills[player.Name] = player.Kills
	}
	for mean, kills := range m.MeansOfDeath {
		matchJson.MeansOfDeath[mean.String()] = kills
	}

	slices.Sort(matchJson.Players)
	matchJson.TotalKills = m.TotalKills

	sort.Slice(matchJson.Leaderboard, func(i, j int) bool {
		p1, p2 := matchJson.Leaderboard[i], matchJson.Leaderboard[j]

		k1 := matchJson.Kills[p1]
		k2 := matchJson.Kills[p2]

		return k1 > k2
	})

	return matchJson
}
