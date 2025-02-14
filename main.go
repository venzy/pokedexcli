package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	registerCommands()
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
		commandEntry, ok := commandRegistry[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := commandEntry.callback()
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

type cliCommand struct {
	name string
	description string
	callback func() error
}

var commandRegistry map[string]cliCommand
func registerCommands() {
	commandRegistry = map[string]cliCommand {
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
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range commandRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}