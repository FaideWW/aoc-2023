package main

import (
	"errors"
	"fmt"
	"os"

	io "github.com/faideww/aoc-2023/lib"
)

type Position struct {
	x int
	y int
}

type Grid struct {
	dims  Position
	start Position
	tiles [][]rune
}

func main() {
	input := io.ReadInputFile(os.Args[1])
	grid := parseInput(input)
	fmt.Printf("grid: %+v\n", grid)

	dist, mainLoopTiles := findFurthestLoop(grid)
	fmt.Printf("furthest: %d\n", dist)

	tiles := findEnclosedTiles(grid, mainLoopTiles)
	fmt.Printf("enclosedTiles: %d\n", tiles)
}

func parseInput(input string) Grid {
	lines := io.TrimAndSplit(input)

	height, width := len(lines), len(lines[0])
	startX, startY := -1, -1
	tiles := make([][]rune, height)

	for y, line := range lines {
		if width == -1 {
			width = len(line)
		}
		tiles[y] = make([]rune, len(line))
		for x, c := range line {
			tiles[y][x] = c
			if c == 'S' {
				startX, startY = x, y
			}
		}
	}

	return Grid{
		Position{width, height},
		Position{startX, startY},
		tiles,
	}
}

func posToStr(p Position) string { return fmt.Sprintf("%d-%d", p.x, p.y) }

// The general idea: BFS, pruning nodes that don't connect to the main loop.
// Track the furthest from the start
func findFurthestLoop(grid Grid) (int, map[string]bool) {
	seen := make(map[string]bool)

	type SearchNode struct {
		pos  Position
		dist int
	}
	frontier := make([]SearchNode, 0)
	frontier = append(frontier, SearchNode{grid.start, 0})

	furthestDistance := 0

	var current SearchNode
	for len(frontier) > 0 {
		current, frontier = frontier[0], frontier[1:]
		seen[posToStr(current.pos)] = true
		if current.dist > furthestDistance {
			furthestDistance = current.dist
		}

		neighbors := getValidNeighbors(current.pos, grid)

		for _, n := range neighbors {
			if _, ok := seen[posToStr(n.pos)]; !ok {
				frontier = append(frontier, SearchNode{n.pos, current.dist + 1})
			}
		}
	}

	return furthestDistance, seen
}

type Neighbor struct {
	pos       Position
	direction string
}

func getValidNeighbors(pos Position, grid Grid) []Neighbor {
	var neighbors []Neighbor
	tile := grid.tiles[pos.y][pos.x]

	north := Neighbor{Position{pos.x, pos.y - 1}, "north"}
	east := Neighbor{Position{pos.x + 1, pos.y}, "east"}
	south := Neighbor{Position{pos.x, pos.y + 1}, "south"}
	west := Neighbor{Position{pos.x - 1, pos.y}, "west"}

	var potentialNeighbors []Neighbor
	switch tile {
	case '|':
		potentialNeighbors = append(potentialNeighbors, north, south)
	case '-':
		potentialNeighbors = append(potentialNeighbors, east, west)
	case 'L':
		potentialNeighbors = append(potentialNeighbors, north, east)
	case 'J':
		potentialNeighbors = append(potentialNeighbors, north, west)
	case '7':
		potentialNeighbors = append(potentialNeighbors, south, west)
	case 'F':
		potentialNeighbors = append(potentialNeighbors, east, south)
	case 'S':
		potentialNeighbors = append(potentialNeighbors, north, east, south, west)
	}

	for _, candidate := range potentialNeighbors {
		if candidate.pos.x < 0 || candidate.pos.x >= grid.dims.x || candidate.pos.y < 0 || candidate.pos.y >= grid.dims.y {
			continue
		}

		candidateTile := grid.tiles[candidate.pos.y][candidate.pos.x]

		switch candidateTile {
		case '|':
			if candidate.direction == "north" || candidate.direction == "south" {
				neighbors = append(neighbors, candidate)
			}
		case '-':
			if candidate.direction == "east" || candidate.direction == "west" {
				neighbors = append(neighbors, candidate)
			}
		case 'L':
			if candidate.direction == "south" || candidate.direction == "west" {
				neighbors = append(neighbors, candidate)
			}
		case 'J':
			if candidate.direction == "south" || candidate.direction == "east" {
				neighbors = append(neighbors, candidate)
			}
		case '7':
			if candidate.direction == "north" || candidate.direction == "east" {
				neighbors = append(neighbors, candidate)
			}
		case 'F':
			if candidate.direction == "north" || candidate.direction == "west" {
				neighbors = append(neighbors, candidate)
			}
		case 'S':
			neighbors = append(neighbors, candidate)
		}

	}

	return neighbors
}

func determineStartTile(grid Grid) (rune, error) {
	neighbors := getValidNeighbors(grid.start, grid)
	fmt.Printf("start neighbors: %+v\n", neighbors)
	// assume: start tile should have exactly two valid neighbors
	if neighbors[0].direction == "north" {
		if neighbors[1].direction == "south" {
			return '|', nil
		}
		if neighbors[1].direction == "east" {
			return 'L', nil
		}
		if neighbors[1].direction == "west" {
			return 'J', nil
		}
	} else if neighbors[0].direction == "east" {
		if neighbors[1].direction == "south" {
			return 'F', nil
		}
		if neighbors[1].direction == "west" {
			return '-', nil
		}
	} else if neighbors[0].direction == "south" {
		if neighbors[1].direction == "west" {
			return '7', nil
		}
	}

	return '.', errors.New("could not determine start tile")
}

// The general idea: scan left-to-right top-to-bottom, and keep track of an
// "inside" value. When we cross a pipe boundary, swap between outside and
// inside.
// Tangent pipes should be considered a single "cross", so when we
// encounter a corner pipe, we scan forward until we find the next corner.
// If the two corners are the same winding (clockwise vs counterclockwise),
// we did not cross. If they form an S, we did cross.
func findEnclosedTiles(grid Grid, mainLoop map[string]bool) int {
	startTile, err := determineStartTile(grid)
	fmt.Printf("start tile: %c\n", startTile)
	if err != nil {
		panic(err)
	}

	// convert the start tile into a normal tile, so we know how to handle it
	grid.tiles[grid.start.y][grid.start.x] = startTile

	enclosedTiles := 0

	inside := false
	for y := 0; y < grid.dims.y; y++ {
		for x := 0; x < grid.dims.x; x++ {
			tile := grid.tiles[y][x]
			_, isMainLoop := mainLoop[posToStr(Position{x, y})]
			if !isMainLoop && inside {
				enclosedTiles++
			} else if isMainLoop {
				if tile == '|' {
					inside = !inside
				} else if tile == 'F' {
					for ; x < grid.dims.x; x++ {
						nextTile := grid.tiles[y][x]
						if nextTile == '7' {
							break
						} else if nextTile == 'J' {
							inside = !inside
							break
						}
					}
				} else if tile == 'L' {
					for ; x < grid.dims.x; x++ {
						nextTile := grid.tiles[y][x]
						if nextTile == '7' {
							inside = !inside
							break
						} else if nextTile == 'J' {
							break
						}
					}
				}
			}
		}
	}

	return enclosedTiles
}
