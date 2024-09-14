package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"container/list"
)

// Room represents a regular room in the ant colony.
type Room struct {
	Name     string
	X, Y     int
	Occupied bool
}

// StartRoom represents the start room (##start).
type StartRoom struct {
	Room
}

// EndRoom represents the end room (##end).
type EndRoom struct {
	Room
}

// Ant represents an ant with an ID and a path.
type Ant struct {
	ID   int
	Path []string
	Pos  int
}

// Graph represents the colony with rooms and tunnels.
type Graph struct {
	Rooms    map[string]*Room
	Tunnels  map[string][]string
	Start    *StartRoom
	End      *EndRoom
}

// Error messages
var (
	ErrInvalidFormat    = errors.New("ERROR: invalid data format")
	ErrMissingStartEnd  = errors.New("ERROR: missing ##start or ##end")
	ErrNoPath           = errors.New("ERROR: no path found between ##start and ##end")
)

// ParseInput parses the input file to build the graph.
func ParseInput(filename string) (*Graph, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)  // Read the  file contents
	rooms := make(map[string]*Room)
	tunnels := make(map[string][]string)
	var startRoom *StartRoom
	var endRoom *EndRoom
	var numAnts int
	var inStart, inEnd bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##start") && !strings.HasPrefix(line, "##end") {
			continue
		}

		if strings.Contains(line, "-") {
			// Parse tunnel (link between rooms)
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, 0, ErrInvalidFormat
			}
			tunnels[parts[0]] = append(tunnels[parts[0]], parts[1])
			tunnels[parts[1]] = append(tunnels[parts[1]], parts[0])
		} else if strings.HasPrefix(line, "##start") {
			inStart = true
			inEnd = false
		} else if strings.HasPrefix(line, "##end") {
			inEnd = true
			inStart = false
		} else if strings.Contains(line, " ") {
			// Parse room
			var name string
			var x, y int
			_, err := fmt.Sscanf(line, "%s %d %d", &name, &x, &y)
			if err != nil {
				return nil, 0, ErrInvalidFormat
			}

			// Assign to start or end room based on previous markers
			if inStart {
				startRoom = &StartRoom{Room{Name: name, X: x, Y: y, Occupied: false}}
				inStart = false
			} else if inEnd {
				endRoom = &EndRoom{Room{Name: name, X: x, Y: y, Occupied: false}}
				inEnd = false
			} else {
				rooms[name] = &Room{Name: name, X: x, Y: y, Occupied: false}
			}
		} else {
			// Parse number of ants
			_, err := fmt.Sscanf(line, "%d", &numAnts)
			if err != nil {
				return nil, 0, ErrInvalidFormat
			}
		}
	}

	if startRoom == nil || endRoom == nil {
		return nil, 0, ErrMissingStartEnd
	}

	graph := &Graph{Rooms: rooms, Tunnels: tunnels, Start: startRoom, End: endRoom}
	return graph, numAnts, nil
}

// BFS finds the shortest path from start to end using Breadth-First Search.
func (g *Graph) BFS() ([]string, error) {
	queue := list.New()
	queue.PushBack([]string{g.Start.Name})
	visited := make(map[string]bool)
	visited[g.Start.Name] = true

	for queue.Len() > 0 {
		path := queue.Remove(queue.Front()).([]string)
		room := path[len(path)-1]

		if room == g.End.Name {
			return path, nil
		}

		for _, neighbor := range g.Tunnels[room] {
			if !visited[neighbor] {
				visited[neighbor] = true
				newPath := append([]string{}, path...)
				newPath = append(newPath, neighbor)
				queue.PushBack(newPath)
			}
		}
	}

	return nil, ErrNoPath
}

// MoveAnt moves an ant along its path, returning the new room.
func MoveAnt(ant *Ant) (string, bool) {
	if ant.Pos < len(ant.Path)-1 {
		ant.Pos++
		return ant.Path[ant.Pos], true
	}
	return "", false
}

// SimulateAnts simulates the movement of ants through the colony and prints the movements.
func SimulateAnts(paths [][]string, numAnts int) {
	ants := make([]*Ant, numAnts)
	for i := 0; i <numAnts; i++ {
		ants[i] = &Ant{ID: i + 1, Path: paths[0], Pos: 0}
	}

	occupied := make(map[string]bool)
	startRoomCapacity := numAnts

	for {
		moved := false
		var moves []string

		for _, ant := range ants {
			if ant.Pos == 0 && startRoomCapacity > 0 { // Ants in the start room
				// Move ant from start to first room in the path if the room is free
				if !occupied[ant.Path[1]] {
					startRoomCapacity--
					occupied[ant.Path[1]] = true
					ant.Pos++
					moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, ant.Path[1]))
					moved = true
				}
			} else if ant.Pos > 0 { // Ants already moving along the path
				if room, canMove := MoveAnt(ant); canMove {
					// Move to the next room only if it's not occupied, or if it's the end room
					if !occupied[room] || room == "##end" {
						occupied[ant.Path[ant.Pos-1]] = false // Free the previous room
						occupied[room] = true                // Occupy the new room
						moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, room))
						moved = true
					}
				}
			}
		}

		if !moved {
			break // If no ant moved, end the simulation
		}

		fmt.Println(strings.Join(moves, " "))
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	filename := os.Args[1]

	// Parse the input file
	graph, numAnts, err := ParseInput(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Find the shortest path using BFS
	path, err := graph.BFS()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Simulate the ants' movements and print only the movement output
	SimulateAnts([][]string{path}, numAnts)
}
