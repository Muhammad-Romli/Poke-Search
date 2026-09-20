package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		finishedFirstWord := cleanInput(input)[0]

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
		fmt.Println(fmt.Errorf("failed reading scan: %w\n", err))
	}
}
