package main

import (
	"bufio"
	"fmt"
	"github.com/venzy/pokedexcli/internal/commands"
	"os"
	"strings"
)

func main() {
	config := commands.NewConfig()
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
		} else if len(tokens) > 1 {
			config.Arguments = tokens[1:]
		}
		command := tokens[0]
		commandEntry, ok := (*commands.GetRegistry())[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := commandEntry.Callback(config)
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