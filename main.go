package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// struct to represent the rooms
type Room struct {
	Name string
	X    int
	Y    int
}

// struct to represent the farm
type Farm struct {
	Rooms map[string]Room     // map of rooms
	Links map[string][]string // map of links between those rooms
}

// struct to represent the path
type Path struct {
	rooms    []string
	roomsNum int
}

// variable declaration
var (
	numAnts      int
	startRoom    Room
	endRoom      Room
	startCounter int
	endCounter   int
	farm         Farm
	coords       = make(map[[2]int]bool) //coordinates
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}
	inputFile := os.Args[1]
	file, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// Initializes the farm variable with empty maps for Rooms and Links
	farm = Farm{
		Rooms: make(map[string]Room),
		Links: make(map[string][]string),
	}

	lines := strings.Split(string(file), "\n") // storing the content of the input file in an array line by line 
	if err := readInputFile(lines); err != nil {
		log.Fatal(err)
		return
	}
	if err := validateInput(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	allPaths := FindAllPaths(startRoom.Name, endRoom.Name)
	if len(allPaths) == 0 {
		fmt.Println("Error: No valid paths found")
		return
	}
	

	filteredPaths := filteredPaths(allPaths)
	line := strings.TrimSpace(string(file))
	fmt.Println(line)
	fmt.Println()
	fmt.Println(allPaths)
	fmt.Println(filteredPaths)

	moveAnts(numAnts, filteredPaths)
	fmt.Println(numAnts)
}
