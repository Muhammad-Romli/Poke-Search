package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var AllCommand = map[string]cliCommand{}

func init() {
	AllCommand = map[string]cliCommand{
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

func commandExit() error {
	fmt.Printf("losing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Printf(`
=================================
Welcome to the Pokedex!
Usage:
`)

	for commandName, commandStruct := range AllCommand {
		fmt.Printf("%s: %s", commandName, commandStruct.description)
	}
	fmt.Println("=================================")

	return nil
}
