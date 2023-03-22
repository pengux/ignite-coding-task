package simulation

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type direction string

const (
	north direction = "north"
	east  direction = "east"
	south direction = "south"
	west  direction = "west"
)

func (d direction) opposite() direction {
	switch d {
	case north:
		return south
	case east:
		return west
	case south:
		return north
	case west:
		return east
	}
	return ""
}

func (d direction) isValid() bool {
	return d == north || d == east || d == south || d == west
}

type city struct {
	name           string
	neighbors      map[direction]*city
	visitingAliens []*alien
	// use flag to differentiate between a destroyed and isolated (all
	// neighbors are destroyed) city
	destroyed bool
}

func (c *city) destroy() error {
	for d, neighbor := range c.neighbors {
		city, ok := neighbor.neighbors[d.opposite()]
		if !ok || city != c {
			return fmt.Errorf("neighbor city %s has no road in direction %s to %s", neighbor.name, d, c.name)
		}

		delete(neighbor.neighbors, d.opposite())
		delete(c.neighbors, d)
	}

	c.destroyed = true

	return nil
}

func (c *city) isDestroyed() bool {
	return c.destroyed
}

// battle simulate a battle between at least 2 aliens in the city. All aliens
// are killed and the city is destroyed in the process
func (c *city) battle() error {
	if len(c.visitingAliens) < 2 {
		return nil
	}

	alienNames := make([]string, len(c.visitingAliens))
	for i, a := range c.visitingAliens {
		a.die()
		alienNames[i] = a.string()
	}
	c.visitingAliens = nil

	log.Printf("city %s has been destroyed by %s", c.name, strings.Join(alienNames, " and "))

	return c.destroy()
}

func (c *city) removeAlien(a *alien) {
	for i, alien := range c.visitingAliens {
		if alien == a {
			c.visitingAliens = append(c.visitingAliens[:i], c.visitingAliens[i+1:]...)
			break
		}
	}
}

// reachableFrom returns true if the city is reachable from the destination. It
// uses a recursive depth first search for the node traversal
func (c *city) reachableFrom(dst *city, visited map[string]bool) bool {
	if c == dst {
		return true
	}

	visited[c.name] = true

	for _, neighbor := range c.neighbors {
		if !visited[neighbor.name] && neighbor.reachableFrom(dst, visited) {
			return true
		}
	}

	return false
}

// string returns a string representation of the city and its neighbors
// the "roads" are sorted by direction to make the output deterministic
func (c *city) string() string {
	builder := strings.Builder{}
	builder.WriteString(c.name)

	if len(c.neighbors) > 0 {
		directions := maps.Keys(c.neighbors)
		slices.Sort(directions)

		roads := make([]string, len(directions))
		for i, dir := range directions {
			roads[i] = string(dir) + "=" + c.neighbors[direction(dir)].name
		}

		builder.WriteString(" ")
		builder.WriteString(strings.Join(roads, " "))
	}

	return builder.String()
}
