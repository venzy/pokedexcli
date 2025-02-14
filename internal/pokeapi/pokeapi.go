package pokeapi

import (
	"encoding/json"
	"github.com/venzy/pokedexcli/internal/pokecache"
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

// Get a page of location area list data - if url is nil or empty the first page will be fetched
func GetLocationAreas(pageUrl *string) (*LocationAreas, error) {
	var url string
	if pageUrl == nil || *pageUrl == "" {
		url = BaseURL + "/location-area"
	} else {
		url = *pageUrl
	}

	bodyBytes, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		cache.Add(url, bodyBytes)
	}

	var data LocationAreas
	err := json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

type LocationAreaDetail struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetLocationAreaDetail(areaName string) (*LocationAreaDetail, error) {
	var url string = BaseURL + "/location-area/" + areaName
	bodyBytes, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		cache.Add(url, bodyBytes)
	}

	var data LocationAreaDetail
	err := json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
