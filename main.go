package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/michaeldebetaz/pokedexcli/internal/pokeapi"
	"github.com/michaeldebetaz/pokedexcli/internal/pokecache"
)

type config struct {
	Next     string
	Previous string
	Cache    pokecache.Cache
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

		inputs := cleanInput(scanner.Text())

		if len(inputs) < 1 {
			continue
		}

		first := inputs[0]
		cmd, ok := cmds[first]
		if ok {
			cmd.Callback()
		} else {
			fmt.Println("Command not found")
		}
	}
}

func cleanInput(s string) []string {
	return strings.Fields(strings.ToLower(s))
}

func createConfig() *config {
	const INTERVAL = 5 * time.Second
	const LIMIT = 20

	return &config{
		Next:     fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?limit=%d", LIMIT),
		Previous: "",
		Cache:    pokecache.NewCache(INTERVAL),
	}
}

func createCommands(config *config) map[string]cliCommand {
	cmds := make(map[string]cliCommand)

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

	cmds["map"] = cliCommand{
		Name:        "map",
		Description: "",
		Callback: func() error {
			url := config.Next
			if url == "" {
				fmt.Println("you're on the last page")
				return nil
			}

			data := pokeapi.GetLocationAreasData(url, config.Cache)
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

			data := pokeapi.GetLocationAreasData(url, config.Cache)
			config.Next = data.Next
			config.Previous = data.Previous

			for _, area := range data.Results {
				fmt.Printf("%s\n", area.Name)
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
