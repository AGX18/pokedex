package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		fmt.Printf("Your command was: %s\n", words[0])

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
