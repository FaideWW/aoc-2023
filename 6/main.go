package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2023/lib"
)

type Race struct {
	time   int
	record int
}

func main() {
	input := io.ReadInputFile(os.Args[1])

	races := parseInput(input)

	fmt.Printf("races: %+v\n", races)

	product := 1
	for _, r := range races {
		product *= findTimeRange(r)
	}

	fmt.Printf("product of ranges: %d\n", product)
}

func parseInput(input string) []Race {
	var races []Race
	lines := strings.Split(strings.TrimSpace(input), "\n")

	// Truncate "Time: " from the start of the string
	timesStr := strings.TrimSpace(lines[0][6:])
	// Truncate "Distance: " from the start of the string
	recordsStr := strings.TrimSpace(lines[1][10:])

	timesStr = strings.ReplaceAll(timesStr, " ", "")
	recordsStr = strings.ReplaceAll(recordsStr, " ", "")

	time, _ := strconv.Atoi(timesStr)
	record, _ := strconv.Atoi(recordsStr)

	race := Race{
		time,
		record,
	}

	races = append(races, race)

	return races
}

func findTimeRange(race Race) int {
	minTime := race.time + 1
	maxTime := -1
	for i := 0; i < race.time; i++ {
		distance := (race.time - i) * i
		if distance > race.record {
			if i < minTime {
				minTime = i
			}
			if i > maxTime {
				maxTime = i
			}
		}
	}

	fmt.Printf("race %+v time range: %d-%d\n", race, minTime, maxTime)

	return (maxTime - minTime) + 1
}
