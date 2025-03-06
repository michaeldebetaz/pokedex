package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/michaeldebetaz/pokedexcli/internal/pokeapi"
	"github.com/michaeldebetaz/pokedexcli/internal/pokecache"
)

type config struct {
	Inputs   []string
	Next     string
	Previous string
	Cache    pokecache.Cache
	Pokedex  map[string]pokeapi.Pokemon
}

type cliCommand struct {
	Name        string
	Description string
	Callback    func() error
}

type cliCommands map[string]cliCommand

func main() {
	config := createConfig()
	cmds := createCommands(config)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()

		config.Inputs = cleanInput(scanner.Text())

		if len(config.Inputs) > 0 {
			first := config.Inputs[0]
			cmd, ok := cmds[first]
			if ok {
				cmd.Callback()
			} else {
				fmt.Println("Command not found")
			}
		}
	}
}

func cleanInput(s string) []string {
	return strings.Fields(strings.ToLower(s))
}

func createConfig() *config {
	return &config{
		Inputs:   []string{},
		Next:     fmt.Sprintf("https://pokeapi.co/api/v2/location-area/"),
		Previous: "",
		Cache:    pokecache.NewCache(5 * time.Second),
		Pokedex:  make(map[string]pokeapi.Pokemon),
	}
}

func createCommands(config *config) map[string]cliCommand {
	cmds := make(map[string]cliCommand)

	cmds["catch"] = cliCommand{
		Name:        "catch",
		Description: "Catch a Pokemon",
		Callback: func() error {
			if len(config.Inputs) < 2 {
				fmt.Println("Please provide a Pokemon name")
				return nil
			}

			pokemonName := config.Inputs[1]
			fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

			url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName + "/"
			pokemon := pokeapi.GetPokemonData(url, config.Cache)

			const BASE = 1000
			n := rand.Intn(BASE)

			if n > (BASE+pokemon.BaseExperience)/2 {
				fmt.Printf("%s was caught!\n", pokemon.Name)
				config.Pokedex[pokemon.Name] = pokemon
			} else {
				fmt.Printf("%s escaped!\n", pokemon.Name)
			}

			return nil
		},
	}

	cmds["explore"] = cliCommand{
		Name:        "explore",
		Description: "Explore a specific location area",
		Callback: func() error {
			if len(config.Inputs) < 2 {
				fmt.Println("Please provide a location area")
				return nil
			}

			areaName := config.Inputs[1]
			fmt.Printf("Exploring %s...\n", areaName)

			url := "https://pokeapi.co/api/v2/location-area/" + areaName + "/"
			data := pokeapi.GetLocationAreaData(url, config.Cache)

			fmt.Println("Found Pokemon:")
			for _, encounter := range data.PokemonEncounters {
				fmt.Printf(" - %s\n", encounter.Pokemon.Name)
			}

			return nil
		},
	}

	cmds["help"] = cliCommand{
		Name:        "help",
		Description: "Displays a help message",
		Callback: func() error {
			fmt.Println("Welcome to the Pokedex!")
			fmt.Println("Usage:")
			fmt.Println()
			for key, cmd := range cmds {
				fmt.Printf("%s: %s\n", key, cmd.Description)
			}
			fmt.Println()
			return nil
		},
	}

	cmds["inspect"] = cliCommand{
		Name:        "inspect",
		Description: "Inspect a Pokemon",
		Callback: func() error {
			if len(config.Inputs) < 2 {
				fmt.Println("Please provide a Pokemon name")
				return nil
			}

			pokemonName := config.Inputs[1]
			fmt.Printf("Inspecting %s in the Pokedex...\n", pokemonName)

			pokemon, ok := config.Pokedex[pokemonName]
			if !ok {
				fmt.Println("you have not caught that pokemon")
				return nil
			}

			fmt.Printf("Name: %s\n", pokemon.Name)
			fmt.Printf("Height: %d\n", pokemon.Height)
			fmt.Printf("Weight: %d\n", pokemon.Weight)

			fmt.Println("Stats:")
			for _, stat := range pokemon.Stats {
				fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
			}

			fmt.Println("Types:")
			for _, t := range pokemon.Types {
				fmt.Printf("  - %s\n", t.Type.Name)
			}

			return nil
		},
	}

	cmds["map"] = cliCommand{
		Name:        "map",
		Description: "",
		Callback: func() error {
			url := config.Next
			if url == "" {
				fmt.Println("you're on the last page")
				return nil
			}

			data := pokeapi.GetLocationAreasPagination(url, config.Cache)
			config.Next = data.Next
			config.Previous = data.Previous

			for _, area := range data.Results {
				fmt.Printf("%s\n", area.Name)
			}

			return nil
		},
	}

	cmds["mapb"] = cliCommand{
		Name:        "mapb",
		Description: "",
		Callback: func() error {
			url := config.Previous
			if url == "" {
				fmt.Println("you're on the first page")
				return nil
			}

			data := pokeapi.GetLocationAreasPagination(url, config.Cache)
			config.Next = data.Next
			config.Previous = data.Previous

			for _, area := range data.Results {
				fmt.Printf("%s\n", area.Name)
			}

			return nil
		},
	}

	cmds["pokedex"] = cliCommand{
		Name:        "pokedex",
		Description: "List all Pokemon in the Pokedex",
		Callback: func() error {
			fmt.Println("Your Pokedex:")

			for name := range config.Pokedex {
				fmt.Printf(" - %s\n", name)
			}

			return nil
		},
	}

	cmds["exit"] = cliCommand{
		Name:        "exit",
		Description: "Exit the Pokedex",
		Callback: func() error {
			fmt.Println("Closing the Pokedex... Goodbye!")
			os.Exit(0)
			return nil
		},
	}

	return cmds
}
