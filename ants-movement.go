package main

import (
	"fmt"
)

// Ant Movement Simulation: Handles the movement of ants along the shortest available paths.
// Error Handling: Identifies and reports issues such as incorrect room formats, duplicate rooms, unrecognized rooms in links, and missing start or end rooms.

// Helper function to update the availability status of rooms in the usedRooms map
func roomAvailbilityFlag(path Path, usedRooms map[string]bool, mark bool) {
	for _, room := range path.rooms[1 : len(path.rooms)-1] { // Excludes the start and end rooms
		usedRooms[room] = mark
	}
}

func moveAnts(numAnts int, paths []Path) {
	// Initialize maps to manage ant paths and their positions
	antPath := make(map[int]int)
	antPosition := make(map[int]int)
	InPath := make([]int, len(paths))

	// Assign ants to paths based on the path with the lowest cost
	for antID := 1; antID <= numAnts; antID++ {
		pathIndex := 0
		minCost := InPath[0] + paths[0].roomsNum // this variable will be used to track the minimum cost
		for i := 1; i < len(paths); i++ {
			cost := paths[i].roomsNum + InPath[i]
			if minCost > cost {
				minCost = cost
				pathIndex = i
			}
		}
		antPath[antID] = pathIndex
		InPath[pathIndex]++
	}

	// Maps to track ants that are outside and inside the paths
	antsOutside := make(map[int][]int)
	for i := 0; i < len(paths); i++ {
		antsOutside[i] = make([]int, 0)
	}
	for i := 1; i <= len(antPath); i++ {
		antsOutside[antPath[i]] = append(antsOutside[antPath[i]], i)
	}
	antsInside := make(map[int][]int)
	var antMoving bool
	var output string

	// Main simulation loop
	for step := 1; ; step++ {
		antMoving = false

		// Move ants that are already on their paths
		for pathIndex := 0; pathIndex < len(paths); pathIndex++ {
			for j := 0; j < len(antsInside[pathIndex]); j++ {
				if antPosition[antsInside[pathIndex][j]] < paths[pathIndex].roomsNum-1 {
					antMoving = true
					antPosition[antsInside[pathIndex][j]]++
					output += fmt.Sprintf("L%d-%s ", antsInside[pathIndex][j], paths[pathIndex].rooms[antPosition[antsInside[pathIndex][j]]])
				}
			}
		}

		// Move ants from outside the paths to their respective paths
		for pathIndex := 0; pathIndex < len(paths); pathIndex++ {
			for len(antsOutside[pathIndex]) != 0 {
				if antPosition[antsOutside[pathIndex][0]] < paths[pathIndex].roomsNum-1 {
					antMoving = true
					antPosition[antsOutside[pathIndex][0]]++
					output += fmt.Sprintf("L%d-%s ", antsOutside[pathIndex][0], paths[pathIndex].rooms[antPosition[antsOutside[pathIndex][0]]])
					antID := antsOutside[pathIndex][0]
					antsInside[pathIndex] = append(antsInside[pathIndex], antID)
					antsOutside[pathIndex] = antsOutside[pathIndex][1:]
					break
				}
			}
		}
		// End the simulation if no ants moved in this step
		if !antMoving {
			break
		}
		// Output the results for this step
		fmt.Println(output)
		output = ""
	}
}
