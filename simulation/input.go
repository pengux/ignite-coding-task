package simulation

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func parseInput(input io.Reader) (map[string]*city, error) {
	intermediateCities := make(map[string][]string)
	cities := make(map[string]*city)

	scanner := bufio.NewScanner(input)
	lineNr := 0
	for scanner.Scan() {
		lineNr++

		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		// each line should have at least one field, the city name. Cities without any neighbors are valid as input
		if len(fields) < 1 {
			return nil, fmt.Errorf("invalid input format: '%s' at line %d", line, lineNr)
		}

		cityName := fields[0]
		if _, exists := intermediateCities[cityName]; exists {
			return nil, fmt.Errorf("duplicate city name: '%s' at line %d", cityName, lineNr)
		}
		intermediateCities[cityName] = fields[1:]
		cities[cityName] = &city{name: cityName, neighbors: make(map[direction]*city)}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	for cityName, directions := range intermediateCities {
		for _, neighborField := range directions {
			neighborFields := strings.SplitN(neighborField, "=", 2)
			if len(neighborFields) != 2 {
				return nil, fmt.Errorf(
					"invalid direction '%s' for city '%s'",
					neighborField,
					cityName,
				)
			}

			neighborName := neighborFields[1]
			if _, exists := cities[neighborName]; !exists {
				return nil, fmt.Errorf(
					"unknown neighbor city '%s' in direction '%s' for city '%s'",
					neighborName,
					neighborFields[0],
					cityName,
				)
			}

			dir := direction(neighborFields[0])
			if !dir.isValid() {
				return nil, fmt.Errorf(
					"invalid direction '%s' for city '%s'",
					neighborFields[0],
					cityName,
				)
			}
			cities[cityName].neighbors[dir] = cities[neighborName]
		}
	}

	// validate cities and neighbors
	for _, c := range cities {

		directions := maps.Keys(c.neighbors)
		slices.Sort(directions)

		for _, d := range directions {
			n := c.neighbors[d]
			neighbor, ok := n.neighbors[d.opposite()]
			if !ok || neighbor != c {
				return nil, fmt.Errorf(
					"neighbor city '%s' has no road in direction '%s' to city '%s'",
					n.name,
					d.opposite(),
					c.name,
				)
			}
		}
	}
	return cities, nil
}
