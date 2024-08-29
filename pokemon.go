package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"math/rand"
	"time"
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

func HandleCatchCommand(pokemonName string, cache *Cache, user *User) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	var pokemonData []byte

	_, pokemonCatched := user.GetPokemon(pokemonName)
	if pokemonCatched {
		fmt.Println("Pokemon already catched, please choose a different pokemon")
		return
	}

	pokemonData, err := FetchPokemonData(url, cache)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Extract area URL from the body
	baseExperience := gjson.GetBytes(pokemonData, "base_experience")
	baseExperientInt := baseExperience.Int()
	if !baseExperience.Exists() {
		fmt.Println("Error: Base experience not found")
		return
	}

	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)

	// If a pokemon has the experience below 20, it will be caught
	// randomFloat := (1 / (float64(baseExperientInt))) * 20
	randomFloat := 1.00
	ceilingRandomNumber := int(float64(baseExperientInt) / randomFloat)
	randomNumber := random.Intn(ceilingRandomNumber)

	fmt.Printf("You have %f chance of catching a pokemon\n", randomFloat*100)
	if randomNumber < int(baseExperientInt) {
		fmt.Println("You caught a pokemon!")
		user.AddPokemon(pokemonName, pokemonData)
	} else {
		fmt.Println("Opps, Try again!")
	}

	fmt.Println("----------------------")
	fmt.Println("Your pokemons list: ")
	for k := range user.PokemonMap {
		fmt.Printf("Pokemon - %s\n", k)
	}

	// catch golduck

}

func handleInspectCommand(pokemonName string,  user *User) {

	pokemonData, pokemonCatched := user.GetPokemon(pokemonName)
	if !pokemonCatched {
		fmt.Println("You can't view stats of a pokemon you haven't catched yet")
		return
	}
	if pokemonData == nil {
		fmt.Println("Error: No data available, please catch the pokemon again")
		return
	}

	// Parse height and weight
	height := gjson.GetBytes(pokemonData, "height").Int()
	weight := gjson.GetBytes(pokemonData, "weight").Int()

	// Extract specific stats using gjson
	hp := gjson.GetBytes(pokemonData, `stats.#(stat.name=="hp").base_stat`).Int()
	attack := gjson.GetBytes(pokemonData, `stats.#(stat.name=="attack").base_stat`).Int()
	defense := gjson.GetBytes(pokemonData, `stats.#(stat.name=="defense").base_stat`).Int()
	specialAttack := gjson.GetBytes(pokemonData, `stats.#(stat.name=="special-attack").base_stat`).Int()
	specialDefense := gjson.GetBytes(pokemonData, `stats.#(stat.name=="special-defense").base_stat`).Int()
	speed := gjson.GetBytes(pokemonData, `stats.#(stat.name=="speed").base_stat`).Int()
	// Extract all types
	types := gjson.GetBytes(pokemonData, "types.#.type.name").Array()

	fmt.Printf("Name: %s\n", pokemonName)
	fmt.Printf("Height: %d\n", height)
	fmt.Printf("Weight: %d\n", weight)

	fmt.Println("Stats: ")
	fmt.Printf("  - HP: %d\n", hp)
	fmt.Printf("  - Attack: %d\n", attack)
	fmt.Printf("  - Defense: %d\n", defense)
	fmt.Printf("  - Special Attack: %d\n", specialAttack)
	fmt.Printf("  - Special Defense: %d\n", specialDefense)
	fmt.Printf("  - Speed: %d\n", speed)

	fmt.Println("Types: ")
	for _, t := range types {
		fmt.Printf("  - Type: %s\n", t)
	}

}
