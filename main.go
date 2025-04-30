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
			description: "PokeAPI location-areas",
			callback:    commandMap,
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
		return err
	}
	data_string := string(data[:])
	fmt.Printf("location-area data\n====================\n%v====================\n", data_string)
	return nil
}
