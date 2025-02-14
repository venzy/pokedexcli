package pokeapi

import (
	"encoding/json"
	"fmt"
	"github.com/venzy/pokedexcli/internal/common"
	"net/http"
)

const BaseURL = "https://pokeapi.co/api/v2"

type LocationAreas struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationAreas(url string, config *common.CliCommandConfig) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var data LocationAreas
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
