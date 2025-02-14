package main

import (
	"bufio"
	"fmt"
	"github.com/venzy/pokedexcli/internal/commands"
	"github.com/venzy/pokedexcli/internal/common"
	"os"
	"strings"
)

func main() {
	config := common.CliCommandConfig{}
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
		commandEntry, ok := (*commands.GetRegistry())[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := commandEntry.Callback(&config)
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