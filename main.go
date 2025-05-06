package main

import (
	"bufio"
	"math/rand"

	// "errors"
	"fmt"
	"os"
	"strings"

	"github.com/AkuPython/pokedexcli/internal/pokeapi"
)


type cliCommand struct {
	name        string
	description string
	callback    func(*string) error
}

var url			string	= "https://pokeapi.co/api/v2/"
var mapOffset	int		= 0


var supportedCommands	map[string]cliCommand
var pokedex				map[string]pokeapi.Pokemon

func main() {
	pokedex = make(map[string]pokeapi.Pokemon)
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
		"explore": {
			name:        "explore <area>",
			description: "PokeAPI location-areas <area>, returns Pokemon in <area>",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <Pokemon>",
			description: "PokeAPI Pokemon <Pokemon>, returns Pokemon details",
			callback:    commandCatch,
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
			callbackWord := ""
			if len(words) > 1 {
				callbackWord = words[1]
			}
			err := c.callback(&callbackWord)
			if err != nil {
				fmt.Printf("%v", err)
			}
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

func commandExit(_ *string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, c := range supportedCommands {
		tabs := "\t"
		if len(c.name) < 5 {
			tabs += "\t"
		}
		if len(c.name) < 10 {
			tabs += "\t"
		}

		fmt.Printf("%v: %v%v\n", c.name, tabs, c.description)
	}
	return nil
}

func commandMap(_ *string) error {
	endpoint := fmt.Sprintf("location-area/?offset=%v", mapOffset)
	mapOffset += 20
	data, err := pokeapi.MakeRequest(url, endpoint)
	if err != nil {
		fmt.Printf("ERROR making request to %v%v\n", url, endpoint)
		return err
	}
	var location_json pokeapi.LocationAreas
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

func commandMapBack(_ *string) error {
	if mapOffset >= 40 {
		mapOffset -= 40
	} else {
		return fmt.Errorf("you're on the first page\n")
	}
	workaround := ""
	err := commandMap(&workaround)
	return err
}

func commandExplore(area *string) error {
	if *area == "" {
		return fmt.Errorf("area not provided!\n")
	}
	endpoint := fmt.Sprintf("location-area/%v/", *area)
	data, err := pokeapi.MakeRequest(url, endpoint)
	if err != nil {
		return fmt.Errorf("ERROR making request to %v%v\n", url, endpoint)
	}
	var location_json pokeapi.LocationArea
	err = pokeapi.Unmarshall(data, &location_json)
	if err != nil {
		fmt.Printf("ERROR unmarshalling: %v\n", err)
		return err
	}
	fmt.Printf("Pokemon in location-area%v\n====================\n", *area)
	for _, v := range location_json.PokemonEncounters {
		fmt.Printf("%v\n", v.Pokemon.Name)
	}
	fmt.Println("====================")
	// data_string := string(data[:])
	// fmt.Printf("location-area data\n====================\n%v\n", string(data))
	return nil
}

func commandCatch(pokemon *string) error {
	if *pokemon == "" {
		return fmt.Errorf("Pokemon not provided!\n")
	}
	endpoint := fmt.Sprintf("Pokemon/%v/", *pokemon)
	data, err := pokeapi.MakeRequest(url, endpoint)
	if err != nil {
		return fmt.Errorf("ERROR making request to %v%v\n", url, endpoint)
	}
	var pokemon_json pokeapi.Pokemon
	err = pokeapi.Unmarshall(data, &pokemon_json)
	if err != nil {
		fmt.Printf("ERROR unmarshalling: %v\n", err)
		return err
	}
	// fmt.Printf("Pokemon %v\n====================\n", *area)
	// for _, v := range location_json.PokemonEncounters {
	// 	fmt.Printf("%v\n", v.Pokemon.Name)
	// }
	fmt.Printf("Throwing a Pokeball at %v...\n", *pokemon)
	experience := pokemon_json.BaseExperience
	randInt := rand.Intn(experience)
	fmt.Printf("Values: %v, %v\n", experience, randInt)
	if randInt > 40 {
		fmt.Printf("%v was caught!\n", *pokemon)
		pokedex[*pokemon] = pokemon_json
	} else {
		fmt.Printf("%v escaped!\n", *pokemon)
	}
	fmt.Println("====================")
	// data_string := string(data[:])
	// fmt.Printf("location-area data\n====================\n%v\n", string(data))
	return nil
}
