package commands

import (
	"fmt"
	"github.com/venzy/pokedexcli/internal/pokeapi"
	"os"
	"sync"
)

type CliCommandConfig struct {
	Arguments []string
	Previous *string
	Next *string
}

type CliCommand struct {
	Name string
	Description string
	Callback func(config *CliCommandConfig) error
}

type Registry map[string]CliCommand
var registryInstance *Registry
var registryOnce sync.Once

func GetRegistry() *Registry {
	registryOnce.Do(func() {
		registryInstance = &Registry{
			"exit": {
				Name: "exit",
				Description: "Exit the Pokedex",
				Callback: commandExit,
			},
			"help": {
				Name: "help",
				Description: "Displays a help message",
				Callback: commandHelp,
			},
			"map": {
				Name: "map",
				Description: "Displays map locations 20 at a time",
				Callback: commandMapNext,
			},
			"mapb": {
				Name: "mapb",
				Description: "Displays previous 20 map locations",
				Callback: commandMapBack,
			},
			"explore": {
				Name: "explore",
				Description: "Returns a list of all Pokémon in a given location",
				Callback: commandExplore,
			},
		}
	})
	return registryInstance
}

func commandExit(config *CliCommandConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *CliCommandConfig) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range *GetRegistry() {
		fmt.Printf("%s: %s\n", command.Name, command.Description)
	}
	return nil
}

func commandMapNext(config *CliCommandConfig) error {
	if config.Next == nil {
		if config.Previous != nil {
			fmt.Println("You're on the last page")
			return nil
		}
	}

	data, err := pokeapi.GetLocationAreas(config.Next)
	if err != nil {
		return err
	}

	config.Previous = data.Previous
	config.Next = data.Next

	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapBack(config *CliCommandConfig) error {
	if config.Previous == nil {
		if config.Next == nil {
			fmt.Println("Must use 'map' command at least once before 'mapb'")
		} else {
			fmt.Println("You're on the first page")
		}
		return nil
	}

	data, err := pokeapi.GetLocationAreas(config.Previous)
	if err != nil {
		return err
	}

	config.Previous = data.Previous
	config.Next = data.Next

	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandExplore(config *CliCommandConfig) error {
	if len(config.Arguments) != 1 {
		return fmt.Errorf("explore command expects 1 argument, the area name")
	}
	areaName := config.Arguments[0]

	detail, err := pokeapi.GetLocationAreaDetail(areaName)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", areaName)
	for _, encounter := range detail.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}