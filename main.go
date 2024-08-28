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
	`

	var result PokemonLocation

	for {
		fmt.Print("Pokedex > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		input = strings.TrimSpace(input)
		switch input {
		case "help":
			fmt.Println(helpMessage)
		case "exit":
			fmt.Println("Exiting...")
			return
		case "map":
			var URL string
			if len(result.Results) == 0 {
				// First time fetching the url
				URL = "https://pokeapi.co/api/v2/location/"
			} else {
				URL = result.Next
			}

			FetchLocation(URL, &result)
			result.Page++

			for _, value := range result.Results {
				fmt.Println(value.Name)
			}
			// fmt.Printf("Result previous: %s\n", result.Previous)

		case "mapb":
			if result.Page == 1 {
				fmt.Println("No previous page")
				continue
			}	
			// fmt.Printf("Result previous: %s\n", result.Previous)
			URL := result.Previous

			FetchLocation(URL, &result)
			result.Page--

			for _, value := range result.Results {
				fmt.Println(value.Name)
			}

		default:
			fmt.Println("Unknown command:", input)
			fmt.Println(helpMessage)
		}
	}

}
