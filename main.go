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
		firstWord := getFirstWord(input)
		fmt.Printf("Your command was: %s\n", firstWord)
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
