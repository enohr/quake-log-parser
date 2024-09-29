package util

import (
	"testing"
)

func TestGenerateNameFNVHash(t *testing.T) {
	testCases := []struct {
		name       string
		playerName string
		expected   int
	}{
		{
			name:       "Hash Jogador 1",
			playerName: "Jogador 1",
			expected:   5590257778569622352,
		},
		{
			name:       "Hash Dono da Bola",
			playerName: "Dono da Bola",
			expected:   -7799584235392420246,
		},
		{
			name:       "Hash Isgalamido",
			playerName: "Isgalamido",
			expected:   -4666113643884491799,
		},

		{
			name:       "Hash Eduardo",
			playerName: "Eduardo",
			expected:   2157764416855725293,
		},

		{
			name:       "Hash Maluquinho",
			playerName: "Maluquinho",
			expected:   -5550577902254583248,
		},
	}

	for _, tc := range testCases {
		hash := GenerateNameFNVHash(tc.playerName)

		if hash != tc.expected {
			t.Errorf("got %d expected %d", hash, tc.expected)
		}
	}

}
