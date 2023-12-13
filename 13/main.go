package main

import (
	"errors"
	"fmt"
	"os"

	io "github.com/faideww/aoc-2023/lib"
)

func main() {
	input := io.ReadInputFile(os.Args[1])

	maps := parseInput(input)

	sum := 0
	for _, m := range maps {
		ord, count := findReflectPoint(m)
		switch ord {
		case "row":
			sum += count * 100
		case "col":
			sum += count
		default:
			panic(errors.New("no reflect found"))
		}
	}

	fmt.Printf("sum of reflects: %d\n", sum)
}

func parseInput(input string) []string {
	return io.TrimAndSplitBy(input, "\n\n")
}

// returns "row" or "col", and the index of the reflection
func findReflectPoint(mapString string) (string, int) {
	rows := io.TrimAndSplit(mapString)

	for i := 0; i < len(rows); i++ {
		errors := isReflected(rows, i)
		if errors == 1 {
			return "row", i + 1
		}
	}

	transposedMap := transposeMap(rows)

	for i := 0; i < len(transposedMap); i++ {
		errors := isReflected(transposedMap, i)
		if errors == 1 {
			return "col", i + 1
		}
	}

	return "", 0
}

// tests if the map is reflected around map[row] and map[row+1]
func isReflected(mapRows []string, rowIdx int) int {
	if rowIdx == len(mapRows)-1 {
		return -1
	}
	low := rowIdx
	high := rowIdx + 1

	errors := 0
	for low >= 0 && high < len(mapRows) {
		for x := 0; x < len(mapRows[0]); x++ {
			if mapRows[low][x] != mapRows[high][x] {
				errors++
			}
		}
		low--
		high++
	}

	return errors
}

func transposeMap(mapRows []string) []string {
	width := len(mapRows[0])
	height := len(mapRows)
	result := make([]string, width)

	for y := 0; y < width; y++ {
		dest := make([]byte, height)
		for x := 0; x < height; x++ {
			dest[x] = mapRows[x][y]
		}
		result[y] = string(dest)
	}

	return result
}
