package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Define the data structures
type Room struct {
	Name string
	X, Y int
}

type Colony struct {
	Rooms map[string]Room
	start string
	end   string
	numAnts int
	Links map[string][]string
}

// Read the colony from the file (this function will be used to test if there a room in certain coordinates,
//
//	if so the input will be considered invalid)
func checkCoordinates(colony *Colony) {
	for room, links := range colony.Links {
		fmt.Printf("%s: %v\n", room, links)
	}
	fmt.Println("Colory:\n", colony)
}

// read the file and return the lines
func readFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening file: %v", err)
	}
	//closes the file after finishing the reading
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal("Error reading file: %v", err)
		}
	}
	return lines
}

// create Rooms
func createRoom(colony *Colony, room string) (*Colony, error) {
	// Room
	parts := strings.Split(room, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid room format: %s", room)
	}
	name, xStr, yStr := parts[0], parts[1], parts[2]

	xCoord, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid x coordinate: %s", xStr)
	}
	yCoord, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid y coordinate: %s", yStr)
	}
	colony.Rooms[name] = Room{Name: name, X: xCoord, Y: yCoord}
	return colony, nil
}

// Create a tunnel
func ceateTunnel(colony *Colony, line string) (*Colony, error) {
	// It's a link
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid tunnel format: %s", line)
	}
	// connecting the links to the rooms
	room1, room2 := parts[0], parts[1]
	// Check if the rooms exist
	_, existsRoom1 := colony.Rooms[room1]
	_, existsRoom2 := colony.Rooms[room2]

	if !existsRoom1 || !existsRoom2 {
		return nil, fmt.Errorf("Linking to unknown rooms: %s", line)
	}
	// Check if the links already exist
	if _, exists := colony.Links[room1]; !exists {
		colony.Links[room1] = []string{}
	}
	// Check if the links already exist
	if _, exists := colony.Links[room2]; !exists {
		colony.Links[room2] = []string{}
	}
	// connecting the links to the rooms
	colony.Links[room1] = append(colony.Links[room1], room2)
	//If line is commited we will be working on directional graph
	colony.Links[room2] = append(colony.Links[room2], room1)
	return colony, nil
}

// readColony reads the colony from the given file
func readColony(fileContent []string) (Colony, int, error) {
	colony := Colony{
		Rooms: make(map[string]Room),
		Links: make(map[string][]string),
	}

	var numAnts int
	firstLine := true

	for index, line := range fileContent {
		if firstLine {
			var err error
			// Handle the number of ants
			numAnts, err = strconv.Atoi(line)
			if err != nil {
				return Colony{}, 0, fmt.Errorf("invalid number of ants: %v", err)
			}

			firstLine = false
			continue
		} else if strings.HasPrefix(line, "#") {
			// Ignore comments
			if strings.HasPrefix(line, "##start") {
				if index != len(fileContent)-1 {
					colony.start = strings.Split(fileContent[index+1], " ")[0]
				} else {
					log.Fatal("Error Wrong file format")
				}
			} else if strings.HasPrefix(line, "##end") {
				if index != len(fileContent)-1 {
					colony.end = strings.Split(fileContent[index+1], " ")[0]
				} else {
					log.Fatal("Error Wrong file format")
				}
			} else {
				continue
			}

		} else if strings.Contains(line, "-") {
			_, err := ceateTunnel(&colony, line)
			if err != nil {
				return Colony{}, 0, err
			}
		} else {
			// Room

			createRoom(&colony, line)
		}
	}

	return colony, numAnts, nil
}

// Function to find all paths from start to end
func findAllPaths(colony Colony, start, end string) [][]string {
	fmt.Println("start: ", start, " end: ", end)
	var allPaths [][]string
	visited := make(map[string]bool)
	var path []string

	var dfs func(current string)
	dfs = func(current string) {
		if current == end {
			// Found a path to end
			// Make a copy of the path and add it to allPaths
			pathCopy := make([]string, len(path))
			copy(pathCopy, path)
			allPaths = append(allPaths, pathCopy)
			return
		}

		// Mark the current node as visited
		visited[current] = true
		path = append(path, current)

		// Explore neighbors
		for _, neighbor := range colony.Links[current] {
			if !visited[neighbor] {
				dfs(neighbor)
			}
		}

		// Backtrack: unmark the current node and remove it from the path
		visited[current] = false
		path = path[:len(path)-1]
	}

	// Start DFS from the start node
	dfs(start)

	return allPaths
}

// Count node appearances in paths
func countNodeCollisions(paths [][]string) map[string]int {
	nodeCollisionCount := make(map[string]int)

	for _, path := range paths {
		for _, node := range path {
			nodeCollisionCount[node]++
		}
	}

	return nodeCollisionCount
}

// Calculate the number of collisions a path has
func pathCollisionCount(path []string, nodeCollisionCount map[string]int) int {
	collisionCount := 0
	for _, node := range path {
		if nodeCollisionCount[node] > 1 {
			collisionCount++
		}
	}
	return collisionCount
}

// Remove paths with the most collisions
func removeMostCollidingPaths(paths [][]string) [][]string {
	nodeCollisionCount := countNodeCollisions(paths)
	pathCollisionCounts := make([]int, len(paths))

	for i, path := range paths {
		pathCollisionCounts[i] = pathCollisionCount(path, nodeCollisionCount)
	}

	// Find maximum collision count
	maxCollisions := -1
	for _, count := range pathCollisionCounts {
		if count > maxCollisions {
			maxCollisions = count
		}
	}

	// Remove paths with the maximum collision count
	var filteredPaths [][]string
	for i, path := range paths {
		if pathCollisionCounts[i] < maxCollisions {
			filteredPaths = append(filteredPaths, path)
		}
	}

	return filteredPaths
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}
	filename := os.Args[1]
	colony, numAnts, err := readColony(readFile(filename))
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	// Print the number of ants
	fmt.Println(numAnts)

	// Print the rooms
	fmt.Println("Rooms:")
	for name, room := range colony.Rooms {
		fmt.Printf("%s: (%d, %d)\n", name, room.X, room.Y)
	}

	// Print the links
	fmt.Println("Links:")
	for room, links := range colony.Links {
		fmt.Printf("%s: %v\n", room, links)
	}

	// Find and print all paths
	if colony.start != "" && colony.end != "" {
		paths := findAllPaths(colony, colony.start, colony.end)
		fmt.Println("Paths from start to end:")
		for _, path := range paths {
			fmt.Println(strings.Join(path, " -> "))
		}

		// Remove paths with the most collisions make sure this function is called only if we have more than one path
		paths = removeMostCollidingPaths(paths)
		fmt.Println("Weights:", calculateWeight(&paths))
		weight := calculateWeight(&paths)
		colony.numAnts = numAnts
		fmt.Println("candidate paths at index",candidatePaths(&paths, &weight,numAnts))
		fmt.Println("Fix the collision function",paths)
		fmt.Println("Filtered paths with least collisions:")
		for _, path := range paths {
			fmt.Println(strings.Join(path, " -> "))
		}
	} else {
		fmt.Println("Start or end room not defined.")
	}

}

func calculateWeight(paths *[][]string) []int {
	weights := make([]int, len(*paths))
	for i, path := range *paths {
		weights[i] = len(path)
	}
	return weights
}

func smallestPath(paths *[][]string, weight *[]int, index int) (int, []string) {
	if len(*paths) == 0 || len(*paths) < index {
		log.Fatal("No paths found")
	}
	candidatePath := (*paths)[index]
	
	indexSmallest := index
	for i := index; i < len(*paths); i++ {
		for j := i + 1; j < len(*weight); j++ {
			if len(candidatePath) > (*weight)[j] {
				candidatePath = (*paths)[j]
				indexSmallest = j
			}
		}

	}
	return indexSmallest, candidatePath
}
func candidatePaths(paths *[][]string, weight *[]int, numAnts int) [][]string {
	if len(*paths) == 0 {
		log.Fatal("No paths found")
	}
	var candidatePaths [][]string
	var i int
	for i = 0; i < len(*paths); i++ {
		if i == 0 {
			_, P := smallestPath(paths, weight, 0)
			candidatePaths = append(candidatePaths, P)
		} else {
			j, p := smallestPath(paths, weight, i)
			
			if len(*&candidatePaths) >=numAnts{
				return candidatePaths
			}
			if vertexCollision(&candidatePaths, (*paths)[j]) {
				fmt.Println("i: ", i, "path: ", p)
				candidatePaths = append(candidatePaths, p)
			} else {
				// Remove the path at index j from *paths
				*paths = append((*paths)[:j], (*paths)[j+1:]...)
				*weight = append((*weight)[:j], (*weight)[j+1:]...)
				// Since the slice has been modified, we need to adjust the loop variable
				i-- // Decrement i to recheck the index of the next path
			}
		}
	}
	return candidatePaths
}

func vertexCollision(candidatePaths *[][]string, path []string) bool {
	for _, candidatePath := range *candidatePaths {
		for i, _ := range candidatePath {
			if i>0&&i < len(path) && candidatePath[i] == path[i] {
				fmt.Println(path,"\nhhh\n",*candidatePaths)
				return false
			}
		}
	}
	return true
}

