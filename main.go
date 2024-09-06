package main

import (
	"bufio"
	"fmt"
	"os"
)

type Command struct {
	name        string
	description string
	callback    func() error
}

func commandHelp() error {
	fmt.Println("")
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println("")
	fmt.Println("help: Displays this help message")
	fmt.Println("exit: Exits the program")
	fmt.Println("")
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func commandMap() error {
	return nil
}

func commandMapb() error {
	return nil
}

func main() {
	commands := map[string]Command{
		"help": {
			name:        "help",
			description: "Displays this help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the program",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the name of 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Map back - display the previous 20 location areas",
			callback:    commandMapb,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if command, exists := commands[input]; exists {
			if err := command.callback(); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command:", input)
		}
	}
}
