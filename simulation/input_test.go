package simulation

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expected      map[string]*city
		expectedError string
	}{
		{
			name: "valid input",
			input: `
				A north=B east=C
				B south=A west=D
				C west=A south=D north=E
				D east=B north=C west=E
				E south=C east=D
			`,
			expected: map[string]*city{
				"A": {
					name: "A",
					neighbors: map[direction]*city{
						north: {name: "B"},
						east:  {name: "C"},
					},
				},
				"B": {
					name: "B",
					neighbors: map[direction]*city{
						south: {name: "A"},
						west:  {name: "D"},
					},
				},
				"C": {
					name: "C",
					neighbors: map[direction]*city{
						west:  {name: "A"},
						south: {name: "D"},
						north: {name: "E"},
					},
				},
				"D": {
					name: "D",
					neighbors: map[direction]*city{
						east:  {name: "B"},
						north: {name: "C"},
						west:  {name: "E"},
					},
				},
				"E": {
					name: "E",
					neighbors: map[direction]*city{
						south: {name: "C"},
						east:  {name: "D"},
					},
				},
			},
		},
		{
			name: "empty input",
			input: `
				
			`,
			expected: map[string]*city{},
		},
		{
			name: "invalid input format",
			input: `
				A
				B invalid-field
			`,
			expectedError: "invalid direction 'invalid-field' for city 'B'",
		},
		{
			name: "duplicate city name",
			input: `
				A north=B
				B south=A
				A east=C
			`,
			expectedError: "duplicate city name: 'A' at line 4",
		},
		{
			name: "unknown neighbor city",
			input: `
				A north=B
				B south=C
			`,
			expectedError: "unknown neighbor city 'C' in direction 'south' for city 'B'",
		},
		{
			name: "invalid direction",
			input: `
				A north=B
				B invalid-dir=A
			`,
			expectedError: "invalid direction 'invalid-dir' for city 'B'",
		},
		{
			name: "invalid city neighbor",
			input: `
				A north=B
				B north=A
			`,
			expectedError: "neighbor city 'B' has no road in direction 'south' to city 'A'",
		},
	}

	toJSON := func(cities map[string]*city) string {
		json, _ := json.Marshal(cities)
		return string(json)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := parseInput(strings.NewReader(tc.input))
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, toJSON(tc.expected), toJSON(actual))
			}
		})
	}
}
