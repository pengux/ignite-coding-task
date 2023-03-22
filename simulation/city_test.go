package simulation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirection_opposite(t *testing.T) {
	tests := []struct {
		input    direction
		expected direction
	}{
		{north, south},
		{south, north},
		{east, west},
		{west, east},
		{direction("invalid"), direction("")},
	}

	for _, test := range tests {
		actual := test.input.opposite()
		assert.Equal(t, test.expected, actual)
	}
}

func TestDirection_isValid(t *testing.T) {
	tests := []struct {
		input    direction
		expected bool
	}{
		{north, true},
		{south, true},
		{east, true},
		{west, true},
		{"", false},
		{"foo", false},
	}

	for _, test := range tests {
		actual := test.input.isValid()
		assert.Equal(t, test.expected, actual)
	}
}

func TestCity_destroy(t *testing.T) {
	city1 := &city{name: "City1"}
	city2 := &city{name: "City2"}
	city3 := &city{name: "City3"}
	city1.neighbors = map[direction]*city{
		north: city2,
		east:  city3,
	}
	city2.neighbors = map[direction]*city{
		south: city1,
		east:  city3,
	}
	city3.neighbors = map[direction]*city{
		west:  city1,
		south: city2,
	}

	err := city1.destroy()
	assert.NoError(t, err, "expected no error when destroying city1")
	assert.True(t, city1.isDestroyed(), "expected city1 to be destroyed")

	for d, n := range city1.neighbors {
		_, ok := n.neighbors[d.opposite()]
		assert.Falsef(t, ok, "expected neighbor %s of city1 to have no road in direction %s", n.name, d)
	}
}

func TestCity_destroy_invalidNeighborRoad(t *testing.T) {
	city1 := &city{name: "City1"}
	city2 := &city{name: "City2"}
	city1.neighbors = map[direction]*city{
		north: city2,
	}
	city2.neighbors = map[direction]*city{
		north: city1,
	}

	err := city1.destroy()
	assert.EqualError(t, err, "neighbor city City2 has no road in direction north to City1")
}

func TestCity_battle(t *testing.T) {
	city := &city{name: "City1"}
	alien1 := &alien{name: 1, currentCity: city}
	alien2 := &alien{name: 2, currentCity: city}
	alien3 := &alien{name: 3, currentCity: city}
	city.visitingAliens = []*alien{alien1, alien2, alien3}

	err := city.battle()
	assert.NoError(t, err, "expected no error when battling city")

	assert.True(t, city.isDestroyed(), "expected city to be destroyed after battle")
	assert.True(t, alien1.isDead(), "expected alien1 to be dead after battle")
	assert.True(t, alien2.isDead(), "expected alien2 to be dead after battle")
	assert.True(t, alien3.isDead(), "expected alien3 to be dead after battle")
}

func TestCity_reachableFrom(t *testing.T) {
	a := &city{name: "A", neighbors: make(map[direction]*city)}
	b := &city{name: "B", neighbors: make(map[direction]*city)}
	c := &city{name: "C", neighbors: make(map[direction]*city)}
	d := &city{name: "D", neighbors: make(map[direction]*city)}

	a.neighbors[east] = b
	b.neighbors[west] = a
	b.neighbors[north] = c
	c.neighbors[south] = b

	testCases := []struct {
		name     string
		src      *city
		dst      *city
		expected bool
	}{
		{"same city", a, a, true},
		{"directly connected", a, b, true},
		{"indirectly connected", a, c, true},
		{"not connected", a, d, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			visited := make(map[string]bool)
			actual := tc.src.reachableFrom(tc.dst, visited)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestCity_string(t *testing.T) {
	testCases := []struct {
		name     string
		city     *city
		expected string
	}{
		{
			name: "without directions",
			city: &city{
				name: "City1",
			},
			expected: "City1",
		},
		{
			name: "with all directions, the output should be sorted",
			city: &city{
				name: "City1",
				neighbors: map[direction]*city{
					north: {name: "City2"},
					south: {name: "City3"},
					east:  {name: "City4"},
					west:  {name: "City5"},
				},
			},
			expected: "City1 east=City4 north=City2 south=City3 west=City5",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.city.string()
			assert.Equal(t, tc.expected, actual)
		})
	}
}
