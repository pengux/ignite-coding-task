package simulation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlien_string(t *testing.T) {
	alien := &alien{name: 1}
	expected := "Alien 1"
	actual := alien.string()
	assert.Equalf(t, expected, actual, "expected alien to be printed as %s", expected)
}

func TestAlien_die(t *testing.T) {
	alien := &alien{name: 1, currentCity: &city{}}
	alien.die()
	assert.True(t, alien.isDead(), "expected alien to be dead after calling die()")
}

func TestAlien_isTrapped(t *testing.T) {
	alien1 := &alien{name: 1, currentCity: &city{name: "City1"}}
	alien2 := &alien{name: 2, currentCity: &city{name: "City2", neighbors: map[direction]*city{}, destroyed: true}}
	alien3 := &alien{name: 3, currentCity: &city{name: "City3", neighbors: map[direction]*city{}}}

	assert.True(t, alien1.isTrapped(), "expected alien1 to be trapped as City1 doesn't have any neighbors")
	assert.False(t, alien2.isTrapped(), "expected alien2 to not be trapped as the city is destroyed")
	assert.True(t, alien3.isTrapped(), "expected alien3 to be trapped")
}

func TestAlien_move(t *testing.T) {
	city1 := &city{name: "City1"}
	city2 := &city{name: "City2"}
	city1.neighbors = map[direction]*city{
		north: city2,
	}
	city2.neighbors = map[direction]*city{
		south: city1,
	}
	alien := &alien{name: 1, currentCity: city1}

	moved := alien.move()

	if moved {
		assert.Equal(t, city2, alien.currentCity, "expected alien to move from City1 to City2")
	} else {
		assert.Equal(t, city1, alien.currentCity, "expected alien to not move")
	}
}

func TestAlien_goToCity(t *testing.T) {
	alien := &alien{name: 1, currentCity: &city{}}
	city := &city{name: "City1"}

	alien.goToCity(city)

	assert.Equal(t, city, alien.currentCity, "expected alien to go to City1")
	assert.Equal(t, 1, len(city.visitingAliens), "expected City1 to have 1 visiting alien")
}
