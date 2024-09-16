package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Checks if a string represents a valid integer
func isNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

// Validates and parses the input lines.
func readInputFile(lines []string) error {
	check := true // Flag to ensure the number of ants is specified

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Handle the first line
		if i == 0 {
			if !isNumber(line) {
				log.Fatal("Number of ants is not correctly specified!")
			}
			numAnts, _ = strconv.Atoi(line)
			if numAnts <= 0 || numAnts > 10000 { // Ensure the number of ants is within a valid range
				return fmt.Errorf("invalid number of ants")
			}
			check = false
			continue
		}

		// Skip empty lines and comments, except special ones
		if line == "" || (strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##start") && !strings.HasPrefix(line, "##end")) {
			continue
		}

		// Ensure no numbers appear after the initial line
		if isNumber(line) && !check {
			return fmt.Errorf("invalid format, check the input file!")
		}

		// Handle the start room definition
		if line == "##start" {
			if startCounter > 0 { // Check for multiple start rooms
				return fmt.Errorf("Invalid format: more than one start room!")
			}
			i++
			line = strings.TrimSpace(lines[i]) // Clean up the line
			err := RoomHandler(line, true, false) // Parse the start room
			if err != nil { // Handle parsing errors
				return err
			}
			startCounter++ // Update start room counter
			continue // Proceed to the next line
		}

		// Handle the end room definition
		if line == "##end" {
			if endCounter > 0 { // Check for multiple end rooms
				return fmt.Errorf("Invalid format: more than one end room!")
			}
			i++ // Move to the next line for room details
			line = strings.TrimSpace(lines[i]) // Clean up the line
			err := RoomHandler(line, false, true) // Parse the end room
			if err != nil { // Handle parsing errors
				return err
			}
			endCounter++ // Update end room counter
			continue // Proceed to the next line
		}

		// Parse lines with coordinates as regular rooms
		if strings.Contains(line, " ") {
			err := RoomHandler(line, false, false) // Parse the room
			if err != nil { // Handle parsing errors
				return err
			}
			continue // Proceed to the next line
		}

		// Parse lines with hyphens as links between rooms
		if strings.Contains(line, "-") {
			err := LinkHandler(line) // Parse the link
			if err != nil { // Handle parsing errors
				return err
			}
		} else {
			return fmt.Errorf("link not found") // Unexpected format
		}
	}
	return nil // No errors found
}

