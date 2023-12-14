package main

import (
	"fmt"
	"os"

	io "github.com/faideww/aoc-2023/lib"
)

type Grid struct {
	width   int
	height  int
	pillars [][]bool
	rocks   [][]bool
}

const CYCLE_COUNT = 1000000000

func main() {
	input := io.ReadInputFile(os.Args[1])
	grid := parseInput(input)

	seenRocks := make([][][]bool, 0)
	remainder := 0

	for i := 0; i < CYCLE_COUNT; i++ {
		seenRocks = append(seenRocks, copyRocks(grid.rocks))

		tiltGridNorth(grid)
		tiltGridWest(grid)
		tiltGridSouth(grid)
		tiltGridEast(grid)

		looped := false
		for loopStart, lastRocks := range seenRocks {
			if compareRocks(lastRocks, grid.rocks) {
				loopEnd := i + 1
				looped = true
				loopSize := loopEnd - loopStart
				fmt.Printf("loop detected: cycle %d = cycle %d (loop size %d)\n", loopEnd, loopStart, loopSize)
				remainder = (CYCLE_COUNT - loopEnd) % loopSize
				fmt.Printf("remaining cycles: %d\n", remainder)
				break
			}
		}

		if looped {
			break
		}
	}

	// execute the remaining cycles
	for i := 0; i < remainder; i++ {
		tiltGridNorth(grid)
		tiltGridWest(grid)
		tiltGridSouth(grid)
		tiltGridEast(grid)
	}

	load := calculateLoad(grid)
	fmt.Printf("total load: %d\n", load)
}

func copyRocks(src [][]bool) [][]bool {
	dest := make([][]bool, len(src))
	for i := range src {
		dest[i] = make([]bool, len(src[i]))
		copy(dest[i], src[i])
	}

	return dest
}

func compareRocks(a, b [][]bool) bool {
	for y := range a {
		for x := range a[y] {
			if a[y][x] != b[y][x] {
				return false
			}
		}
	}
	return true
}

func parseInput(input string) Grid {
	lines := io.TrimAndSplit(input)
	height := len(lines)
	width := len(lines[0])
	pillars := make([][]bool, height)
	for i := 0; i < height; i++ {
		pillars[i] = make([]bool, width)
	}

	rocks := make([][]bool, height)
	for i := 0; i < height; i++ {
		rocks[i] = make([]bool, width)
	}

	for y, line := range lines {
		for x, c := range line {
			switch c {
			case 'O':
				rocks[y][x] = true
			case '#':
				pillars[y][x] = true
			}
		}
	}

	return Grid{
		width,
		height,
		pillars,
		rocks,
	}
}

func printGrid(g Grid) {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if g.rocks[y][x] {
				fmt.Printf("%s", "O")
			} else if g.pillars[y][x] {
				fmt.Printf("%s", "#")
			} else {
				fmt.Printf("%s", ".")
			}
		}
		fmt.Println()
	}
}

// slide all rocks up until they can't be moved anymore
func tiltGridNorth(g Grid) {
	for y := 1; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if !g.rocks[y][x] {
				continue
			}
			// scan up until there is an obstacle
			nextY := y - 1
			for ; nextY >= 0; nextY-- {
				if g.rocks[nextY][x] || g.pillars[nextY][x] {
					break
				}
			}
			g.rocks[y][x] = false
			g.rocks[nextY+1][x] = true
		}
	}
}

func tiltGridSouth(g Grid) {
	for y := g.height - 2; y >= 0; y-- {
		for x := 0; x < g.width; x++ {
			if !g.rocks[y][x] {
				continue
			}
			// scan down until there is an obstacle
			nextY := y + 1
			for ; nextY < g.height; nextY++ {
				if g.rocks[nextY][x] || g.pillars[nextY][x] {
					break
				}
			}
			g.rocks[y][x] = false
			g.rocks[nextY-1][x] = true
		}
	}
}

func tiltGridEast(g Grid) {
	for x := g.width - 2; x >= 0; x-- {
		for y := 0; y < g.height; y++ {
			if !g.rocks[y][x] {
				continue
			}
			// scan down until there is an obstacle
			nextX := x + 1
			for ; nextX < g.width; nextX++ {
				if g.rocks[y][nextX] || g.pillars[y][nextX] {
					break
				}
			}
			g.rocks[y][x] = false
			g.rocks[y][nextX-1] = true
		}
	}
}

func tiltGridWest(g Grid) {
	for x := 1; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			if !g.rocks[y][x] {
				continue
			}
			// scan up until there is an obstacle
			nextX := x - 1
			for ; nextX >= 0; nextX-- {
				if g.rocks[y][nextX] || g.pillars[y][nextX] {
					break
				}
			}
			g.rocks[y][x] = false
			g.rocks[y][nextX+1] = true
		}
	}
}

func calculateLoad(g Grid) int {
	load := 0
	for y := 0; y < g.height; y++ {
		loadPerRock := g.height - y
		for x := 0; x < g.width; x++ {
			if g.rocks[y][x] {
				load += loadPerRock
			}
		}
	}
	return load
}
