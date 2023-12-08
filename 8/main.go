package main

import (
	"fmt"
	"os"
	"strings"

	io "github.com/faideww/aoc-2023/lib"
)

type Node struct {
	name  string
	left  string
	right string
}

type Network struct {
	instructions []rune
	nodes        map[string]Node
}

func main() {
	input := io.ReadInputFile(os.Args[1])
	network := parseInput(input)

	fmt.Printf("network: %+v\n", network)

	// stepCount := findExit(network)
	// fmt.Printf("found exit in %d steps\n", stepCount)
	stepCount2 := findExits(network)
	fmt.Printf("found all exits in %d steps\n", stepCount2)
}

func parseInput(input string) Network {
	lines := strings.Split(strings.ReplaceAll(strings.TrimSpace(input), "\r\n", "\n"), "\n")

	directionStr, mapsStr := lines[0], lines[2:]
	nodes := make(map[string]Node)

	for _, mapStr := range mapsStr {
		nodeName, left, right := mapStr[:3], mapStr[7:10], mapStr[12:15]
		nodes[nodeName] = Node{
			nodeName,
			left,
			right,
		}
	}

	return Network{
		[]rune(directionStr),
		nodes,
	}

}

func findExit(network Network) int {
	loc := "AAA"
	steps := 0
	pc := 0

	for loc != "ZZZ" {
		nextStep := network.instructions[pc]
		var nextLoc string
		if nextStep == 'L' {
			nextLoc = network.nodes[loc].left
		} else {
			nextLoc = network.nodes[loc].right
		}

		loc = nextLoc

		pc = (pc + 1) % len(network.instructions)
		steps++
	}

	return steps
}

func findExits(network Network) int {
	var starts []string
	for nodeName := range network.nodes {
		if nodeName[2] == 'A' {
			starts = append(starts, nodeName)
		}
	}

	stepCounts := make([]int, len(starts))

	for i, nodeName := range starts {
		loc := nodeName
		steps := 0
		pc := 0

		for loc[2] != 'Z' {
			nextStep := network.instructions[pc]
			var nextLoc string
			if nextStep == 'L' {
				nextLoc = network.nodes[loc].left
			} else {
				nextLoc = network.nodes[loc].right
			}

			loc = nextLoc
			pc = (pc + 1) % len(network.instructions)
			steps++
		}
		stepCounts[i] = steps
	}

	// find least common multiple of all steps

	var gcd func(a, b int) int
	gcd = func(a, b int) int {
		if a != 0 {
			return gcd(b%a, a)
		}
		return b
	}

	var lcm func(a, b int) int
	lcm = func(a, b int) int {
		return a * b / gcd(a, b)
	}

	fullSteps := stepCounts[0]
	for i := 1; i < len(stepCounts); i++ {
		fullSteps = lcm(fullSteps, stepCounts[i])
	}

	return fullSteps
}
