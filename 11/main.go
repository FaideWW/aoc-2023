package main

import (
	"fmt"
	"os"

	io "github.com/faideww/aoc-2023/lib"
)

type Position struct {
	x int
	y int
}

type Universe struct {
	dims     Position
	galaxies []Position
}

// 1 empty row/column expands by this factor
const EXPAND_FACTOR = 1000000

func main() {
	input := io.ReadInputFile(os.Args[1])

	universe := parseInput(input)
	rows, cols := expandUniverse(universe)

	distances := getDistances(universe, rows, cols)

	sum := 0
	for _, d := range distances {
		sum += d
	}

	fmt.Printf("sum of distances: %d\n", sum)
}

func parseInput(input string) Universe {
	lines := io.TrimAndSplit(input)

	var galaxies []Position

	height := len(lines)
	width := len(lines[0])

	for y, line := range lines {
		for x, c := range line {
			if c == '#' {
				galaxies = append(galaxies, Position{x, y})
			}
		}
	}

	return Universe{
		Position{width, height},
		galaxies,
	}
}

func expandUniverse(u Universe) ([]int, []int) {
	rowMap := make(map[int]bool)
	colMap := make(map[int]bool)

	for _, g := range u.galaxies {
		rowMap[g.y] = true
		colMap[g.x] = true
	}

	var emptyRows []int
	for y := 0; y < u.dims.y; y++ {
		if _, ok := rowMap[y]; !ok {
			emptyRows = append(emptyRows, y)
		}
	}

	var emptyCols []int
	for x := 0; x < u.dims.y; x++ {
		if _, ok := colMap[x]; !ok {
			emptyCols = append(emptyCols, x)
		}
	}

	return emptyRows, emptyCols
}

// distances are manhattan distances
func getDistances(u Universe, expandRows, expandCols []int) []int {
	var distances []int

	for i := 0; i < len(u.galaxies); i++ {
		g1 := u.galaxies[i]
		for j := i + 1; j < len(u.galaxies); j++ {
			g2 := u.galaxies[j]

			minX := min(g2.x, g1.x)
			maxX := max(g2.x, g1.x)
			minY := min(g2.y, g1.y)
			maxY := max(g2.y, g1.y)

			distance := (maxX - minX) + (maxY - minY)

			for _, y := range expandRows {
				if y > minY && y < maxY {
					distance += EXPAND_FACTOR - 1
				}
			}

			for _, x := range expandCols {
				if x > minX && x < maxX {
					distance += EXPAND_FACTOR - 1
				}
			}

			distances = append(distances, distance)
		}
	}

	return distances
}
