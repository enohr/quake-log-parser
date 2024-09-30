package output

import (
	"reflect"
	"testing"

	"github.com/enohr/quake-log-parser/internal/model"
)

func TestParsingMatchesToMatchesJSON(t *testing.T) {

	testCases := []struct {
		name     string
		input    map[string]*model.Match
		expected map[string]model.MatchJSON
	}{
		{
			name: "User 1 reconnected",
			input: map[string]*model.Match{
				"game_1": &model.Match{
					TotalKills: 5,
					Players: map[int]*model.Player{
						1: {Name: "Jogador 1", Kills: 2},
						2: {Name: "Jogador 2", Kills: 3},
					},
					MeansOfDeath: map[model.MeanOfDeath]int{
						7:  2,
						10: 3,
					},
				},
			},

			expected: map[string]model.MatchJSON{
				"game_1": model.MatchJSON{
					TotalKills: 5,
					Players: []string{
						"Jogador 1",
						"Jogador 2",
					},
					Kills: map[string]int{
						"Jogador 1": 2,
						"Jogador 2": 3,
					},
					Leaderboard: []string{
						"Jogador 2",
						"Jogador 1",
					},
					MeansOfDeath: map[string]int{
						"MOD_ROCKET_SPLASH": 2,
						"MOD_RAILGUN":       3,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		matchesJSON, err := matchToJSON(tc.input)

		if err != nil {
			t.Errorf("function returned with an error %s", err)
		}

		if !reflect.DeepEqual(matchesJSON, tc.expected) {
			t.Errorf("got %+v expected %+v", matchesJSON, tc.expected)
		}

	}

}
