package main

import "fmt"

// Validates the parsed information
func validateInput() error {
	if startCounter != 1 {
		return fmt.Errorf("start room not correctly identified")
	}
	if endCounter != 1 {
		return fmt.Errorf("end room not correctly identified")
	}
	if len(farm.Rooms) == 0 {
		return fmt.Errorf("no rooms found")
	}
	if len(farm.Links) == 0 {
		return fmt.Errorf("no links found")
	}
	if numAnts < 0 {
		return fmt.Errorf("invalid number of ants")
	}
	return nil
}
