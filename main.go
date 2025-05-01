package main

import (
	"bufio"
	// "errors"
	"fmt"
	"os"
	"strings"
	"github.com/AkuPython/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var url string = "https://pokeapi.co/api/v2/"
var mapOffset int = 0


var supportedCommands map[string]cliCommand

func main() {
	supportedCommands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "PokeAPI location-areas, next page",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "PokeAPI location-areas, previous page",
			callback:    commandMapBack,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		inp := scanner.Scan()
		if !inp {
			fmt.Printf("Scan Fail: %v\n", scanner.Err())
			break
		}
		words := cleanInput(scanner.Text())
		if len(words) < 1 {
			fmt.Printf("***No command provided***\n")
			continue
		}
		firstWord := words[0]

		if c, valid := supportedCommands[firstWord]; valid {
			// fmt.Printf("command: %v -- %v\n", firstWord, c.description)
			c.callback()
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	lower_text := strings.ToLower(text)
	words := strings.Fields(lower_text)
	return words
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, c := range supportedCommands {
		fmt.Printf("%v: %v\n", c.name, c.description)
	}
	return nil
}

func commandMap() error {
	endpoint := fmt.Sprintf("location-area/?offset=%v", mapOffset)
	mapOffset += 20
	data, err := pokeapi.MakeRequest(url, endpoint)
	if err != nil {
		fmt.Printf("ERROR making request to %v%v\n", url, endpoint)
		return err
	}
	var location_json pokeapi.LocationArea
	err = pokeapi.Unmarshall(data, &location_json)
	if err != nil {
		fmt.Printf("ERROR unmarshalling: %v\n", err)
		return err
	}
	// data_string := string(data[:])
	fmt.Printf("location-area data - pg %v\n====================\n", mapOffset/20)
	for _, v := range location_json.Results {
		fmt.Printf("%v\n", v.Name)
	}
	fmt.Println("====================")
	return nil
}

func commandMapBack() error {
	if mapOffset >= 40 {
		mapOffset -= 40
	} else {
		fmt.Println("you're on the first page")
		return fmt.Errorf("you're on the first page\n")
	}
	err := commandMap()
	return err
}
