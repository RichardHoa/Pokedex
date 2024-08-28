package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	// "net/url"
	// "strconv"
)

type PokemonLocation struct {
	Page     int
	Count    int          `json:"count"`
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Results  []LocationResult `json:"results"`
}


type LocationResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}


func FetchLocation(URL string, result *PokemonLocation)  {
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


// func GetOffsetFromURL(urlString string) int {
// 	// Parse the URL string
// 	parsedURL, err := url.Parse(urlString)
// 	if err != nil {
// 		fmt.Println("Error parsing URL:", err)
// 		return 0
// 	}

// 	// Extract query parameters
// 	queryParams := parsedURL.Query()

// 	// Get the offset value, if present
// 	offsetStr := queryParams.Get("offset")
// 	if offsetStr == "" {
// 		return 0
// 	}

// 	// Convert the offset value to an integer
// 	// offset, err := strconv.Atoi(offsetStr)
// 	// if err != nil {
// 	// 	fmt.Println("Error converting offset to integer:", err)
// 	// 	return 0
// 	// }

// 	return 1
// }