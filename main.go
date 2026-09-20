package main

type config struct {
	commands map[string]cliCommand
}

func main() {
	cfg := &config{
		commands: AllCommands,
	}
	startRepl(cfg)
}
