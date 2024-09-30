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

// Creates a new instance of Match
func NewMatch() *Match {
	match := &Match{
		Players:      make(map[int]*Player),
		MeansOfDeath: make(map[MeanOfDeath]int),
	}
	return match
}

// Add a player to the current Match
func (m *Match) AddPlayer(playerID int) error {
	if _, ok := m.Players[playerID]; !ok {
		m.Players[playerID] = &Player{}
	}
	return nil
}

// Update nickname of a player on the current match
// Check if the name hash of the user is the same as a user that
// has disconnected before. If so, restore the kill count
func (m *Match) UpdateUserInfo(playerID int, playerName string) error {
	oldHashID := util.GenerateNameFNVHash(playerName)

	// verify if user was connected before
	// If it was, restore the kill count
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

// Disconnect a player for the current match
// Generate a hash based on player name to allow
// kill count restore if user reconnects
func (m *Match) DisconnectPlayer(playerID int) error {
	if player, ok := m.Players[playerID]; ok {
		name := player.Name
		oldHashID := util.GenerateNameFNVHash(name)

		// Move the player to other map position
		// With this, it can restore the kill count
		m.Players[int(oldHashID)] = player
		delete(m.Players, playerID)
	}

	return nil
}

// Process a kill event
// If the killer is World, decrement the victim kill count
// Update the match kill count and the deaths by weapon count
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

// Transform the Match model to the desired output format
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

	// Sort players alphabetically
	slices.Sort(matchJson.Players)
	matchJson.TotalKills = m.TotalKills

	// Sort players based on kill count to create the match leaderboard
	sort.Slice(matchJson.Leaderboard, func(i, j int) bool {
		p1, p2 := matchJson.Leaderboard[i], matchJson.Leaderboard[j]

		k1 := matchJson.Kills[p1]
		k2 := matchJson.Kills[p2]

		return k1 > k2
	})

	return matchJson
}
