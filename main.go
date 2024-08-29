package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	helpMessage := `Help: 
	Using help you can get help
	Using exit you can exit
	Using map to get 20 pokemon location
	Using mapb to get previous 20 pokemon location
	Using explore <city> to explore a specific city
	Using catch <pokemon name> to catch a pokemon
	Using inspect <pokemon name> to inspect stats of a pokemon
	`
	// Creating the cache to store the pokemon location for 10 seconds
	Cache := NewCache(60)

	var location PokemonLocation

	user := User{
		PokemonMap: make(map[string]Pokemon), // Properly initialize the map
	}

	// Start an infinite loop to read user input
	for {
		// Prompt to read user input
		fmt.Print("Pokedex > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Trim white space
		input = strings.TrimSpace(input)

		// Parse the input to get command and arguments
		command, args := ParseInput(input)

		// Handle each command with a function
		switch command {
		case "help":
			fmt.Println(helpMessage)
		case "exit":
			fmt.Println("Exiting...")
			return
		case "map":
			HandleMapCommand(&location, Cache)
		case "mapb":
			HandleMapbCommand(&location, Cache)
		case "explore":
			if len(args) < 1 {
				fmt.Println("Usage: explore <city>")
			} else {
				HandleExploreCommand(args[0], Cache)
			}
		case "catch":
			if len(args) < 1 {
				fmt.Println("Usage: catch <pokemon>")
			} else {
				HandleCatchCommand(args[0], Cache, &user)
			}
		case "inspect":
			if len(args) < 1 {
				fmt.Println("Usage: inspect <pokemon>")
			} else {
				handleInspectCommand(args[0], &user)
			}
		case "pokedex":
			HandlePokedexCommand(&user)

		default:
			fmt.Println("Unknown command:", command)
			fmt.Println(helpMessage)
		}
	}
}
