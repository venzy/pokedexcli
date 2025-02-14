package commands

import (
	"fmt"
	"github.com/venzy/pokedexcli/internal/common"
	"github.com/venzy/pokedexcli/internal/pokeapi"
	"os"
	"sync"
)

type CliCommand struct {
	Name string
	Description string
	Callback func(config *common.CliCommandConfig) error
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
		}
	})
	return registryInstance
}

func commandExit(config *common.CliCommandConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *common.CliCommandConfig) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range *GetRegistry() {
		fmt.Printf("%s: %s\n", command.Name, command.Description)
	}
	return nil
}

func commandMapNext(config *common.CliCommandConfig) error {
	var url string
	if config.Next == nil {
		if config.Previous == nil {
			url = pokeapi.BaseURL + "/location-area"
		} else {
			fmt.Println("You're on the last page")
			return nil
		}
	} else {
		url = *config.Next
	}
	return pokeapi.GetLocationAreas(url, config)
}

func commandMapBack(config *common.CliCommandConfig) error {
	if config.Previous == nil {
		if config.Next == nil {
			fmt.Println("Must use 'map' command at least once before 'mapb'")
		} else {
			fmt.Println("You're on the first page")
		}
		return nil
	}
	return pokeapi.GetLocationAreas(*config.Previous, config)
}