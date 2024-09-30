package parser

import (
	"reflect"
	"testing"

	"github.com/enohr/quake-log-parser/internal/model"
)

// Tests the parse of a line with user connection
func TestParseUserConnectLine(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		expected *model.Match
	}{
		{
			name: "User with ID 2 connected",
			line: " 1:47 ClientConnect: 2",
			expected: &model.Match{
				TotalKills:   0,
				Players:      map[int]*model.Player{2: {Name: "", Kills: 0}},
				MeansOfDeath: make(map[model.MeanOfDeath]int),
			},
		},
	}

	for _, tc := range testCases {
		match := model.NewMatch()
		parseLine(tc.line, match)

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}
}

// Tests the parse of a line with user updating info
func TestParseUserUpdateInfoLine(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		expected *model.Match
	}{
		{
			name: "User with ID 2 changed nickname",
			line: ` 1:47 ClientUserinfoChanged: 2 n\Dono da Bola\t\0\model\sarge\hmodel\sarge\g_redteam\\g_blueteam\\c1\4\c2\5\hc\95\w\0\l\0\tt\0\tl\0"`,
			expected: &model.Match{
				TotalKills:   0,
				Players:      map[int]*model.Player{2: {Name: "Dono da Bola", Kills: 0}},
				MeansOfDeath: make(map[model.MeanOfDeath]int),
			},
		},
	}

	for _, tc := range testCases {
		match := model.NewMatch()
		match.AddPlayer(2)
		parseLine(tc.line, match)

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}
}

// Tests the parse of a line with user disconnection
func TestParseUserDisconnectLine(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		expected *model.Match
	}{
		{
			name: "User with ID 2 disconnected",
			line: " 13:26 ClientDisconnect: 2",
			expected: &model.Match{
				TotalKills:   0,
				Players:      map[int]*model.Player{-3750763034362895579: {Name: "", Kills: 0}},
				MeansOfDeath: make(map[model.MeanOfDeath]int),
			},
		},
	}

	for _, tc := range testCases {
		match := model.NewMatch()
		match.AddPlayer(2)
		parseLine(tc.line, match)

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}
}

// Tests the parse of a line with user kill
func TestParseUserKillLine(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		expected *model.Match
	}{
		{
			name: "User with ID 2 killed user with ID 5",
			line: " 14:15 Kill: 2 5 10: Zeh killed Assasinu Credi by MOD_RAILGUN",
			expected: &model.Match{
				TotalKills: 1,
				Players: map[int]*model.Player{
					2: {Name: "", Kills: 1},
					5: {Name: "", Kills: 0},
				},
				MeansOfDeath: map[model.MeanOfDeath]int{
					10: 1,
				},
			},
		},
	}

	for _, tc := range testCases {
		match := model.NewMatch()
		match.AddPlayer(2)
		match.AddPlayer(5)
		parseLine(tc.line, match)

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}
}

// Tests the parse of a line with world kill
func TestParseWorldKillLine(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		expected *model.Match
	}{
		{
			name: "World killed user with ID 5",
			line: " 14:02 Kill: 1022 5 22: <world> killed Assasinu Credi by MOD_TRIGGER_HURT",
			expected: &model.Match{
				TotalKills: 1,
				Players: map[int]*model.Player{
					5: {Name: "", Kills: -1},
				},
				MeansOfDeath: map[model.MeanOfDeath]int{
					22: 1,
				},
			},
		},
	}

	for _, tc := range testCases {
		match := model.NewMatch()
		match.AddPlayer(5)
		parseLine(tc.line, match)

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}
}

// Tests the parse of multiple lines with kills
func TestParseMultipleUserKillLines(t *testing.T) {
	testCases := []struct {
		name     string
		lines    []string
		expected *model.Match
	}{
		{
			name: "Multiple User Kills",
			lines: []string{
				" 13:27 Kill: 6 2 7: Zeh killed Isgalamido by MOD_ROCKET_SPLASH",
				" 13:27 Kill: 6 5 7: Zeh killed Assasinu Credi by MOD_ROCKET_SPLASH",
				" 13:37 Kill: 3 4 7: Oootsimo killed Dono da Bola by MOD_ROCKET_SPLASH",
				" 13:43 Kill: 2 5 6: Isgalamido killed Assasinu Credi by MOD_ROCKET",
				" 13:45 Kill: 1022 7 22: <world> killed Mal by MOD_TRIGGER_HURT",
				" 13:46 Kill: 4 3 7: Dono da Bola killed Oootsimo by MOD_ROCKET_SPLASH",
				" 13:46 Kill: 6 2 6: Zeh killed Isgalamido by MOD_ROCKET",
			},
			expected: &model.Match{
				TotalKills: 7,
				Players: map[int]*model.Player{
					2: {Name: "", Kills: 1},
					3: {Name: "", Kills: 1},
					4: {Name: "", Kills: 1},
					6: {Name: "", Kills: 3},
					7: {Name: "", Kills: -1},
				},
				MeansOfDeath: map[model.MeanOfDeath]int{
					7:  4,
					6:  2,
					22: 1,
				},
			},
		},
	}

	for _, tc := range testCases {
		match := model.NewMatch()
		match.AddPlayer(2)
		match.AddPlayer(3)
		match.AddPlayer(4)
		match.AddPlayer(6)
		match.AddPlayer(7)

		for _, line := range tc.lines {
			parseLine(line, match)
		}

		if !reflect.DeepEqual(match, tc.expected) {
			t.Errorf("got %+v expected %+v", match, tc.expected)
		}
	}

}
