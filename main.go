package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if ! scanner.Scan() {
			break
		}
		text := scanner.Text()
		tokens := cleanInput(text)
		if len(tokens) > 0 {
			fmt.Printf("Your command was: %s\n", tokens[0])
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