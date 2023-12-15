package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	io "github.com/faideww/aoc-2023/lib"
)

type Lens struct {
	label string
	fl    int
}

type Box []Lens

func main() {
	input := io.ReadInputFile(os.Args[1])
	instructions := parseInput(input)

	sum := 0
	for _, instruct := range instructions {
		currentValue := hash(instruct)
		sum += currentValue
	}

	fmt.Printf("sum: %d\n", sum)

	// part 2
	boxes := fillBoxes(instructions)
	totalFocus := 0
	for boxId, b := range boxes {
		for lensOrder, lens := range b {
			totalFocus += (boxId + 1) * (lensOrder + 1) * lens.fl
		}
	}

	fmt.Printf("focusing power: %d\n", totalFocus)
}

func parseInput(input string) []string {
	return io.TrimAndSplitBy(input, ",")

}

func hash(str string) int {
	currentValue := 0
	for _, c := range str {
		ascii := int(c)
		currentValue = ((currentValue + ascii) * 17) % 256
	}
	return currentValue
}

func fillBoxes(instructions []string) []Box {
	boxes := make([]Box, 256)
	re := regexp.MustCompile(`([a-z]+)([=-])(\d)?`)
	for _, instruct := range instructions {
		components := re.FindStringSubmatch(instruct)
		label, op := components[1], components[2]
		boxId := hash(label)
		switch op {
		case "=":
			fl, _ := strconv.Atoi(components[3])
			lens := Lens{label, fl}
			idx := -1
			for i := 0; i < len(boxes[boxId]); i++ {
				l := boxes[boxId][i]
				if l.label == label {
					idx = i
				}
			}
			if idx != -1 {
				boxes[boxId][idx] = lens
			} else {
				boxes[boxId] = append(boxes[boxId], lens)
			}
		case "-":
			idx := -1
			for i := 0; i < len(boxes[boxId]); i++ {
				l := boxes[boxId][i]
				if l.label == label {
					idx = i
				}
			}
			if idx != -1 {
				boxes[boxId] = append(boxes[boxId][:idx], boxes[boxId][idx+1:]...)
			}
		}
	}
	return boxes
}
