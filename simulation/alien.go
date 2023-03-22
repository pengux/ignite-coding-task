package simulation

import (
	"fmt"
	"log"
	"math/rand"
)

type alien struct {
	name        int
	currentCity *city
}

func (a *alien) string() string {
	return fmt.Sprintf("Alien %d", a.name)
}

func (a *alien) die() {
	a.currentCity = nil
}

func (a *alien) isDead() bool {
	return a.currentCity == nil
}

// isTrapped returns true if the alien is trapped in a city which is not
// destroyed but doesn't have any neighbors
func (a *alien) isTrapped() bool {
	return !a.isDead() && !a.currentCity.isDestroyed() && len(a.currentCity.neighbors) == 0
}

// move moves the alien randomly from a city to one of its neighbors if
// possible. The alien may also stay at the same city, the return boolean value
// indicates if the alien moved or not
func (a *alien) move() bool {
	if a.isDead() || a.isTrapped() {
		return false
	}

	// randomly pick a neighbor
	directions := make([]direction, len(a.currentCity.neighbors))
	i := 0
	for d := range a.currentCity.neighbors {
		directions[i] = d
		i++
	}

	moveDirectionIndex := rand.Intn(len(directions)+1) - 1
	// alien decided to stay
	if moveDirectionIndex < 0 {
		return false
	}

	a.goToCity(a.currentCity.neighbors[directions[moveDirectionIndex]])

	return true
}

func (a *alien) goToCity(c *city) {
	a.currentCity.removeAlien(a)
	a.currentCity = c
	a.currentCity.visitingAliens = append(a.currentCity.visitingAliens, a)

	log.Printf("%s moved to %s", a.string(), c.name)
}
