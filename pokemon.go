package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

func FetchLocationData(URL string) []byte {
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return nil
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	// fmt.Printf("data: %s\n", data)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	return data
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

func ParseInput(input string) (string, []string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

// HandleExploreCommand handles the "explore" command with a city argument
func HandleExploreCommand(city string, cache *Cache) {
	URL := "https://pokeapi.co/api/v2/location/" + city

	areaUrl, err := fetchAndExtractAreaURL(URL, cache)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	pokemonList, err := fetchPokemonNames(areaUrl, cache)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Found pokemon: ")
	for _, name := range pokemonList {
		fmt.Printf("- %s\n", name)
	}
}

func fetchAndExtractAreaURL(url string, cache *Cache) (areaURL string, err error) {
	var body []byte

	// Check if data is in cache
	body, dataExist := cache.Get(url)
	if !dataExist {
		fmt.Printf("URL has never been cached: %s\n", url)
		// Data is not cached, fetch it from the URL
		resp, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("failed to fetch URL: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("unexpected status code: %d | Possibly wrong location name", resp.StatusCode)
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %v", err)
		}

		// Add the newly fetched data to the cache
		cache.Add(url, body)
	} else {
		fmt.Printf("Url has been cached %s\n", url)
	}

	// Extract area URL from the body
	result := gjson.GetBytes(body, "areas.0.url")
	if !result.Exists() {
		return "", fmt.Errorf("area URL not found in JSON")
	}

	return result.String(), nil
}

func fetchPokemonNames(url string, cache *Cache) ([]string, error) {
	var body []byte

	// Check if data is in cache
	body, dataExist := cache.Get(url)

	if !dataExist {
		fmt.Printf("URL has never been cached: %s\n", url)
		// Data is not cached, fetch it from the URL
		resp, err := http.Get(url)
		if err != nil {
			return []string{}, fmt.Errorf("failed to fetch URL: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return []string{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return []string{}, fmt.Errorf("failed to read response body: %v", err)
		}
		// Add the newly fetched data to the cache
		cache.Add(url, body)
	} else {
		fmt.Printf("Url has been cached %s\n", url)
	}

	var pokemonNames []string
	result := gjson.GetBytes(body, "pokemon_encounters.#.pokemon.name")
	result.ForEach(func(_, value gjson.Result) bool {
		pokemonNames = append(pokemonNames, value.String())
		return true
	})

	return pokemonNames, nil
}
