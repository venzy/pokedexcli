package commands

import (
	"fmt"
	"github.com/venzy/pokedexcli/internal/pokeapi"
	"math/rand"
	"os"
	"sync"
)

const debug = true

type CliCommandConfig struct {
	Arguments []string
	Previous *string
	Next *string
	Caught map[string]*pokeapi.PokemonDetail
}

func NewConfig() *CliCommandConfig {
	config := CliCommandConfig{}
	config.Arguments = []string{}
	config.Caught = map[string]*pokeapi.PokemonDetail{}
	return &config
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
			"catch": {
				Name: "catch",
				Description: "Attempt to catch a given Pokémon",
				Callback: commandCatch,
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

func commandCatch(config *CliCommandConfig) error {
	if len(config.Arguments) != 1 {
		return fmt.Errorf("catch command expects 1 argument, the Pokémon name")
	}
	pokemonName := config.Arguments[0]

	if _, ok := config.Caught[pokemonName]; ok {
		fmt.Printf("%s already caught!\n", pokemonName)
		return nil
	}

	detail, err := pokeapi.GetPokemonDetail(pokemonName)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	// Generate a random 'user experience level' in the range 50 - 608
	// Minimum Pokemon exp is probably 20 but we guarantee a minimum level of 50
	const minExp = 50
	const maxExp = 608
	expToCatch := rand.Intn(maxExp - minExp + 1) + minExp
	if debug {
		fmt.Printf("DEBUG: base exp: %v, expToCatch: %v\n", detail.BaseExperience, expToCatch)
	}
	if detail.BaseExperience <= expToCatch {
		// Caught
		fmt.Printf("%s was caught!\n", pokemonName)

		config.Caught[pokemonName] = detail
	} else {
		// Escaped
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}