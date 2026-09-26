package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

var AllCommands = map[string]cliCommand{}

func init() {
	AllCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Find all command available in Pokedex",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Get 20 location areas in the Pokemon world. call it again to move to the next page.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Move to previous 20 location areas page",
			callback:    commandMapB,
		},
	}
}

func commandExit(configP *config, additional string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(configP *config, additional string) error {
	fmt.Printf(`
=================================
Welcome to the Pokedex!
Usage:

`)

	for commandName, commandStruct := range AllCommands {
		fmt.Printf("%s: %s\n", commandName, commandStruct.description)
	}
	fmt.Println("\n=================================")

	return nil
}

type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func getRequest(fullUrl string, locStruct *LocationAreaResponse, config *config) error {
	val, ok := config.cache.Get(fullUrl)
	if ok {
		if err := json.Unmarshal(val, locStruct); err != nil {
			return fmt.Errorf("failed unmarshaling cached data: %w", err)
		}
		return nil
	}

	resp, err := http.Get(fullUrl)

	if err != nil {
		return fmt.Errorf("error making get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to reading response: %w", err)
	}
	if err = json.Unmarshal(data, locStruct); err != nil {
		return fmt.Errorf("failed unmarshaling data, data not cached yet: %w", err)
	}
	config.cache.Add(fullUrl, data)
	return nil
}

func commandMap(configP *config, additional string) error {
	var locStruct LocationAreaResponse
	fullUrl := ""

	if configP.Next == "" {
		page := 1
		shown := 20
		total := (page * shown) - shown
		query := fmt.Sprintf("?limit=%d&offset=%d", shown, total)
		url := "https://pokeapi.co/api/v2/location-area/"
		fullUrl = fmt.Sprintf("%s%s", url, query)
	} else {
		fullUrl = configP.Next
	}

	err := getRequest(fullUrl, &locStruct, configP)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	configP.Next = locStruct.Next
	configP.Previous = locStruct.Previous

	for _, result := range locStruct.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapB(configP *config, additional string) error {
	var locStruct LocationAreaResponse
	if configP.Previous == "" {
		fmt.Printf("You are on the first page you can't go back")
		return nil
	}
	err := getRequest(configP.Previous, &locStruct, configP)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	configP.Next = locStruct.Next
	configP.Previous = locStruct.Previous

	for _, result := range locStruct.Results {
		fmt.Println(result.Name)
	}
	return nil
}
