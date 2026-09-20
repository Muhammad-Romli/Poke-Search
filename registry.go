package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}
}

func commandExit(configP *config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(configP *config) error {
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
