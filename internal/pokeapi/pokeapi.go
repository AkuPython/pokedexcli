package pokeapi

import (
	"fmt"
	"io"
	"encoding/json"
	"net/http"
)

type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func MakeRequest(url, endpoint string) ([]byte, error) {
	full_url := url + endpoint
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
	return data, nil
}

func Unmarshall(data []byte, dataType interface{}) error {
	err := json.Unmarshal(data, dataType)
	return err
}

