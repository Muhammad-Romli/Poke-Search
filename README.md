# Poke-Search

A command-line Pokedex written in Go. It runs as an interactive REPL, fetches data from [PokeAPI](https://pokeapi.co/), and caches responses so repeated lookups don't hit the network again.

## Features

- **Browse locations**: page forward and backward through location areas, 20 at a time.
- **Explore**: list the Pokemon that can be encountered in a location area.
- **Catch**: throw a Pokeball at a Pokemon. Success is random and harder for Pokemon with higher base experience.
- **Inspect**: view the name, height, weight, stats, and types of a Pokemon you caught.
- **Pokedex**: list every Pokemon you have caught in the current session.
- **Response caching**: API responses are cached in memory, keyed by request URL.

## Requirements

- Go 1.26.5 or newer
- An internet connection (PokeAPI is a public API, no API key needed)

## Installation

```bash
git clone https://github.com/muhammad-romli/Poke-Search.git
cd Poke-Search
go build -o bin/poke-search .
```

## Usage

Start the REPL:

```bash
./bin/poke-search
# or, without building:
go run .
```

Then type commands at the prompt.

### Commands

| Command | Description |
|---|---|
| `help` | Show all available commands |
| `map` | Show the next 20 location areas (call again for the next page) |
| `mapb` | Show the previous 20 location areas |
| `explore <location-name>` | List the Pokemon found in a location area |
| `catch <pokemon-name>` | Try to catch a Pokemon |
| `inspect <pokemon-name>` | Show details of a Pokemon you caught |
| `pokedex` | List all the Pokemon you have caught |
| `exit` | Quit the Pokedex |

### Example session

```text
> map
canalave-city-area
eterna-city-area
...

> explore canalave-city-area
tentacool
tentacruel
...

> catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!

> inspect pikachu
Name: pikachu
Height: 4
Weight: 60
Stats:
  -hp: 35
  -attack: 55
  ...
Types:
  - electric

> pokedex
pikachu
```

> Location names come from the `map` output. Pokemon names are lowercase, as PokeAPI uses them.

## How it works

### HTTP

All data comes from PokeAPI v2 over plain `net/http` GET requests:

- `GET /api/v2/location-area/?limit=20&offset=0`: paginated location list (`map`, `mapb`)
- `GET /api/v2/location-area/{name}`: Pokemon encounters for one location (`explore`)
- `GET /api/v2/pokemon/{name}`: Pokemon details (`catch`, `inspect`)

A single generic helper, `getRequest[T]`, handles every request. It checks the cache first, then makes the HTTP call, rejects any status code above 299, reads the body, unmarshals the JSON into the target struct, and stores the raw response in the cache.

### Caching

Responses are cached as raw JSON bytes, keyed by the full request URL. When a URL was already requested, the cached bytes are unmarshaled directly and no network call is made. This makes going back and forth with `map` / `mapb`, or re-exploring a location, instant after the first fetch.

The cache lives in the `pokecache` package (`internal/pokecache`). Each entry is deleted after 1 minute, and a mutex protects the cache so it is safe for concurrent use. To change how long entries live, edit the interval in `main.go`.

### Config

A shared `config` struct is passed to every command callback. It holds:

- `Next` / `Previous`: the pagination URLs returned by PokeAPI, used by `map` and `mapb`
- `cache`: the response cache
- `pokemonsCaught`: the names of the Pokemon caught in this session

Commands are registered in a map (`AllCommands`) of name, description, and callback, so adding a new command means writing one function and adding one map entry. `help` is generated from that same map.

### Catch mechanic

The chance to catch is `100 - (baseExperience / 10)` percent. Pokemon with higher base experience are harder to catch.

## Notes

- Caught Pokemon are stored in memory only. They are lost when you exit.
- `help` prints commands in random order (Go map iteration order).

## Acknowledgements

Data provided by [PokeAPI](https://pokeapi.co/).
