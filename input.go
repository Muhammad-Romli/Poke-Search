package main

import (
	"strings"
)

func cleanInput(text string) []string {
	words := strings.Split(text, " ")
	newWords := []string{}
	for _, word := range words {
		word = strings.TrimSpace(word)
		word = strings.ToLower(word)
		if word == "" {
			continue
		}
		newWords = append(newWords, word)
	}
	return newWords
}
