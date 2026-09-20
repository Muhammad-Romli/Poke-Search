package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		firstWord := getFirstWord(input)
		trimmedFirstWord := strings.Trim(firstWord, " ")
		finishedFirstWord := strings.ToLower(trimmedFirstWord)

		for commandName, commandStruct := range AllCommand {
			if finishedFirstWord == commandName {
				if err := commandStruct.callback(); err != nil {
					fmt.Println("Error running the command: %w", err)
					os.Exit(0)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(fmt.Errorf("failed reading scan: %w", err))
	}
}

func getFirstWord(s string) string {
	for i := range s {
		if s[i] == ' ' {
			return s[:i]
		}
	}
	return s
}
