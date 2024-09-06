package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"pokedexcli/internal/pokecache"
)

var cache = pokecache.NewCache(5 * time.Minute)

type config struct {
	Next     *string
	Previous *string
}

type LocationResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

type Command struct {
	name        string
	description string
	callback    func(cfg *config, args []string) error
}

type ExploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("")
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println("")
	fmt.Println("help: Displays this help message")
	fmt.Println("map: Displays 20 location areas")
	fmt.Println("mapb: Displays the previous 20 location areas")
	fmt.Println("explore <area>: Explore a location area and list Pokemon")
	fmt.Println("exit: Exits the program")
	fmt.Println("")
	return nil
}

func commandExit(cfg *config, args []string) error {
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func fetchLocationAreas(url string) (*LocationResponse, error) {
	//Check if the response is in the cache
	if cachedData, found := cache.Get(url); found {
		fmt.Println("us9ing cached data for", url)

		var locations LocationResponse
		if err := json.Unmarshal(cachedData, &locations); err != nil {
			return nil, err
		}
		return &locations, nil
	}

	// If not in cache, make the network request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var locations LocationResponse
	if err := json.Unmarshal(body, &locations); err != nil {
		return nil, err
	}
	return &locations, nil
}

func commandMap(cfg *config, args []string) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.Next != nil {
		url = *cfg.Next
	}

	locations, err := fetchLocationAreas(url)
	if err != nil {
		return fmt.Errorf("failed to fetch locations: %w", err)
	}

	// Print the names of the location areas
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	// Update the config for pagination
	cfg.Next = locations.Next
	cfg.Previous = locations.Previous

	return nil
}

func commandMapb(cfg *config, args []string) error {
	if cfg.Previous == nil {
		return fmt.Errorf("no previous locations to display")
	}

	locations, err := fetchLocationAreas(*cfg.Previous)
	if err != nil {
		return fmt.Errorf("failed to fetch locations: %w", err)
	}

	// Print the names of the location areas
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	// Update the config for pagination
	cfg.Next = locations.Next
	cfg.Previous = locations.Previous

	return nil
}

func fetchPokemonInLocation(area string) (*ExploreResponse, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", area)

	// Check the cache first
	if cachedData, found := cache.Get(url); found {
		fmt.Println("Using cached data for", area)

		var exploreResponse ExploreResponse
		if err := json.Unmarshal(cachedData, &exploreResponse); err != nil {
			return nil, err
		}
		return &exploreResponse, nil
	}

	// If not in cache, fetch from the API
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response into the ExploreResponse struct
	var exploreResponse ExploreResponse
	if err := json.Unmarshal(body, &exploreResponse); err != nil {
		return nil, err
	}

	// Cache the response
	cache.Add(url, body)

	return &exploreResponse, nil
}

// The explore command
func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a location area to explore")
	}

	area := args[0]
	fmt.Printf("Exploring %s...\n", area)

	// Fetch Pokémon in the given area
	exploreResponse, err := fetchPokemonInLocation(area)
	if err != nil {
		return fmt.Errorf("failed to explore location: %w", err)
	}

	// Print the found Pokémon
	fmt.Println("Found Pokemon:")
	for _, encounter := range exploreResponse.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

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
		"explore": {
			name:        "explore",
			description: "Explore a location area and list Pokemon",
			callback:    commandExplore,
		},
	}

	cfg := &config{}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		commandName := input[0]
		args := input[1:]
		if command, exists := commands[commandName]; exists {
			if err := command.callback(cfg, args); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command:", input)
		}
	}
}
