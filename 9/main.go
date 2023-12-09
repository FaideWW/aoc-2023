package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2023/lib"
)

func main() {
	input := io.ReadInputFile(os.Args[1])
	sequences := parseInput(input)

	sum := 0
	reversedSum := 0
	for _, seq := range sequences {
		nextValue := findNextValue(seq)
		sum += nextValue
		slices.Reverse(seq)
		prevValue := findNextValue(seq)
		reversedSum += prevValue
	}

	fmt.Printf("sum of next values: %d\n", sum)
	fmt.Printf("sum of prev values: %d\n", reversedSum)
}

func parseInput(input string) [][]int {
	lines := strings.Split(strings.ReplaceAll(strings.TrimSpace(input), "\r\n", "\n"), "\n")

	sequences := make([][]int, len(lines))
	for i, line := range lines {
		vals := strings.Fields(line)
		nums := make([]int, len(vals))
		for i, val := range vals {
			nums[i], _ = strconv.Atoi(val)
		}

		sequences[i] = nums
	}

	return sequences
}

func findNextValue(seq []int) int {
	allZeroes := true
	for _, n := range seq {
		if n != 0 {
			allZeroes = false
			break
		}
	}
	if allZeroes {
		return 0
	}

	subSequence := make([]int, len(seq)-1)
	for i := 0; i < len(subSequence); i++ {
		subSequence[i] = seq[i+1] - seq[i]
	}

	nextSubseqValue := findNextValue(subSequence)

	return seq[len(seq)-1] + nextSubseqValue
}
