package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"github.com/tidwall/gjson"
)



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



func ParseInput(input string) (string, []string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}


func FetchAndExtractAreaURL(url string, cache *Cache) (areaURL string, err error) {
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

func FetchPokemonNames(url string, cache *Cache) ([]string, error) {
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

func FetchPokemonData(url string, cache *Cache) ([]byte, error) {

	var pokemonData []byte

	// Check if data is in cache
	pokemonData, dataExist := cache.Get(url)

	if !dataExist {
		fmt.Printf("URL has never been cached: %s\n", url)
		// Data is not cached, fetch it from the URL
		resp, err := http.Get(url)
		if err != nil {
			return []byte{}, fmt.Errorf("failed to fetch URL: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return []byte{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		pokemonData, err = io.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, fmt.Errorf("failed to read response body: %v", err)
		}
		// Add the newly fetched data to the cache
		cache.Add(url, pokemonData)
	} else {
		fmt.Printf("Url has been cached %s\n", url)
	}

	return pokemonData, nil


}