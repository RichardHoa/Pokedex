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
	// Creating the cache to store the pokemon location for 10 seconds
	PokemonLocationCache := NewCache(10)

	var location PokemonLocation

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

		switch input {
		// If user seek help, display help message
		case "help":
			fmt.Println(helpMessage)
		// Exit the program
		case "exit":
			fmt.Println("Exiting...")
			return
		// Get 20 pokemon locations with cahe
		case "map":
			var URL string

			// For fetching in the first time, use the default url
			// The second times use the next url in the respond
			if len(location.Results) == 0 {
				URL = "https://pokeapi.co/api/v2/location/"
			} else {
				URL = location.Next
			}
			// Check if the url has been cached
			objectExist := PokemonLocationCache.Get(URL, &location)
			// If it's not cached, fetch it and cache it
			if !objectExist{
				fmt.Printf("Url: %s has never been cached before\n", URL)
				FetchLocation(URL, &location)
				PokemonLocationCache.Add(URL, location)
				fmt.Printf("Url: %s has been cached\n", URL)
			// If it has been cached, do nothing because the location has been updated
			} else {
				fmt.Printf("URL: %s has been cached\n", URL)
			}
			
			// Print the value of each location
			for _, value := range location.Results {
				fmt.Println(value.Name)
			}
			// Increment the page number 
			location.Page++

		// Goes back to the previous 20 pokemon locations
		case "mapb":
			// Current page to keep track of the current page 
			currentPage := location.Page
			// If the page is 1 or 0, there is no previous page
			if location.Page == 1 || location.Page == 0 {
				fmt.Println("No previous page")
				continue
			}	

			// Get the previous url
			URL := location.Previous
			// Check if the url respond has been check
			objectExist := PokemonLocationCache.Get(URL, &location)
			// If it's not cached, fetch it and cache it
			if !objectExist{
				fmt.Printf("Url: %s has never been cached before\n", URL)
				FetchLocation(URL, &location)
				PokemonLocationCache.Add(URL, location)
				fmt.Printf("Url: %s has been cached\n", URL)
			// If it has been cached, set the location.Page to the current page
			// If do not set the location get messed up
			} else {
				fmt.Printf("URL: %s has been cached\n", URL)
				location.Page = currentPage
			}
			// Print the value of each location
			for _, value := range location.Results {
				fmt.Println(value.Name)
			}
			// Decrease the location page
			location.Page--

		// If user input is unknown, display help message
		default:
			fmt.Println("Unknown command:", input)
			fmt.Println(helpMessage)
		}
	}

}
