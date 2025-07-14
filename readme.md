# Pokedex CLI 
A command-line Pokedex written in Go. This tool interacts with the public [PokeAPI](https://pokeapi.co/) to fetch Pokémon data and caches results locally to reduce redundant API calls.

## Features

- Search and explore Pokémon regions, locations, and Pokémon entries.
- In-memory caching with expiration to optimize API usage.
- Command-line interface with basic commands like `map`, `explore`, and `inspect`.
- Concurrent-safe cache using goroutines and `sync.RWMutex`.

## Available Commands
- exit: Exit the Pokedex
- help: Displays a help message
- map: Fetches the map of locations
- mapb: Fetches the previous map of locations
- explore: Explore a specific location area by name
- catch: Catch a specific Pokemon by name
- inspect: Inspect a specific Pokemon by name
- pokedex: Display all caught Pokemon


## Getting Started

### Prerequisites
- Go 1.20+

### Run the CLI
```bash
go run main.go
```

## Testing
The caching layer is unit-tested using Go’s built-in testing package. Other components were tested manually by interacting with the CLI and observing responses.

## Technologies Used
- Go (Golang)
- RESTful API (PokeAPI)
- Concurrency primitives (goroutines, RWMutex)

## Upcoming Features
- [ ] Update the CLI to support the "up" arrow to cycle through previous commands
- [ ] Simulate battles between pokemon
- [ ] Add more unit tests
- [ ] Refactor your code to organize it better and make it more testable
- [ ] Keep pokemon in a "party" and allow them to level up
- [ ] Allow for pokemon that are caught to evolve after a set amount of time
- [ ] Persist a user's Pokedex to disk so they can save progress between sessions
- [ ] Use the PokeAPI to make exploration more interesting. For example, rather than typing the names of areas, maybe you are given choices of areas and just type "left" or "right"
- [ ] Random encounters with wild pokemon
- [ ] Adding support for different types of balls (Pokeballs, Great Balls, Ultra Balls, etc), which have different chances of catching pokemon
