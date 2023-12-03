package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2023/lib"
)

func main() {
	input := io.ReadInputFile(os.Args[1])
	grid := strings.Split(strings.TrimSpace(input), "\n")
	numbers := findNumbers(grid)
	symbols := findSymbols(grid)

	partNumbers, gearRatios := findPartNumbersAndGears(grid, numbers, symbols)

	partsSum := 0
	for _, n := range partNumbers {
		partsSum += n
	}

	gearsSum := 0
	for _, n := range gearRatios {
		gearsSum += n
	}

	fmt.Printf("parts sum: %d\n", partsSum)
	fmt.Printf("gears sum: %d\n", gearsSum)
}

func findNumbers(lines []string) [][][]int {
	var nums [][][]int
	re := regexp.MustCompile("[0-9]+")
	for row, line := range lines {
		nums = append(nums, [][]int{})
		matches := re.FindAllStringIndex(line, -1)
		for _, matches := range matches {
			nums[row] = append(nums[row], []int{matches[0], matches[1]})
		}
	}

	return nums
}

func findSymbols(lines []string) [][]int {
	var symbols [][]int
	re := regexp.MustCompile("[^0-9.]+")
	for row, line := range lines {
		symbols = append(symbols, []int{})
		matches := re.FindAllStringIndex(line, -1)
		for _, matches := range matches {
			symbols[row] = append(symbols[row], matches[0])
		}
	}

	return symbols
}

func findPartNumbersAndGears(lines []string, numbers [][][]int, symbols [][]int) ([]int, []int) {
	var partNumbers []int
	var gears []int
	symbolParts := make(map[string][]int)

	for y, row := range numbers {
		for _, num := range row {
			xStart := num[0]
			xEnd := num[1]
			numStr := lines[y][xStart:xEnd]
			partNumber, err := strconv.Atoi(numStr)
			if err != nil {
				panic(err)
			}

			hasSymbol := false
			for x := xStart; x < xEnd; x++ {
				symX, symY := hasNeighboringSymbols(symbols, x, y)
				if symX != -1 && symY != -1 {
					symKey := posToStr(symX, symY)
					hasSymbol = true
					if _, ok := symbolParts[symKey]; !ok {
						symbolParts[symKey] = []int{}
					}
					symbolParts[symKey] = append(symbolParts[symKey], partNumber)
					break
				}
			}

			if hasSymbol {
				partNumbers = append(partNumbers, partNumber)
			}
		}
	}

	for _, potentialGear := range symbolParts {
		if len(potentialGear) == 2 {
			gears = append(gears, potentialGear[0]*potentialGear[1])
		}
	}

	return partNumbers, gears
}

func posToStr(x, y int) string {
	return fmt.Sprintf("%d:%d", x, y)
}

func hasNeighboringSymbols(symbols [][]int, targetX, targetY int) (int, int) {
	// search the 3x3 grid around the target for a symbol
	for y := targetY - 1; y <= targetY+1; y++ {
		if y < 0 || y >= len(symbols) {
			continue
		}
		for x := targetX - 1; x <= targetX+1; x++ {
			for _, symX := range symbols[y] {
				if symX == x {
					return x, y
				}
			}
		}
	}
	return -1, -1
}
