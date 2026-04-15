# lem-in

`lem-in` is a Go command-line program that simulates moving ants through a graph of rooms and links.

It:
- parses a farm map from an input file,
- validates rooms, coordinates, and links,
- finds all paths from `##start` to `##end`,
- chooses non-overlapping paths (excluding start/end),
- prints ant movements step by step in `L<ant>-<room>` format.

## Requirements

- Go `1.22.2` (or compatible Go 1.22.x)

## Run

From the project root:

```bash
go run . <input_file>
```

Example:

```bash
go run . ex1.txt
```

## Input Format

The input file must follow this structure:

1. First line: number of ants (integer `1..10000`)
2. Room declarations: `<name> <x> <y>`
3. Special markers:
   - `##start` followed by one room declaration
   - `##end` followed by one room declaration
4. Links: `<room1>-<room2>`
5. Comments start with `#` (except `##start` / `##end`)

### Rules Enforced

- Exactly one start room and one end room
- Unique room names
- Unique room coordinates
- Links must connect existing rooms
- No duplicate links
- No self-links (`a-a`)

## Output

The program prints:
- the original input,
- a blank line,
- movement lines where each token is `L<ant_id>-<room_name>`.

Example movement lines:

```text
L1-t L2-h L3-0 
L1-E L2-A L3-o L4-t L5-h L6-0 
...
```

## Project Files

- `main.go`: program entry point and orchestration
- `readInput.go`: input parsing
- `inputValidation.go`: global validation checks
- `roomConstruction.go`: room/link creation and validation helpers
- `FindPaths.go`: DFS path discovery and path filtering
- `ants-movement.go`: ant assignment and movement simulation
- `ex1.txt`: sample input
