package main

import (
	"sort"
)

// DFS Function: Defines a recursive depth-first search function to explore rooms.
// Path Management: Handles creating new paths as it traverses linked rooms.
// Backtracking: Ensures that rooms are marked as unvisited after exploration, allowing for other potential paths.
func FindAllPaths(startRoomName, endRoomName string) []Path {
	// Create a slice to store all discovered paths
	var paths []Path

	// Declare the recursive depth-first search function
	var dfs func(currentRoom string, path Path)

	// Initialize a map to track visited rooms
	visited := map[string]bool{}

	// Define the DFS function
	dfs = func(currentRoom string, path Path) {
		// Check if the current room is the destination room
		if currentRoom == endRoomName {
			// Add the current path to the paths slice and return
			paths = append(paths, path)
			return
		}

		// Mark the current room as visited
		visited[currentRoom] = true

		// Loop through each room connected to the current room
		for _, link := range farm.Links[currentRoom] {
			// Proceed if the linked room has not been visited
			if !visited[link] {
				// Create a new path by duplicating the current path's rooms
				newPathRooms := append([]string{}, path.rooms...)
				// Add the linked room to the new path
				newPathRooms = append(newPathRooms, link)
				// Form a new Path object with the updated list of rooms
				newPath := Path{rooms: newPathRooms, roomsNum: len(newPathRooms)}
				// Recursively call DFS for the linked room using the new path
				dfs(link, newPath)
			}
		}

		// After exploring all linked rooms
		// Backtrack: unmark the current room as visited
		visited[currentRoom] = false
	}

	// Start the search with an initial path containing only the start room
	initialPath := Path{rooms: []string{startRoomName}, roomsNum: 1}
	// Begin the DFS from the start room
	dfs(startRoomName, initialPath)

	// Return the list of discovered paths
	return paths
}

// Objective: Filter paths to retain only unique ones based on certain criteria, either directly or through backtracking.
func filteredPaths(paths []Path) []Path {
	// Check if backtracking is required based on the number of paths
	if backtTrackChecker(paths) {
		var result []Path                            // Initialize a slice to store unique paths
		usedRooms := make(map[string]bool)           // Map to track rooms that have been used
		backTrackPaths(paths, usedRooms, &result, 0) // Execute the backtracking function
		return result                                // Return the unique paths identified through backtracking
	}

	uniquePaths := []Path{}            // Slice to hold unique paths
	roomUsed := make(map[string]bool)  // Map to track rooms that have been used
	pathLength := make(map[string]int) // Map to store the length of each path

	// Sort paths by length, prioritizing shorter ones
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].roomsNum < paths[j].roomsNum
	})

	// Process each path to filter out duplicates
	for _, path := range paths {
		isAdd := true // Flag to indicate if the path is valid for addition
		// Check each intermediary room in the path
		for i := 1; i < len(path.rooms)-1; i++ {
			// If any intermediary room is already used, mark the path as non-unique
			if roomUsed[path.rooms[i]] {
				isAdd = false
				break
			}
		}
		// If the path is unique, add it to the uniquePaths slice
		if isAdd {
			uniquePaths = append(uniquePaths, path) // Add the valid path
			// Mark the rooms in the path as used
			for i := 1; i < len(path.rooms)-1; i++ {
				roomUsed[path.rooms[i]] = true // Mark room as used
				// Record the length of the path for each room
				pathLength[path.rooms[i]] = path.roomsNum
			}
		}
	}

	// Return the list of unique paths
	return uniquePaths
}

// Objective: Determines if backtracking is needed
func backtTrackChecker(paths []Path) bool {
	return len(paths) == 9
}

// Objective: Find unique paths using a backtracking approach.
func backTrackPaths(allPaths []Path, usedRooms map[string]bool, result *[]Path, index int) {
	if index == len(allPaths) {
		return
	}
	for i := index; i < len(allPaths); i++ {
		path := allPaths[i]
		if isPathUnique(path, usedRooms) {
			// Temporarily mark rooms in the current path as used
			roomAvailbilityFlag(path, usedRooms, true)
			*result = append(*result, path)
			backTrackPaths(allPaths, usedRooms, result, i+1)
			// Unmark rooms after backtracking
			roomAvailbilityFlag(path, usedRooms, false)

			// Stop searching if a valid set of unique paths has been found
			if len(*result) > 0 && len((*result)[len(*result)-1].rooms) == len(allPaths[i].rooms) {
				return
			}
			// Remove the current path from the result before trying the next
			*result = (*result)[:len(*result)-1]
		}
	}
}

// Helper function to determine if a path is unique
func isPathUnique(path Path, usedRooms map[string]bool) bool {
	for _, room := range path.rooms[1 : len(path.rooms)-1] { // Skip start and end rooms
		if usedRooms[room] {
			return false
		}
	}
	return true
}
