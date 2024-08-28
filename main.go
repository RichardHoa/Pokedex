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
	PokemonLocationCache := NewCache(10)

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
			// fmt.Println("Before deferred increment, Page:", result.Page)
			// fmt.Printf("Len of the result: %d\n", len(result.Results))
			if len(result.Results) == 0 {
				// First time fetching the url
				URL = "https://pokeapi.co/api/v2/location/"
			} else {
				URL = result.Next
			}
			// fmt.Printf("We are handling: %s\n", URL)
			objectExist := PokemonLocationCache.Get(URL, &result)
			if !objectExist{
				fmt.Printf("Url: %s has never been cached before\n", URL)
				FetchLocation(URL, &result)
				PokemonLocationCache.Add(URL, result)
				fmt.Printf("Url: %s has been cached\n", URL)
			} else {
				fmt.Printf("URL: %s has been cached\n", URL)
			}
			
			for _, value := range result.Results {
				fmt.Println(value.Name)
			}
			result.Page++
			// fmt.Println("After deferred increment, Page:", result.Page)
			
			// fmt.Printf("Result previous: %s\n", result.Previous)

		case "mapb":

			// fmt.Println("Before deferred increment, Page:", result.Page)
			
			currentPage := result.Page
			if result.Page == 1 || result.Page == 0 {
				fmt.Println("No previous page")
				continue
			}	
			URL := result.Previous

			objectExist := PokemonLocationCache.Get(URL, &result)
			if !objectExist{
				fmt.Printf("Url: %s has never been cached before\n", URL)
				FetchLocation(URL, &result)
				PokemonLocationCache.Add(URL, result)
				fmt.Printf("Url: %s has been cached\n", URL)

			} else {
				fmt.Printf("URL: %s has been cached\n", URL)
				result.Page = currentPage
			}


			for _, value := range result.Results {
				fmt.Println(value.Name)
			}
			result.Page--
			// fmt.Println("After deferred increment, Page:", result.Page)

		default:
			fmt.Println("Unknown command:", input)
			fmt.Println(helpMessage)
		}
	}

}
