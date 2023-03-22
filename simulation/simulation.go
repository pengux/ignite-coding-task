package simulation

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type simState string

// simulation states
const (
	// simulation is still running
	SimStateRunning simState = "STATE_RUNNING"
	// all aliens are dead or trapped
	SimStateAllAliensDeadOrTrapped = "STATE_ALL_ALIENS_DEAD_OR_TRAPPED"
	// only one alien left (cannot do battles)
	SimStateOnlyOneAlienLeft = "STATE_ONLY_ONE_ALIEN_LEFT"
	// all alive aliens cannot reach each other
	SimStateAliveAliensDisconnected = "STATE_ALIVE_ALIENS_DISCONNECTED"
	// all cities are destroyed (cannot reach this state)
	// SimStateAllCitiesDestroyed
)

type Simulation struct {
	// cities contains all cities and their neighbors
	cities map[string]*city
	// aliens contains all aliens in the world
	aliens    []*alien
	iteration int
}

// NewSimulation creates a new simulation from the input string and number of
// aliens to randomly place in the world
func NewSimulation(input io.Reader, nrOfAliens int) (*Simulation, error) {
	cities, err := parseInput(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	// - create aliens
	aliens := make([]*alien, nrOfAliens)
	for i := 0; i < nrOfAliens; i++ {
		aliens[i] = &alien{name: i + 1}
	}

	cityKeys := make([]string, len(cities))
	i := 0
	for cityName := range cities {
		cityKeys[i] = cityName
		i++
	}

	// randomly place aliens in cities
	for _, alien := range aliens {
		// get random city
		rndCityIndex := rand.Intn(len(cityKeys))
		cities[cityKeys[rndCityIndex]].visitingAliens = append(
			cities[cityKeys[rndCityIndex]].visitingAliens,
			alien,
		)
		alien.currentCity = cities[cityKeys[rndCityIndex]]
	}

	return &Simulation{cities: cities, aliens: aliens}, nil
}

// Run runs the simulation until it ends
func (s *Simulation) Run() (simState, error) {
	for {
		s.iteration++

		log.Printf("iteration %d", s.iteration)

		// move aliens
		for _, alien := range s.aliens {
			alien.move()
		}

		// simulate battles
		for _, city := range s.cities {
			err := city.battle()
			if err != nil {
				return SimStateRunning, fmt.Errorf("failed to simulate battle: %w", err)
			}
		}

		simState := s.checkEndState()
		if simState != SimStateRunning {
			return simState, nil
		}
	}
}

func (s *Simulation) checkEndState() simState {
	var state simState

	// check if all aliens are dead or trapped
	state = SimStateAllAliensDeadOrTrapped
	for _, alien := range s.aliens {
		log.Printf("alien %d is dead: %t, is trapped: %t", alien.name, alien.isDead(), alien.isTrapped())
		if !alien.isDead() && !alien.isTrapped() {
			state = SimStateRunning
			break
		}
	}

	if state != SimStateRunning {
		return state
	}

	// check if only one alien left
	state = SimStateOnlyOneAlienLeft
	aliveAliens := make([]*alien, 0)
	for _, alien := range s.aliens {
		if !alien.isDead() {
			aliveAliens = append(aliveAliens, alien)
			log.Printf("alien %d is alive", alien.name)
		} else {
			log.Printf("alien %d is dead", alien.name)
		}
	}
	if len(aliveAliens) > 1 {
		state = SimStateRunning
	}

	if state != SimStateRunning {
		return state
	}

	// check if alive aliens are able to reach each other
	state = SimStateAliveAliensDisconnected
	for _, alien := range aliveAliens {
		for _, target := range aliveAliens {
			if alien == target {
				continue
			}
			visited := make(map[string]bool)
			log.Printf("alien %d in city %s is able to reach alien %d in city %s", alien.name, alien.currentCity.name, target.name, target.currentCity.name)
			if alien.currentCity.reachableFrom(target.currentCity, visited) {
				state = SimStateRunning
				return state
			}
		}
	}

	// check if all cities are destroyed
	// state = SimStateAllCitiesDestroyed
	// for _, city := range s.cities {
	// 	if !city.isDestroyed() {
	// 		state = SimStateRunning
	// 		break
	// 	}
	// }
	//
	// if state != SimStateRunning {
	// 	return state
	// }

	return state
}

// SurvivedCities returns all cities that are not destroyed
func (s *Simulation) SurvivedCities() map[string]*city {
	survivedCities := make(map[string]*city)
	for cityName, city := range s.cities {
		if !city.isDestroyed() {
			survivedCities[cityName] = city
		}
	}
	return survivedCities
}

// CitiesToString returns a string of the cities and their neighbors in the
// format:
/*
Foo north=Bar east=Baz
Bar south=Foo
Baz west=Foo north=Qux
Qux south=Baz
*/
// The city names will be sorted before being added to the string to make the
// output deterministic
func (s *Simulation) CitiesToString(cities map[string]*city) string {
	cityNames := maps.Keys(cities)
	slices.Sort(cityNames)

	builder := strings.Builder{}
	for _, cityName := range cityNames {
		builder.WriteString(cities[cityName].string())
		builder.WriteString("\n")
	}

	return builder.String()
}
