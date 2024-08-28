package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"io"
	"github.com/tidwall/gjson"
	// "net/url"
	// "strconv"
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

func FetchLocation(URL string, result *PokemonLocation) {
	respond, err := http.Get(URL)
	if err != nil {
		fmt.Println("Error happen when fetching data")
		fmt.Println(err)
	}
	defer respond.Body.Close()

	err = json.NewDecoder(respond.Body).Decode(result)
	if err != nil {
		fmt.Println("Error happen when parsing data into object")
		fmt.Println(err)
	}
}

// handleMapCommand handles the "map" command logic
func HandleMapCommand(location *PokemonLocation, cache *Cache) {
	var URL string

	// For fetching the first time, use the default url
	// The second time use the next url in the response
	if len(location.Results) == 0 {
		URL = "https://pokeapi.co/api/v2/location/"
	} else {
		URL = location.Next
	}
	// Check if the URL has been cached
	objectExist := cache.GetLocation(URL, location)
	// If it's not cached, fetch it and cache it
	if !objectExist {
		fmt.Printf("URL: %s has never been cached before\n", URL)
		FetchLocation(URL, location)
		cache.AddLocation(URL, location)
		fmt.Printf("URL: %s has been cached\n", URL)
	} else {
		fmt.Printf("URL: %s has been cached\n", URL)
	}

	// Print the value of each location
	for _, value := range location.Results {
		fmt.Println(value.Name)
	}
	// Increment the page number
	location.Page++
}

// handleMapbCommand handles the "mapb" command logic
func HandleMapbCommand(location *PokemonLocation, cache *Cache) {
	// Current page to keep track of the current page
	currentPage := location.Page
	// If the page is 1 or 0, there is no previous page
	if location.Page == 1 || location.Page == 0 {
		fmt.Println("No previous page")
		return
	}

	// Get the previous URL
	URL := location.Previous
	// Check if the URL has been cached
	objectExist := cache.GetLocation(URL, location)
	// If it's not cached, fetch it and cache it
	if !objectExist {
		fmt.Printf("URL: %s has never been cached before\n", URL)
		FetchLocation(URL, location)
		cache.AddLocation(URL, location)
		fmt.Printf("URL: %s has been cached\n", URL)
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
}

func ParseInput(input string) (string, []string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

// handleExploreCommand handles the "explore" command with a city argument
func HandleExploreCommand(city string) {
	// Implement the logic for exploring a specific city
	URL := "https://pokeapi.co/api/v2/location/" + city
	// fmt.Printf("Exploring city: %s\n", URL)
	// Fetch and process the city information
	areaUrl, err := fetchAndExtractAreaURL(URL)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// fmt.Println("Area URL:", areaUrl)
	pokemonList, err := fetchPokemonNames(areaUrl)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Found pokemon: ")
	for _, name := range pokemonList {
		fmt.Printf("- %s\n", name)
	}
	// Add your explore city logic here
}

func fetchAndExtractAreaURL(url string) (areaURL string, err error) {
	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "",fmt.Errorf("failed to fetch URL: %v", err)

	}
	defer resp.Body.Close()

	// Check for a successful response status
	if resp.StatusCode != http.StatusOK {
		return "",fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "",fmt.Errorf("failed to read response body: %v", err)
	}

	result := gjson.GetBytes(body, "areas.0.url")

	if !result.Exists() {
		return "", fmt.Errorf("area URL not found in JSON")
	}

	return result.String(), nil
}


func fetchPokemonNames(url string) ([]string, error) {
	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	// Check for a successful response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Extract Pokémon names using gjson
	var pokemonNames []string
	result := gjson.GetBytes(body, "pokemon_encounters.#.pokemon.name")
	result.ForEach(func(_, value gjson.Result) bool {
		pokemonNames = append(pokemonNames, value.String())
		return true
	})

	return pokemonNames, nil
}