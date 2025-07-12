package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := Config{
		NextURL: "https://pokeapi.co/api/v2/location-area/?limit=20&offset=0",
		PrevURL: "",
		Offset:  0,  // Offset for pagination
		Limit:   20, // Default limit for pagination
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

type Config struct {
	// Add configuration fields as needed
	NextURL string
	PrevURL string
	Offset  int
	Limit   int
}

type LocationAreaListResponse struct {
	Count    int                   `json:"count"`
	Next     *string               `json:"next"`
	Previous *string               `json:"previous"`
	Results  []LocationAreaSummary `json:"results"`
}

// For the simple list items
type LocationAreaSummary struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationArea struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
