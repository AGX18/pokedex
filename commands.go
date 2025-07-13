package main

import "fmt"

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
