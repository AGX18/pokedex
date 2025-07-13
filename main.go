package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AGX18/pokedex/internal/pokecache"
)

var supportedCommands map[string]cliCommand

var cache *pokecache.Cache = pokecache.NewCache(5 * time.Second)

func init() {
	supportedCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Fetches the map of locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Fetches the previous map of locations",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore a specific location area by name",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a specific Pokemon by name",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a specific Pokemon by name",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Display all caught Pokemon",
			callback:    commandPokedex,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := Config{
		NextURL:     "https://pokeapi.co/api/v2/location-area/?limit=20&offset=0",
		PrevURL:     "",
		Offset:      0,  // Offset for pagination
		Limit:       20, // Default limit for pagination
		AreaName:    "", // For searching by area name
		AreaID:      0,  // For searching by area ID
		Pokedex:     make(map[string]Pokemon),
		PokemonName: "", // For catching a specific Pokemon
	}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		command, ok := supportedCommands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if command.name == "explore" {
			if len(words) < 2 {
				fmt.Println("Please provide an area name or ID to explore.")
				continue
			}
			config.AreaName = words[1] // Set the area name from the input
		} else if command.name == "catch" || command.name == "inspect" {
			if len(words) < 2 {
				fmt.Println("Please provide a Pokemon name.")
				continue
			}
			config.PokemonName = words[1] // Set the Pokemon name from the input
		} else {
			config.AreaName = ""    // Reset area name for other commands
			config.PokemonName = "" // Reset Pokemon name for other commands
		}
		// Execute the command callback
		err := command.callback(&config) // Call the command's callback function
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

// split the user's input into "words" based on whitespace. It should also lowercase the input and trim any leading or trailing whitespace.
// example: "Hello World" -> ["hello", "world"]
func cleanInput(text string) []string {
	// Trim leading and trailing spaces
	trimmed := strings.TrimSpace(text)
	// Convert to lowercase
	lowered := strings.ToLower(trimmed)
	// Split by whitespace
	words := strings.Fields(lowered)
	return words
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config) error
}
