package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
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
		"explore": {
			name:        "explore",
			description: "Reveal all the pokemons in specified location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Throw a pokeball to try to catch specified pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect pokemon you already caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "See every pokemon you have caught",
			callback:    commandPokedex,
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

type LocationAreaList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaDetailstruct struct {
	ID                int                `json:"id"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonDetailStruct struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`

	Abilities []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`

	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`

	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

type NamedAPIResource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonEncounter struct {
	Pokemon NamedAPIResource `json:"pokemon"`
}

func getCachedData[T any](fullUrl string, dataStruct *T, configP *config) (bool, error) {
	val, ok := configP.cache.Get(fullUrl)
	if !ok {
		return false, nil
	}
	if err := json.Unmarshal(val, dataStruct); err != nil {
		return true, fmt.Errorf("failed unmarshal cached data: %w", err)
	}
	return true, nil
}

func getRequest[T any](fullUrl string, dataStruct *T, configP *config) error {
	exist, err := getCachedData(fullUrl, dataStruct, configP)
	if err != nil {
		return err
	}
	if exist {
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
	if err = json.Unmarshal(data, dataStruct); err != nil {
		return fmt.Errorf("failed unmarshaling data, data not cached yet: %w", err)
	}
	configP.cache.Add(fullUrl, data)
	return nil
}

func commandMap(configP *config, additional string) error {
	var locStruct LocationAreaList
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
	var locStruct LocationAreaList
	if configP.Previous == "" {
		fmt.Printf("Error: You are on the first page you can't go back")
		// This is intended, the prefix error just to warn user, but not kick the user from REPL
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

func commandExplore(configP *config, additional string) error {
	var locDetailStruct LocationAreaDetailstruct
	if additional == "" {
		fmt.Printf("Error: You need to specify the location of the map after the command(e.g: explore <location-name> )\n")
		// This is intended, the prefix error just to warn user, but not kick the user from REPL
		return nil
	}
	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	fullUrl := fmt.Sprintf("%s%s", baseUrl, additional)
	if err := getRequest(fullUrl, &locDetailStruct, configP); err != nil {
		return fmt.Errorf("failed when getting request: %w", err)
	}

	for _, pokemon := range locDetailStruct.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(configP *config, additional string) error {
	var pokeStruct PokemonDetailStruct
	if additional == "" {
		fmt.Printf("Error: You need to specify the pokemon name after the command(e.g: catch <pokemon-name> )\n")
		// This is intended, the prefix error just to warn user, but not kick the user from REPL
		return nil
	}
	baseUrl := "https://pokeapi.co/api/v2/pokemon/"
	fullUrl := fmt.Sprintf("%s%s", baseUrl, additional)
	if err := getRequest(fullUrl, &pokeStruct, configP); err != nil {
		return fmt.Errorf("failed when getting request: %w", err)
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokeStruct.Name)
	// CATCH CHANCE
	baseChance := 100
	expChance := pokeStruct.BaseExperience / 10
	totalChance := baseChance - expChance
	rolled := rand.IntN(baseChance)
	if rolled < totalChance {
		fmt.Printf("%s was caught!\n", pokeStruct.Name)
		configP.pokemonsCaught = append(configP.pokemonsCaught, pokeStruct.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokeStruct.Name)
	}
	return nil
}

func commandInspect(configP *config, additional string) error {
	var pokeStruct PokemonDetailStruct
	if additional == "" {
		fmt.Printf("Error: You need to specify the pokemon name after the command(e.g: catch <pokemon-name> )\n")
		// This is intended, the prefix error just to warn user, but not kick the user from REPL
		return nil
	}
	baseUrl := "https://pokeapi.co/api/v2/pokemon/"
	fullUrl := fmt.Sprintf("%s%s", baseUrl, additional)
	exist, err := getCachedData(fullUrl, &pokeStruct, configP)
	if err != nil {
		return fmt.Errorf("Error when getting cached data for caught pokemon\n")
	}
	if !exist {
		fmt.Printf("you have not caught that pokemon\n")
		return nil
	}

	fmt.Printf(`
Name: %s
Height: %d
Weight: %d
`, pokeStruct.Name, pokeStruct.Height, pokeStruct.Weight)

	fmt.Printf("Stats:")
	for _, s := range pokeStruct.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, t := range pokeStruct.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(configP *config, additional string) error {
	if len(configP.pokemonsCaught) == 0 {
		fmt.Println("You haven't caught any pokemons")
		return nil
	}
	for _, pokemon := range configP.pokemonsCaught {
		fmt.Println(pokemon)
	}
	return nil
}
