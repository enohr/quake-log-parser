package model

import (
	"reflect"
	"testing"
)

func TestUserConnect(t *testing.T) {
	testCases := []struct {
		name     string
		id       int
		expected *Match
	}{
		{
			name: "User 1 connected",
			id:   1,
			expected: &Match{
				TotalKills:   0,
				Players:      map[int]*Player{1: {Name: "", Kills: 0}},
				MeansOfDeath: make(map[MeanOfDeath]int),
			},
		},
	}

	for _, tc := range testCases {
		match := NewMatch()
		match.AddPlayer(tc.id)

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}

}

func TestUpdateUserInfo(t *testing.T) {
	testCases := []struct {
		name       string
		id         int
		playerName string
		expected   *Match
	}{
		{
			name:       "User 1 changed nickname",
			id:         1,
			playerName: "Jogador",
			expected: &Match{
				TotalKills:   0,
				Players:      map[int]*Player{1: {Name: "Jogador", Kills: 0}},
				MeansOfDeath: make(map[MeanOfDeath]int),
			},
		},
	}

	for _, tc := range testCases {
		match := NewMatch()
		match.AddPlayer(tc.id)
		match.UpdateUserInfo(tc.id, tc.playerName)

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}
}

func TestUserChangeNicknameKeepKills(t *testing.T) {
	testCases := []struct {
		name          string
		id            int
		oldPlayerName string
		newPlayerName string
		expected      *Match
	}{
		{
			name:          "User 1 changed nickname",
			id:            1,
			oldPlayerName: "Jogador",
			newPlayerName: "Novo nome",
			expected: &Match{
				TotalKills: 1,
				Players: map[int]*Player{
					1: {Name: "Novo nome", Kills: 1},
					5: {Name: "", Kills: 0},
				},
				MeansOfDeath: map[MeanOfDeath]int{
					11: 1,
				},
			},
		},
	}

	for _, tc := range testCases {
		match := NewMatch()
		match.AddPlayer(tc.id)
		match.UpdateUserInfo(tc.id, tc.oldPlayerName)
		match.AddPlayer(5)
		match.ProcessKill(tc.id, 5, 11)

		match.UpdateUserInfo(tc.id, tc.newPlayerName)

		if !reflect.DeepEqual(match, tc.expected) {

			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}
}

func TestUserDisconnect(t *testing.T) {
	testCases := []struct {
		name     string
		id       int
		expected bool
	}{
		{
			name:     "User 1 disconnected",
			id:       1,
			expected: false,
		},
	}

	for _, tc := range testCases {
		match := NewMatch()
		match.AddPlayer(tc.id)
		match.DisconnectPlayer(tc.id)

		_, ok := match.Players[tc.id]

		if ok != tc.expected {
			t.Errorf("got %t expected %t", ok, tc.expected)
		}
	}
}

func TestUserKill(t *testing.T) {
	testCases := []struct {
		name     string
		killerID int
		victimID int
		meanID   int
		expected *Match
	}{
		{
			name:     "1 killed 2 with 5",
			killerID: 1,
			victimID: 2,
			meanID:   5,
			expected: &Match{
				TotalKills: 1,
				Players: map[int]*Player{
					1: {Name: "", Kills: 1},
					2: {Name: "", Kills: 0},
				},
				MeansOfDeath: map[MeanOfDeath]int{
					5: 1,
				},
			},
		},
	}

	for _, tc := range testCases {
		match := NewMatch()
		match.AddPlayer(tc.killerID)
		match.AddPlayer(tc.victimID)

		match.ProcessKill(tc.killerID, tc.victimID, tc.meanID)

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}
}

func TestWorldKill(t *testing.T) {
	testCases := []struct {
		name     string
		killerID int
		victimID int
		meanID   int
		expected *Match
	}{
		{
			name:     "1022 killed 2 with 5",
			killerID: 1022,
			victimID: 2,
			meanID:   22,
			expected: &Match{
				TotalKills: 1,
				Players: map[int]*Player{
					2: {Name: "", Kills: -1},
				},
				MeansOfDeath: map[MeanOfDeath]int{
					22: 1,
				},
			},
		},
	}

	for _, tc := range testCases {
		match := NewMatch()
		match.AddPlayer(tc.victimID)

		match.ProcessKill(tc.killerID, tc.victimID, tc.meanID)

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}
}

func TestMultipleUserKills(t *testing.T) {
	type KillEvent struct {
		killerID int
		victimID int
		meanID   int
	}

	killEvents := []KillEvent{
		{
			killerID: 1,
			victimID: 2,
			meanID:   5,
		},

		{
			killerID: 1,
			victimID: 3,
			meanID:   5,
		},

		{
			killerID: 2,
			victimID: 1,
			meanID:   8,
		},

		{
			killerID: 3,
			victimID: 1,
			meanID:   10,
		},
		{
			killerID: 1,
			victimID: 3,
			meanID:   8,
		},

		{
			killerID: 2,
			victimID: 1,
			meanID:   5,
		},

		{
			killerID: 2,
			victimID: 3,
			meanID:   5,
		},
	}

	testCases := []struct {
		name      string
		killEvent []KillEvent
		expected  *Match
	}{
		{
			name:      "Multiple kills",
			killEvent: killEvents,
			expected: &Match{
				TotalKills: 7,
				Players: map[int]*Player{
					1: {Name: "", Kills: 3},
					2: {Name: "", Kills: 3},
					3: {Name: "", Kills: 1},
				},
				MeansOfDeath: map[MeanOfDeath]int{
					5:  4,
					8:  2,
					10: 1,
				},
			},
		},
	}
	match := NewMatch()

	for _, tc := range testCases {
		for _, kill := range tc.killEvent {
			match.AddPlayer(kill.killerID)
			match.AddPlayer(kill.victimID)
			match.ProcessKill(kill.killerID, kill.victimID, kill.meanID)
		}

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}
}

func TestMultipleWorldKills(t *testing.T) {
	type KillEvent struct {
		killerID int
		victimID int
		meanID   int
	}

	killEvents := []KillEvent{
		{
			killerID: 1022,
			victimID: 2,
			meanID:   22,
		},

		{
			killerID: 1022,
			victimID: 3,
			meanID:   22,
		},

		{
			killerID: 1022,
			victimID: 1,
			meanID:   22,
		},

		{
			killerID: 1022,
			victimID: 1,
			meanID:   22,
		},
		{
			killerID: 1022,
			victimID: 3,
			meanID:   22,
		},
	}

	testCases := []struct {
		name      string
		killEvent []KillEvent
		expected  *Match
	}{
		{
			name:      "Multiple world kills",
			killEvent: killEvents,
			expected: &Match{
				TotalKills: 5,
				Players: map[int]*Player{
					1: {Name: "", Kills: -2},
					2: {Name: "", Kills: -1},
					3: {Name: "", Kills: -2},
				},
				MeansOfDeath: map[MeanOfDeath]int{
					22: 5,
				},
			},
		},
	}
	match := NewMatch()

	for _, tc := range testCases {
		for _, kill := range tc.killEvent {
			match.AddPlayer(kill.victimID)
			match.ProcessKill(kill.killerID, kill.victimID, kill.meanID)
		}

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}

}

func TestUserDisconnectKeepingHistory(t *testing.T) {
	testCases := []struct {
		name       string
		id         int
		playerName string
		expected   *Match
	}{
		{
			name:       "User 1 disconnected",
			id:         1,
			playerName: "User 1",
			expected: &Match{
				TotalKills: 0,
				Players: map[int]*Player{
					-1883664181901461531: {Name: "User 1", Kills: 0},
				},
				MeansOfDeath: make(map[MeanOfDeath]int),
			},
		},
	}

	for _, tc := range testCases {
		match := NewMatch()
		match.AddPlayer(tc.id)
		match.UpdateUserInfo(tc.id, tc.playerName)
		match.DisconnectPlayer(tc.id)

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}

}

func TestUserReconnects(t *testing.T) {
	testCases := []struct {
		name       string
		id         int
		playerName string
		expected   *Match
	}{
		{
			name:       "User 1 reconnected",
			id:         1,
			playerName: "User 1",
			expected: &Match{
				TotalKills: 1,
				Players: map[int]*Player{
					50: {Name: "User 1", Kills: 1},
					5:  {Name: "", Kills: 0},
				},
				MeansOfDeath: map[MeanOfDeath]int{
					10: 1,
				},
			},
		},
	}

	for _, tc := range testCases {
		match := NewMatch()
		match.AddPlayer(tc.id)
		match.AddPlayer(5)
		match.UpdateUserInfo(tc.id, tc.playerName)
		match.ProcessKill(tc.id, 5, 10)
		match.DisconnectPlayer(tc.id)

		match.AddPlayer(50)
		match.UpdateUserInfo(50, tc.playerName)

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}
}

func TestParsingMatchToMatchJSON(t *testing.T) {

	testCases := []struct {
		name     string
		input    *Match
		expected MatchJSON
	}{
		{
			name: "User 1 reconnected",
			input: &Match{
				TotalKills: 5,
				Players: map[int]*Player{
					1: {Name: "Jogador 1", Kills: 3},
					2: {Name: "Jogador 2", Kills: 2},
				},
				MeansOfDeath: map[MeanOfDeath]int{
					7:  2,
					10: 3,
				},
			},

			expected: MatchJSON{
				TotalKills: 5,
				Players: []string{
					"Jogador 1",
					"Jogador 2",
				},
				Kills: map[string]int{
					"Jogador 1": 3,
					"Jogador 2": 2,
				},
				MeansOfDeath: map[string]int{
					"MOD_ROCKET_SPLASH": 2,
					"MOD_RAILGUN":       3,
				},
			},
		},
	}

	for _, tc := range testCases {
		matchJSON := tc.input.ToMatchJSON()
		if !reflect.DeepEqual(matchJSON, tc.expected) {
			t.Errorf("got %+v expected %+v", matchJSON, tc.expected)
		}
	}

}
