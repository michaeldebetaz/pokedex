package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/michaeldebetaz/pokedexcli/internal/pokecache"
)

type responseData struct {
	Next     string
	Previous string
	Results  []results
}

type results struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetLocationAreasData(url string, cache pokecache.Cache) responseData {
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
	}

	var data responseData

	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	cache.Add(url, body)

	return data
}
