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
	Rooms    map[string]Room
	start    string
	end      string
	numAnts  int
	Links    map[string][]string
}

// Read the colony from the file
func checkCoordinates(colony *Colony) {
	for room, links := range colony.Links {
		fmt.Printf("%s: %v\n", room, links)
	}
	fmt.Println("Colony:\n", colony)
}

// Read the file and return the lines
func readFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening file: %v", err)
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading file: %v", err)
	}
	return lines
}

// Create Rooms
func createRoom(colony *Colony, room string) (*Colony, error) {
	parts := strings.Split(room, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid room format: %s", room)
	}
	name, xStr, yStr := parts[0], parts[1], parts[2]

	xCoord, err := strconv.Atoi(xStr)
	if err != nil {
		return nil, fmt.Errorf("invalid x coordinate: %s", xStr)
	}
	yCoord, err := strconv.Atoi(yStr)
	if err != nil {
		return nil, fmt.Errorf("invalid y coordinate: %s", yStr)
	}
	colony.Rooms[name] = Room{Name: name, X: xCoord, Y: yCoord}
	return colony, nil
}

// Create a tunnel
func createTunnel(colony *Colony, line string) (*Colony, error) {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid tunnel format: %s", line)
	}
	room1, room2 := parts[0], parts[1]
	_, existsRoom1 := colony.Rooms[room1]
	_, existsRoom2 := colony.Rooms[room2]

	if !existsRoom1 || !existsRoom2 {
		return nil, fmt.Errorf("linking to unknown rooms: %s", line)
	}
	if _, exists := colony.Links[room1]; !exists {
		colony.Links[room1] = []string{}
	}
	if _, exists := colony.Links[room2]; !exists {
		colony.Links[room2] = []string{}
	}
	colony.Links[room1] = append(colony.Links[room1], room2)
	colony.Links[room2] = append(colony.Links[room2], room1)
	return colony, nil
}

// Read colony from the given file
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
			numAnts, err = strconv.Atoi(line)
			if err != nil {
				return Colony{}, 0, fmt.Errorf("invalid number of ants: %v", err)
			}
			firstLine = false
			continue
		} else if strings.HasPrefix(line, "#") {
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
			_, err := createTunnel(&colony, line)
			if err != nil {
				return Colony{}, 0, err
			}
		} else {
			createRoom(&colony, line)
		}
	}
	return colony, numAnts, nil
}

// Find all paths from start to end
func findAllPaths(colony Colony, start, end string) [][]string {
	fmt.Println("start: ", start, " end: ", end)
	var allPaths [][]string
	visited := make(map[string]bool)
	var path []string

	var dfs func(current string)
	dfs = func(current string) {
		if current == end {
			pathCopy := make([]string, len(path))
			copy(pathCopy, path)
			allPaths = append(allPaths, pathCopy)
			return
		}
		visited[current] = true
		path = append(path, current)
		for _, neighbor := range colony.Links[current] {
			if !visited[neighbor] {
				dfs(neighbor)
			}
		}
		visited[current] = false
		path = path[:len(path)-1]
	}
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
	maxCollisions := -1
	for _, count := range pathCollisionCounts {
		if count > maxCollisions {
			maxCollisions = count
		}
	}
	var filteredPaths [][]string
	for i, path := range paths {
		if pathCollisionCounts[i] < maxCollisions {
			filteredPaths = append(filteredPaths, path)
		}
	}
	return filteredPaths
}

// Calculate the weight of paths
func calculateWeight(paths *[][]string) []int {
	weights := make([]int, len(*paths))
	for i, path := range *paths {
		weights[i] = len(path)
	}
	return weights
}

// Find the smallest path based on weight
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

// Determine candidate paths based on the number of ants
func candidatePaths(paths *[][]string, weight *[]int, numAnts int) [][]string {
	if len(*paths) == 0 {
		log.Fatal("No paths found")
	}
	var candidatePaths [][]string
	var i int
	for i = 0; i < len(*paths); i++ {
		j, p := smallestPath(paths, weight, 0)
		if len(candidatePaths) >= numAnts {
			return candidatePaths
		}
		if vertexCollision(&candidatePaths, (*paths)[j]) {
			candidatePaths = append(candidatePaths, p)
			*paths = append((*paths)[:j], (*paths)[j+1:]...)
			*weight = append((*weight)[:j], (*weight)[j+1:]...)
			i--
		} else {
			*paths = append((*paths)[:j], (*paths)[j+1:]...)
			*weight = append((*weight)[:j], (*weight)[j+1:]...)
			i--
		}
	}
	return candidatePaths
}

func removeStartAddEnd (paths *[][]string, end string)*[][]string{
	for i := range *paths {
		(*paths)[i] = append((*paths)[i], end)
		(*paths)[i] =  (*paths)[i][1:]
	}
	return paths
}

// Check for vertex collisions
func vertexCollision(candidatePaths *[][]string, path []string) bool {
	for _, candidatePath := range *candidatePaths {
		for i := range candidatePath {
			if i > 0 && i < len(path) && candidatePath[i] == path[i] {
				return false
			}
		}
	}
	return true
}

// Print ant movements through candidate paths
func printAntMovements(candidatePaths [][]string, numAnts int) {
	if len(candidatePaths) == 0 || numAnts == 0 {
		fmt.Println("No paths or ants available.")
		return
	}

	antsPositions := make([]int, numAnts)
	for i := 0; i < numAnts; i++ {
		if i < len(candidatePaths) {
			antsPositions[i] = 0
		}
	}

	maxPathLength := 0
	for _, path := range candidatePaths {
		if len(path) > maxPathLength {
			maxPathLength = len(path)
		}
	}

	for step := 0; step < maxPathLength; step++ {
		var stepMovements []string
		for ant := 0; ant < numAnts; ant++ {
			if ant < len(candidatePaths) {
				path := candidatePaths[ant]
				if antsPositions[ant] < len(path) {
					stepMovements = append(stepMovements, fmt.Sprintf("L%d-%s", ant+1, path[antsPositions[ant]]))
					antsPositions[ant]++
				}
			}
		}
		if len(stepMovements) > 0 {
			fmt.Println(strings.Join(stepMovements, " "))
		}
	}
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

	fmt.Println(numAnts)

	fmt.Println("Rooms:")
	for name, room := range colony.Rooms {
		fmt.Printf("%s: (%d, %d)\n", name, room.X, room.Y)
	}

	fmt.Println("Links:")
	for room, links := range colony.Links {
		fmt.Printf("%s: %v\n", room, links)
	}

	if colony.start != "" && colony.end != "" {
		paths := findAllPaths(colony, colony.start, colony.end)
		fmt.Println("Paths from start to end:")
		for _, path := range paths {
			fmt.Println(strings.Join(path, " -> "))
		}

		weight := calculateWeight(&paths)
		colony.numAnts = numAnts
		fmt.Println("Weights:", weight)

		candidatePaths := candidatePaths(&paths, &weight, numAnts)
		removeStartAddEnd(&candidatePaths,colony.end)
		fmt.Println("Candidate paths:",candidatePaths)
		for _, path := range candidatePaths {
			fmt.Println(strings.Join(path, " -> "))
		}

		printAntMovements(candidatePaths, numAnts)
	} else {
		fmt.Println("Start or end room not defined.")
	}
}
