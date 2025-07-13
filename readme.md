# Pokedex 
Pokedex  in a command-line REPL. used the incredible PokéAPI to fetch all of the data using GET requests. A Pokedex is just a make-believe device that lets us look up information about Pokemon - things like their name, type, and stats.

## Commands
- exit: Exit the Pokedex
- help: Displays a help message
- map: Fetches the map of locations
- mapb: Fetches the previous map of locations
- explore: Explore a specific location area by name
- catch: Catch a specific Pokemon by name
- inspect: Inspect a specific Pokemon by name
- pokedex: Display all caught Pokemon

## Ideas for Extending the Project
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
