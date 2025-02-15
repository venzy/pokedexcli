package commands

import (
	"fmt"
	"github.com/venzy/pokedexcli/internal/pokeapi"
	"math"
	"math/rand"
	"os"
	"sync"
)

const debug = false

type CliCommandContext struct {
	Arguments []string
	Previous *string
	Next *string
	Caught map[string]*pokeapi.PokemonDetail
}

func NewContext() *CliCommandContext {
	context := CliCommandContext{}
	context.Arguments = []string{}
	context.Caught = map[string]*pokeapi.PokemonDetail{}
	return &context
}

type CliCommand struct {
	Name string
	Description string
	Callback func(context *CliCommandContext) error
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
				Description: "Display a help message",
				Callback: commandHelp,
			},
			"map": {
				Name: "map",
				Description: "Display map locations 20 at a time",
				Callback: commandMapNext,
			},
			"mapb": {
				Name: "mapb",
				Description: "Display previous 20 map locations",
				Callback: commandMapBack,
			},
			"explore": {
				Name: "explore",
				Description: "Return a list of all Pokémon in a given location",
				Callback: commandExplore,
			},
			"catch": {
				Name: "catch",
				Description: "Attempt to catch a given Pokémon",
				Callback: commandCatch,
			},
			"inspect": {
				Name: "inspect",
				Description: "Inspect a given Pokémon you've already caught",
				Callback: commandInspect,
			},
			"pokedex": {
				Name: "pokedex",
				Description: "List the Pokémon you've already caught",
				Callback: commandPokedex,
			},
		}
	})
	return registryInstance
}

func commandExit(context *CliCommandContext) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(context *CliCommandContext) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range *GetRegistry() {
		fmt.Printf("%s: %s\n", command.Name, command.Description)
	}
	return nil
}

func commandMapNext(context *CliCommandContext) error {
	if context.Next == nil {
		if context.Previous != nil {
			fmt.Println("You're on the last page")
			return nil
		}
	}

	data, err := pokeapi.GetLocationAreas(context.Next)
	if err != nil {
		return err
	}

	context.Previous = data.Previous
	context.Next = data.Next

	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapBack(context *CliCommandContext) error {
	if context.Previous == nil {
		if context.Next == nil {
			fmt.Println("Must use 'map' command at least once before 'mapb'")
		} else {
			fmt.Println("You're on the first page")
		}
		return nil
	}

	data, err := pokeapi.GetLocationAreas(context.Previous)
	if err != nil {
		return err
	}

	context.Previous = data.Previous
	context.Next = data.Next

	for _, result := range data.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandExplore(context *CliCommandContext) error {
	if len(context.Arguments) != 1 {
		return fmt.Errorf("explore command expects 1 argument, the area name")
	}
	areaName := context.Arguments[0]

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

func commandCatch(context *CliCommandContext) error {
	if len(context.Arguments) != 1 {
		return fmt.Errorf("catch command expects 1 argument, the Pokémon name")
	}
	pokemonName := context.Arguments[0]

	if _, ok := context.Caught[pokemonName]; ok {
		fmt.Printf("%s already caught!\n", pokemonName)
		return nil
	}

	detail, err := pokeapi.GetPokemonDetail(pokemonName)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	// This is the range we expect Pokemon base experience to fall within
	const minExp = 20
	const maxExp = 608

	// This is Boots' suggested algorithm, inverting the Pokemon's experience on a normalised scale
	// See repo history for my original naiive algorithm.
	scaledDifficulty := 1.0 - (float64(detail.BaseExperience - minExp) / float64(maxExp - minExp))

	// Minimum chance to catch of 20%
	catchChance := math.Max(0.2, scaledDifficulty)

	// Roll for Pokemon to escape
	roll := rand.Float64()

	if debug {
		fmt.Printf("DEBUG: escape roll: %v, catchChance: %v\n", roll, catchChance)
	}
	if roll <= catchChance {
		// Caught
		fmt.Printf("%s was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")

		context.Caught[pokemonName] = detail
	} else {
		// Escaped
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(context *CliCommandContext) error {
	if len(context.Arguments) != 1 {
		return fmt.Errorf("inspect command expects 1 argument, the Pokémon name")
	}
	pokemonName := context.Arguments[0]

	detail, ok := context.Caught[pokemonName]
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Printf("Name: %v\n", detail.Name)
	fmt.Printf("Height: %v\n", detail.Height)
	fmt.Printf("Weight: %v\n", detail.Weight)
	fmt.Println("Stats:")
	for _, stat := range detail.Stats {
		fmt.Printf("  - %s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typeInfo := range detail.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}

func commandPokedex(context *CliCommandContext) error {
	fmt.Println("Your Pokedex:")
	for name := range context.Caught {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}
