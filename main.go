package main

import (
	"time"

	"github.com/muhammad-romli/Poke-Search/internal/pokecache"
)

type config struct {
	cache          pokecache.Cache
	commands       map[string]cliCommand
	pokemonsCaught []string
	Next           string
	Previous       string
}

func main() {
	cfg := &config{
		cache:    *pokecache.NewCache(IntervalTimer * time.Second),
		commands: AllCommands,
	}
	startRepl(cfg)
}

