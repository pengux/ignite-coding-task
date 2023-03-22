package simulation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSimulation(t *testing.T) {
	input := `Foo north=Bar east=Baz
Bar south=Foo
Baz west=Foo north=Qux
Qux south=Baz`
	nrOfAliens := 5

	sim, err := NewSimulation(strings.NewReader(input), nrOfAliens)
	assert.NoError(t, err, "expected no error when creating simulation")
	assert.Len(t, sim.cities, 4, "expected 4 cities in simulation world")
	assert.Lenf(t, sim.aliens, nrOfAliens, "expected %d aliens in simulation world", nrOfAliens)

	// assert that aliens are found in the cities
	for _, alien := range sim.aliens {
		assert.NotNil(t, alien.currentCity, "expected alien to have a non-nil `currentCity`")
		assert.Contains(t, alien.currentCity.visitingAliens, alien, "expected alien to be in a city's `visitingAliens` slice")
	}
}

func TestSimulation_Run(t *testing.T) {
	tests := []struct {
		name       string
		simulation *Simulation
		expected   simState
	}{
		{
			name: "all aliens are trapped",
			simulation: func() *Simulation {
				city1 := &city{name: "City1"}
				city2 := &city{name: "City2"}
				city3 := &city{name: "City3"}

				return &Simulation{
					aliens: []*alien{
						{name: 1, currentCity: city1},
						{name: 2, currentCity: city2},
						{name: 3, currentCity: city3},
					},
					cities: map[string]*city{
						"City1": city1,
						"City2": city2,
						"City3": city3,
					},
				}
			}(),
			expected: SimStateAllAliensDeadOrTrapped,
		},
		{
			name: "all aliens dead",
			simulation: &Simulation{
				aliens: []*alien{
					{name: 1},
					{name: 2},
					{name: 3},
				},
			},
			expected: SimStateAllAliensDeadOrTrapped,
		},
		{
			name: "only 1 alien alive and not trapped",
			simulation: func() *Simulation {
				alien1 := &alien{name: 1}
				alien2 := &alien{name: 2}
				alien3 := &alien{name: 3}

				city1 := &city{name: "City1"}
				city2 := &city{name: "City2"}
				city3 := &city{name: "City3"}

				alien1.currentCity = city1
				city1.visitingAliens = []*alien{alien1}
				city1.neighbors = map[direction]*city{
					north: city2,
					east:  city3,
				}
				city2.neighbors = map[direction]*city{
					south: city1,
				}
				city3.neighbors = map[direction]*city{
					west: city1,
				}

				return &Simulation{
					aliens: []*alien{
						alien1,
						alien2,
						alien3,
					},
					cities: map[string]*city{
						"City1": city1,
						"City2": city2,
						"City3": city3,
					},
				}
			}(),
			expected: SimStateOnlyOneAlienLeft,
		},
		{
			name: "alive aliens cannot reach each other",
			simulation: func() *Simulation {
				alien1 := &alien{name: 1}
				alien2 := &alien{name: 2}

				city1 := &city{name: "City1"}
				city2 := &city{name: "City2"}
				city3 := &city{name: "City3"}
				city4 := &city{name: "City4"}

				alien1.currentCity = city1
				city1.visitingAliens = []*alien{alien1}
				city1.neighbors = map[direction]*city{
					north: city2,
				}
				city2.neighbors = map[direction]*city{
					south: city1,
				}

				alien2.currentCity = city3
				city3.visitingAliens = []*alien{alien2}
				city3.neighbors = map[direction]*city{
					west: city4,
				}
				city4.neighbors = map[direction]*city{
					east: city3,
				}

				return &Simulation{
					aliens: []*alien{
						alien1,
						alien2,
					},
					cities: map[string]*city{
						"City1": city1,
						"City2": city2,
						"City3": city3,
						"City4": city4,
					},
				}
			}(),
			expected: SimStateAliveAliensDisconnected,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.simulation.Run()
			assert.NoError(t, err, "expected no error when running simulation")
			assert.Equal(t, test.expected, actual, "expected simulation to be in state %s", test.expected)
		})
	}
}

func TestSimulation_SurvivedCities(t *testing.T) {
	sim := &Simulation{
		cities: map[string]*city{
			"City1": {name: "City1"},
			"City2": {name: "City2"},
			"City3": {name: "City3"},
		},
	}

	sim.cities["City1"].destroyed = true
	sim.cities["City2"].destroyed = true

	actual := sim.SurvivedCities()
	assert.Len(t, actual, 1, "expected 1 city to be alive")

	city3, ok := sim.cities["City3"]
	assert.True(t, ok, "expected city3 to be in simulation")
	assert.Equal(t, "City3", city3.name, "expected survived city to be City3")
}

func TestSimulation_CitiesToString(t *testing.T) {
	city1 := &city{name: "City1"}
	city2 := &city{name: "City2"}
	city3 := &city{name: "City3"}
	city4 := &city{name: "City4"}

	city1.neighbors = map[direction]*city{
		north: city2,
	}
	city2.neighbors = map[direction]*city{
		south: city1,
	}
	city3.neighbors = map[direction]*city{
		west: city4,
	}
	city4.neighbors = map[direction]*city{
		east: city3,
	}

	sim := &Simulation{
		cities: map[string]*city{
			"City1": city1,
			"City2": city2,
			"City3": city3,
			"City4": city4,
		},
	}

	sim.cities["City1"].destroyed = true
	sim.cities["City2"].destroyed = true

	actual := sim.CitiesToString(sim.SurvivedCities())
	assert.Equal(t, `City3 west=City4
City4 east=City3
`, actual)
}
