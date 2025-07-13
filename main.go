package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
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
		} else if command.name == "catch" {
			if len(words) < 2 {
				fmt.Println("Please provide a Pokemon name to catch.")
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

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: pokedex [command]")
	for cmd, command := range supportedCommands {
		fmt.Printf("- %s: %s\n", cmd, command.description)
	}
	return nil
}

func commandMap(config *Config) error {
	if config.NextURL == "" {
		fmt.Println("No more location areas available.")
		return nil
	}
	// Check if the URL is cached
	var response LocationAreaListResponse
	cachedData, found := cache.Get(config.NextURL)
	if found {
		err := json.Unmarshal(cachedData, &response)
		if err != nil {
			return fmt.Errorf("error decoding cached JSON: %w", err)
		}

	} else {
		res, err := http.Get(config.NextURL)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return fmt.Errorf("error decoding JSON: %w", err)
		}
	}

	// Update config URLs
	if response.Previous != nil {
		config.PrevURL = *response.Previous
	} else {
		config.PrevURL = ""
	}
	if response.Next != nil {
		config.NextURL = *response.Next
	} else {
		config.NextURL = ""
	}
	displayLocationAreas(response.Results)
	return nil
}

func commandMapBack(config *Config) error {
	if config.PrevURL == "" {
		fmt.Println("No previous location areas available.")
		return nil
	}
	var response LocationAreaListResponse
	cachedData, found := cache.Get(config.NextURL)
	if found {
		err := json.Unmarshal(cachedData, &response)
		if err != nil {
			return fmt.Errorf("error decoding cached JSON: %w", err)
		}

	} else {
		res, err := http.Get(config.PrevURL)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return fmt.Errorf("error decoding JSON: %w", err)
		}

		// Cache the response
		data, err := json.Marshal(response)
		if err != nil {
			return fmt.Errorf("error encoding JSON for cache: %w", err)
		}
		cache.Add(config.NextURL, data)
	}

	// Update config URLs
	if response.Previous != nil {
		config.PrevURL = *response.Previous
	} else {
		config.PrevURL = ""
	}
	if response.Next != nil {
		config.NextURL = *response.Next
	} else {
		config.NextURL = ""
	}
	displayLocationAreas(response.Results)
	return nil
}

func displayLocationAreas(locationAreas []LocationAreaSummary) {
	if len(locationAreas) == 0 {
		fmt.Println("No more location areas found.")
		return
	}
	for _, area := range locationAreas {
		fmt.Println(area.Name)
	}
}

func commandExplore(config *Config) error {
	if config.AreaName == "" && config.AreaID == 0 {
		fmt.Println("Please provide an area name or ID to explore.")
		return nil
	}

	var url string
	if config.AreaName != "" {
		url = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", config.AreaName)
		fmt.Printf("Exploring %s...\n", config.AreaName)
	} else {
		url = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d", config.AreaID)
	}
	// Check if the URL is cached
	cachedData, found := cache.Get(url)
	var area LocationArea
	if found {
		err := json.Unmarshal(cachedData, &area)
		if err != nil {
			return fmt.Errorf("error decoding cached JSON: %w", err)
		}
		fmt.Printf("Using cached data for %s\n", area.Name)
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error fetching location area: %w", err)
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&area)
		if err != nil {
			return fmt.Errorf("error decoding JSON: %w", err)
		}
		data, err := json.Marshal(area)
		if err != nil {
			return fmt.Errorf("error encoding JSON for cache: %w", err)
		}
		cache.Add(url, data)
	}

	// Cache the response

	fmt.Printf("Found Pokemon:\n")
	for _, encounter := range area.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *Config) error {
	if config.PokemonName == "" {
		fmt.Println("Please provide a Pokemon name to catch.")
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", config.PokemonName)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", config.PokemonName)
	// Check if the URL is cached
	var pokemon Pokemon
	err := GetWithCache(url, cache, &pokemon)
	if err != nil {
		return err
	}

	CatchProbability := catchProbability(pokemon.BaseExperience)

	if CatchProbability >= 5 {
		fmt.Printf("Caught %s!\n", pokemon.Name)
		config.Pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func catchProbability(max int) int {
	if max <= 1 {
		return 0
	}
	factor := 1.0 / (1 + float64(max)/10.0) // decay function
	return int(float64(max) * rand.Float64() * factor)
}

func GetWithCache[T any](url string, cache *pokecache.Cache, target *T) error {
	cachedData, found := cache.Get(url)
	if found {
		err := json.Unmarshal(cachedData, target)
		if err != nil {
			return fmt.Errorf("error decoding cached JSON: %w", err)
		}
		return nil
	}

	// Fetch from HTTP if not in cache
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching data from %s: %w", url, err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	data, err := json.Marshal(target)
	if err != nil {
		return fmt.Errorf("error encoding JSON for cache: %w", err)
	}

	cache.Add(url, data)
	return nil
}
