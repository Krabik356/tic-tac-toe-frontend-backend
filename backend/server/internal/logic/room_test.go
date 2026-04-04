package logic

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsCorrectMove(t *testing.T) {
	checks := []struct {
		name     string
		field    [3][3]string
		who      string
		turn     int
		row      int
		col      int
		expected error
	}{
		{
			name: "Correct move",
			field: [3][3]string{{"*", "*", "*"},
				{"0", "*", "X"},
				{"*", "*", "*"}},
			who:      "Bohdan",
			turn:     0,
			row:      1,
			col:      1,
			expected: nil,
		},
		{
			name: "Not your turn",
			field: [3][3]string{{"*", "*", "*"},
				{"*", "*", "*"},
				{"*", "*", "*"}},
			who:      "Stas",
			turn:     0,
			row:      1,
			col:      1,
			expected: errors.New("Not your turn"),
		},
		{
			name: "Incorrect move",
			field: [3][3]string{{"*", "*", "*"},
				{"*", "*", "*"},
				{"*", "*", "X"}},
			who:      "Bohdan",
			turn:     0,
			row:      2,
			col:      2,
			expected: errors.New("Incorrect move"),
		},
	}

	room := &Room{}

	client1 := &Client{Name: "Bohdan"}
	client2 := &Client{Name: "Stas"}
	room.Clients = [2]*Client{client1, client2}
	for id, check := range checks {
		room.Turn = check.turn
		room.Field = check.field
		t.Run(fmt.Sprintf("Test №%d %s", id, check.name), func(t *testing.T) {
			actual := room.isCorrectMove(check.row, check.col, check.who)
			if check.expected == nil {
				assert.NoError(t, actual)
			} else {
				assert.EqualError(t, actual, check.expected.Error())
			}
		})
	}
}

func TestIsPlayable(t *testing.T) {
	checks := []struct {
		name     string
		field    [3][3]string
		expected bool
	}{
		{
			name: "Playable1",
			field: [3][3]string{{"*", "X", "0"},
				{"*", "*", "X"},
				{"X", "X", "X"}},
			expected: true,
		},
		{
			name: "Playable2",
			field: [3][3]string{{"*", "0", "*"},
				{"0", "X", "*"},
				{"*", "X", "X"}},
			expected: true,
		},
		{
			name: "Playable3",
			field: [3][3]string{{"0", "X", "X"},
				{"0", "*", "X"},
				{"0", "0", "X"}},
			expected: true,
		},
		{
			name: "Not playable1",
			field: [3][3]string{{"0", "X", "X"},
				{"0", "0", "X"},
				{"0", "0", "X"}},
			expected: false,
		},
		{
			name: "Not playable2",
			field: [3][3]string{{"X", "X", "X"},
				{"0", "0", "X"},
				{"X", "0", "0"}},
			expected: false,
		},
	}

	room := &Room{}
	for id, check := range checks {
		room.Field = check.field

		t.Run(fmt.Sprintf("Test №%d %s", id, check.name), func(t *testing.T) {
			actual := room.isPlayable()
			assert.Equal(t, actual, check.expected)
		})
	}
}

func TestWinnerId(t *testing.T) {
	checks := []struct {
		name     string
		marker   string
		expected int
	}{
		{
			name:     "Is correct id by marker1",
			marker:   "X",
			expected: 0,
		},
		{
			name:     "Is correct id by marker1",
			marker:   "0",
			expected: 1,
		},
	}

	for id, check := range checks {
		t.Run(fmt.Sprintf("Test №%d %s", id, check.name), func(t *testing.T) {
			actual := winnerId(check.marker)
			assert.Equal(t, actual, check.expected)
		})
	}
}

func TestIsWin(t *testing.T) {
	checks := []struct {
		name             string
		field            [3][3]string
		expectedIsWin    bool
		expectedWinnerId int
	}{
		{
			name: "Win1",
			field: [3][3]string{{"*", "*", "X"},
				{"0", "*", "X"},
				{"0", "*", "X"}},
			expectedIsWin:    true,
			expectedWinnerId: 0,
		},
		{
			name: "Win2",
			field: [3][3]string{{"*", "*", "0"},
				{"X", "X", "X"},
				{"0", "*", "0"}},
			expectedIsWin:    true,
			expectedWinnerId: 0,
		},
		{
			name: "Win3",
			field: [3][3]string{{"X", "*", "0"},
				{"X", "0", "X"},
				{"0", "*", "0"}},
			expectedIsWin:    true,
			expectedWinnerId: 1,
		},
		{
			name: "Not win1",
			field: [3][3]string{{"X", "*", "0"},
				{"X", "*", "X"},
				{"0", "*", "0"}},
			expectedIsWin:    false,
			expectedWinnerId: 0,
		},
		{
			name: "Not win2",
			field: [3][3]string{{"0", "*", "0"},
				{"*", "X", "X"},
				{"0", "*", "*"}},
			expectedIsWin:    false,
			expectedWinnerId: 0,
		},
	}

	room := Room{}
	for id, check := range checks {
		room.Field = check.field
		t.Run(fmt.Sprintf("Test №%d %s", id, check.name), func(t *testing.T) {
			actualId, actualIs := room.isWin()
			assert.Equal(t, actualIs, check.expectedIsWin)
			assert.Equal(t, actualId, check.expectedWinnerId)
		})
	}
}
