package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"

	"github.com/AGX18/pokedex/internal/pokecache"
)

func commandInspect(config *Config) error {
	if config.PokemonName == "" {
		fmt.Println("Please provide a Pokemon name to inspect.")
		return nil
	}

	if pokemon, found := config.Pokedex[config.PokemonName]; found {
		printInfo(pokemon)

	} else {
		fmt.Println("you have not caught that pokemon")
	}

	return nil
}

func printInfo(pokemon Pokemon) {
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("- %s\n", t.Type.Name)
	}
}

func commandPokedex(config *Config) error {
	if len(config.Pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Catch some Pokemon first!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range config.Pokedex {
		fmt.Printf("- %s\n", name)
	}

	return nil
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
	err := GetWithCache(config.NextURL, cache, &response)
	if err != nil {
		return fmt.Errorf("error fetching location areas: %w", err)
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
	err := GetWithCache(config.PrevURL, cache, &response)
	if err != nil {
		return fmt.Errorf("error fetching previous location areas: %w", err)
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
	var area LocationArea
	err := GetWithCache(url, cache, &area)
	if err != nil {
		return fmt.Errorf("error fetching area data: %w", err)
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
		fmt.Println("You may now inspect it with the inspect command.")
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
