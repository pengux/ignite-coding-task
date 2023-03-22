package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"codingtask/simulation"
)

func main() {
	// First argument is the number of aliens to run
	nrOfAliens, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("failed to parse number of aliens: %v", err)
	}

	if nrOfAliens < 2 {
		log.Fatalf("number of aliens must be greater than 1")
	}

	// Read and parse input from STDIN
	sim, err := simulation.NewSimulation(os.Stdin, nrOfAliens)
	if err != nil {
		log.Fatalf("failed to create simulation: %v", err)
	}

	// Run simulation
	state, err := sim.Run()
	if err != nil {
		log.Fatalf("failed to run simulation: %v", err)
	}

	log.Printf("simulation ended with state: %s", state)

	// Print simulation result
	fmt.Println(sim.CitiesToString(sim.SurvivedCities()))
}
