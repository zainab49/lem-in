package main

import (
	"fmt"
	"strconv"
	"strings"
)

func RoomHandler(line string, isStart bool, isEnd bool) error {
	// Split the line into components separated by spaces
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return fmt.Errorf("incorrect format: %s", line) 
	}
	name := parts[0] 

	if name[0] == '#' {
		return fmt.Errorf("invalid room name: %s", name)
	}

	// Convert the x coordinate from string to integer
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid x coordinate for room %s", name) // Handle conversion error
	}

	// Convert the y coordinate from string to integer
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid y coordinate for room %s", name) // Handle conversion error
	}

	// Check if the room name already exists
	if _, exists := farm.Rooms[name]; exists {
		return fmt.Errorf("duplicate room name: %s", name)
	}

	// Create a coordinate tuple and ensure it is not duplicated
	coord := [2]int{x, y}
	if coords[coord] {
		return fmt.Errorf("duplicate room coordinates: %d %d", x, y)
	}

	// Create and add a new Room instance to the farm
	room := Room{Name: name, X: x, Y: y}
	farm.Rooms[name] = room // Add room to the farm

	// Mark the coordinates as occupied
	coords[coord] = true

	// If the room is designated as the start room, set the global startRoom variable
	if isStart {
		startRoom = room
	}

	// If the room is designated as the end room, set the global endRoom variable
	if isEnd {
		endRoom = room
	}

	return nil
}

func LinkHandler(line string) error {
	// Split the line into two parts using the hyphen as a separator
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return fmt.Errorf("incorrect link format: %s", line) // Ensure the link has exactly two parts
	}
	room1, room2 := parts[0], parts[1] // Extract room names from the parts

	// Check if the first room is known in the farm
	if _, exists := farm.Rooms[room1]; !exists {
		return fmt.Errorf("unknown room in link: %s", room1) // Error if room1 is unknown
	}

	// Check if the second room is known in the farm
	if _, exists := farm.Rooms[room2]; !exists {
		return fmt.Errorf("unknown room in link: %s", room2) // Error if room2 is unknown
	}

	if linkExists(farm, room1, room2) {
		return fmt.Errorf("invalid format: multiple links between %s and %s", room1, room2)
	}

	if checker(room1, room2) {
		return fmt.Errorf("cannot link a room to itself")
	}

	// Add the link to the farm, establishing a bidirectional connection
	farm.Links[room1] = append(farm.Links[room1], room2) // Link room1 to room2
	farm.Links[room2] = append(farm.Links[room2], room1) // Link room2 to room1

	return nil
}

func linkExists(farm Farm, room1, room2 string) bool {
	// Check if room1 is linked to room2
	for _, link := range farm.Links[room1] {
		if link == room2 {
			return true
		}
	}
	// Check if room2 is linked to room1
	for _, link := range farm.Links[room2] {
		if link == room1 {
			return true
		}
	}
	return false
}

func checker(room1, room2 string) bool {
	if room1 == room2 || room2 == room1 {
		return true
	}
	return false
}
