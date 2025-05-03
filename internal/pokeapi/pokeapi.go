package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AkuPython/pokedexcli/internal/pokecache"
)

type LocationAreas struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}


var cache *pokecache.Cache

func init() {
	cache = pokecache.NewCache(5 * time.Second)
}

func MakeRequest(url, endpoint string) ([]byte, error) {
	full_url := url + endpoint
	
	if val, found := cache.Get(full_url); found {
		fmt.Println("Cache hit: ", full_url)
		return val, nil
	}

	resp, err := http.Get(full_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Code: %v, Reason: %v", resp.StatusCode, resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(full_url, data)
	return data, nil
}

func Unmarshall(data []byte, dataType interface{}) error {
	err := json.Unmarshal(data, dataType)
	return err
}

