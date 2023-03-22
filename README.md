# Alien Invasion Simulation

[Background](task_description.pdf)

## Assumptions
- City names cannot have whitespace in it
- Two cities that are connected to each other must have roads to each other in the opposite directions, e.g. `Foo north=Bar\nBar south=Foo`
- A city cannot have duplicates in the map input
- For each iteration, an alien can either move OR stay at the same city. This is to avoid the situation when there are 2 aliens left and each one is in a city that is direct connected to each other, thus making the simulation runs forever.

## Design choices
- Input data for the map is accepted through STDIN and the number of aliens for the simulation is expected to be the first argument
- Only Simulation struct with its methods are public, all other types and methods are private
- Use of `golang.org/x/exp` for generic functions
- Use recursive depth-first search to check for the case where all aliens are isolated from each other.

## Usage example

```sh
cat testdata/input.txt | go run ./main.go 10 > output.txt
cat output.txt
````

## Fmt, lint, test and coverage

```sh
make fmt
make lint
make test
make test-coverprofile
make show-coverage
```

## Can be improved
- Add a fuzz test for the `parseInput` function
- Try to increase test coverage
