package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/michaeldebetaz/pokedexcli/internal/pokecache"
)

type pagination struct {
	Next     string
	Previous string
	Results  []paginationResults
}

type paginationResults struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetLocationAreasPagination(url string, cache pokecache.Cache) pagination {
	body := getCacheOrResponseBody(url, cache)
	var data pagination
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}
	return data
}

type locationAreaData struct {
	Id                int                `json:"id"`
	Name              string             `json:"name"`
	PokemonEncounters []pokemonEncounter `json:"pokemon_encounters"`
}

type Pokemon struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

type pokemonEncounter struct {
	Pokemon Pokemon
}

func GetLocationAreaData(url string, cache pokecache.Cache) locationAreaData {
	body := getCacheOrResponseBody(url, cache)
	var data locationAreaData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}
	return data
}

func GetPokemonData(url string, cache pokecache.Cache) Pokemon {
	body := getCacheOrResponseBody(url, cache)
	var data Pokemon
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}
	return data
}

func getCacheOrResponseBody(url string, cache pokecache.Cache) []byte {
	var body []byte

	body, ok := cache.Get(url)

	if !ok {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code %d and\nbody: %s\n", res.StatusCode, body)
		}

		cache.Add(url, body)
	}

	return body
}
