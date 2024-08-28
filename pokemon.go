package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"

)

type PokemonLocation struct {
	Page     int
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []LocationResult `json:"results"`
}

type LocationResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}



func HandleMapCommand(location *PokemonLocation, cache *Cache) {
	var URL string
	var data []byte

	if len(location.Results) == 0 {
		URL = "https://pokeapi.co/api/v2/location/"
	} else {
		URL = location.Next
	}

	cachedData, objectExist := cache.Get(URL)
	if !objectExist {
		fmt.Printf("URL: %s has never been cached before\n", URL)
		data = FetchLocationData(URL)
		if data != nil {
			cache.Add(URL, data)
			fmt.Printf("URL: %s has been cached\n", URL)
		}
	} else {
		fmt.Printf("URL: %s has been cached\n", URL)
		data = cachedData
	}

	if data == nil {
		fmt.Println("Error: No data available")
		return
	}

	// Decode the data into location
	err := json.Unmarshal(data, location)
	if err != nil {
		fmt.Println("Error parsing cached data into object:", err)
		return
	}

	// Print the value of each location
	for _, value := range location.Results {
		fmt.Println(value.Name)
	}
	location.Page++
}

func HandleMapbCommand(location *PokemonLocation, cache *Cache) {
	var data []byte
	currentPage := location.Page
	if location.Page == 1 || location.Page == 0 {
		fmt.Println("No previous page")
		return
	}

	URL := location.Previous
	cachedData, objectExist := cache.Get(URL)
	if !objectExist {
		fmt.Printf("URL: %s has never been cached before\n", URL)
		data = FetchLocationData(URL)
		if data != nil {
			cache.Add(URL, data)
			fmt.Printf("URL: %s has been cached\n", URL)
		}
	} else {
		fmt.Printf("URL: %s has been cached\n", URL)
		data = cachedData
	}

	if data == nil {
		fmt.Println("Error: No data available")
		return
	}

	// Decode the data into location
	err := json.Unmarshal(data, location)
	location.Page = currentPage
	if err != nil {
		fmt.Println("Error parsing cached data into object:", err)
		return
	}

	// Print the value of each location
	for _, value := range location.Results {
		fmt.Println(value.Name)
	}
	location.Page--
}



// HandleExploreCommand handles the "explore" command with a city argument
func HandleExploreCommand(city string, cache *Cache) {
	URL := "https://pokeapi.co/api/v2/location/" + city

	areaUrl, err := FetchAndExtractAreaURL(URL, cache)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	pokemonList, err := FetchPokemonNames(areaUrl, cache)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Found pokemon: ")
	for _, name := range pokemonList {
		fmt.Printf("- %s\n", name)
	}
}

func HandleCatchCommand(pokemonName string, cache *Cache) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	pokemonData, err := FetchPokemonData(url, cache)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// fmt.Printf("Pokemon: %s\n", pokemonData)

	// Extract area URL from the body
	baseExperience := gjson.GetBytes(pokemonData, "base_experience")
	if !baseExperience.Exists() {
		fmt.Println("Error: Base experience not found")
	}
	fmt.Printf("Base experience: %d\n", baseExperience.Int())



}

