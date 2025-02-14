package pokeapi

import (
	"encoding/json"
	"fmt"
	"github.com/venzy/pokedexcli/internal/pokecache"
	"github.com/venzy/pokedexcli/internal/common"
	"io"
	"net/http"
	"time"
)

const BaseURL = "https://pokeapi.co/api/v2"

var cache *pokecache.Cache = pokecache.NewCache(5 * time.Second)

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
	bodyBytes, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cache.Add(url, bodyBytes)
	}

	var data LocationAreas
	err := json.Unmarshal(bodyBytes, &data)
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
