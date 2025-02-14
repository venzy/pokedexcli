package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

const baseURL = "https://pokeapi.co/api/v2"

func main() {
	config := cliCommandConfig{}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if ! scanner.Scan() {
			break
		}
		text := scanner.Text()
		tokens := cleanInput(text)
		if len(tokens) == 0 {
			continue
		}
		command := tokens[0]
		commandEntry, ok := (*getCommandRegistry())[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := commandEntry.callback(&config)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(text string) []string {
	if len(text) < 1 {
		return []string{}
	}
	cleaned := strings.Fields(text)
	for i, word := range cleaned {
		cleaned[i] = strings.ToLower(word)
	}
	return cleaned
}

type cliCommandConfig struct {
	Previous *string
	Next *string
}

type cliCommand struct {
	name string
	description string
	callback func(config *cliCommandConfig) error
}

type commandRegistry map[string]cliCommand
var commandRegistryInstance *commandRegistry
var commandRegistryOnce sync.Once

func getCommandRegistry() *commandRegistry {
	commandRegistryOnce.Do(func() {
		commandRegistryInstance = &commandRegistry{
			"exit": {
				name: "exit",
				description: "Exit the Pokedex",
				callback: commandExit,
			},
			"help": {
				name: "help",
				description: "Displays a help message",
				callback: commandHelp,
			},
			"map": {
				name: "map",
				description: "Displays map locations 20 at a time",
				callback: commandMapNext,
			},
			"mapb": {
				name: "mapb",
				description: "Displays previous 20 map locations",
				callback: commandMapBack,
			},
		}
	})
	return commandRegistryInstance
}

func commandExit(config *cliCommandConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *cliCommandConfig) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range *getCommandRegistry() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

type locationAreas struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMapNext(config *cliCommandConfig) error {
	var url string
	if config.Next == nil {
		if config.Previous == nil {
			url = baseURL + "/location-area"
		} else {
			fmt.Println("You're on the last page")
			return nil
		}
	} else {
		url = *config.Next
	}
	return getMaps(url, config)
}

func commandMapBack(config *cliCommandConfig) error {
	if config.Previous == nil {
		if config.Next == nil {
			fmt.Println("Must use 'map' command at least once before 'mapb'")
		} else {
			fmt.Println("You're on the first page")
		}
		return nil
	}
	return getMaps(*config.Previous, config)
}

func getMaps(url string, config *cliCommandConfig) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var data locationAreas
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&data)
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