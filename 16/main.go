package main

import (
	"fmt"
	"os"

	io "github.com/faideww/aoc-2023/lib"
)

type Pos struct {
	x int
	y int
}

type Beam struct {
	pos Pos
	dir int
}

const (
	North int = iota
	East
	South
	West
)

func main() {
	input := io.ReadInputFile(os.Args[1])

	grid := parseInput(input)

	maxEnergy := 0
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[0]); x++ {
			if y > 0 && y < len(grid)-1 && x > 0 && x <= len(grid[0])-1 {
				continue
			}
			pos := Pos{x, y}
			var dir int
			if x == 0 {
				dir = East
			} else if x == len(grid[0])-1 {
				dir = West
			} else if y == 0 {
				dir = South
			} else if y == len(grid)-1 {
				dir = North
			}
			tiles := energize(grid, pos, dir)
			if maxEnergy < tiles {
				fmt.Printf("new max: %+v %d (%d tiles)\n", pos, dir, tiles)
				maxEnergy = tiles
			}
		}
	}

	fmt.Printf("max energized tiles: %d\n", maxEnergy)
}

func parseInput(input string) []string {
	return io.TrimAndSplit(input)
}

func inBounds(grid []string, p Pos) bool {
	if p.x < 0 || p.x >= len(grid[0]) || p.y < 0 || p.y >= len(grid) {
		return false
	}
	return true
}

func energize(grid []string, start Pos, dir int) int {

	seenBeams := make(map[Beam]bool)
	seenTiles := make(map[Pos]bool)

	frontier := make([]Beam, 1)
	frontier[0] = Beam{start, dir}
	for len(frontier) > 0 {
		var curr Beam
		curr, frontier = frontier[0], frontier[1:]
		_, seen := seenBeams[curr]
		if seen || curr.pos.x < 0 || curr.pos.x >= len(grid[0]) || curr.pos.y < 0 || curr.pos.y >= len(grid) {
			continue
		}
		tile := grid[curr.pos.y][curr.pos.x]
		seenBeams[curr] = true
		seenTiles[curr.pos] = true

		var next Beam
		next.dir = curr.dir
		switch tile {
		case '/':
			// east/west is rotated counterclockwise
			if curr.dir == East || curr.dir == West {
				next.dir = (4 + (curr.dir - 1)) % 4
				// north/south is rotated clockwise
			} else if curr.dir == North || curr.dir == South {
				next.dir = (curr.dir + 1) % 4
			}
		case '\\':
			// east/west is rotated clockwise
			if curr.dir == East || curr.dir == West {
				next.dir = (curr.dir + 1) % 4
				// north/south is rotated counterclockwise
			} else if curr.dir == North || curr.dir == South {
				next.dir = (4 + (curr.dir - 1)) % 4
			}
		case '|':
			if curr.dir == East || curr.dir == West {
				frontier = append(frontier, Beam{Pos{curr.pos.x, curr.pos.y - 1}, North})
				next.dir = South
			}
		case '-':
			if curr.dir == North || curr.dir == South {
				frontier = append(frontier, Beam{Pos{curr.pos.x - 1, curr.pos.y}, West})
				next.dir = East
			}
		}
		next.pos.x = curr.pos.x
		next.pos.y = curr.pos.y
		switch next.dir {
		case North:
			next.pos.y--
		case East:
			next.pos.x++
		case South:
			next.pos.y++
		case West:
			next.pos.x--
		}
		frontier = append(frontier, next)
	}

	return len(seenTiles)
}
