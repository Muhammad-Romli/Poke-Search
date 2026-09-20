package main

type config struct {
	commands map[string]cliCommand
	Next     string
	Previous string
}

func main() {
	cfg := &config{
		commands: AllCommands,
	}
	startRepl(cfg)
}
