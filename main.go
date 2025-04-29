package main

import (
	"bufio"
	// "errors"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}

var supportedCommands map[string]cliCommand

func main()  {
	supportedCommands = map[string]cliCommand{
		"help": {
		name:		"help",
		description:	"Displays a help message",
		callback:	commandHelp,
		},
		"exit": {
		name:		"exit",
		description:	"Exit the Pokedex",
		callback:	commandExit,
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

