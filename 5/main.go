package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2023/lib"
)

type Range struct {
	low  int
	high int
}

type RangeMap struct {
	destStart int
	srcStart  int
	length    int
}

func main() {
	input := io.ReadInputFile(os.Args[1])

	seeds, maps := parseInput(input)
	locations := mapSeedsToLocations(seeds, maps)

	minLoc := locations[0].low
	for _, loc := range locations {
		if minLoc > loc.low {
			minLoc = loc.low
		}
	}

	fmt.Printf("min location: %d\n", minLoc)
}

func parseInput(input string) ([]Range, [][]RangeMap) {
	var seeds []Range
	var rangeMaps [][]RangeMap

	// windows hack
	groups := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n\n")

	seedsGroup, mapGroups := groups[0], groups[1:]

	seedStrs := strings.Fields(seedsGroup)[1:]
	for i := 0; i < len(seedStrs); i += 2 {
		low, _ := strconv.Atoi(seedStrs[i])
		length, _ := strconv.Atoi(seedStrs[i+1])
		seeds = append(seeds, Range{low, low + length})
	}

	for _, mapGroup := range mapGroups {
		ranges := []RangeMap{}
		mapStrings := strings.Split(mapGroup, "\n")[1:]

		for _, mapStr := range mapStrings {
			components := strings.Fields(mapStr)
			destStart, _ := strconv.Atoi(components[0])
			srcStart, _ := strconv.Atoi(components[1])
			length, _ := strconv.Atoi(components[2])
			ranges = append(ranges, RangeMap{
				destStart,
				srcStart,
				length,
			})
		}

		rangeMaps = append(rangeMaps, ranges)
	}

	return seeds, rangeMaps
}

/*
*

	Maintain a list of ranges at the current "stage" (current)
	For each stage:
	  For each range in current:
	    Map the range to the next stage, splitting the range into multiple if needed
*/
func mapSeedsToLocations(seeds []Range, maps [][]RangeMap) []Range {

	currentList := seeds
	for _, stage := range maps {
		var nextList []Range

		for _, r := range currentList {
			nextRanges := r.mapToRanges(stage)
			nextList = append(nextList, nextRanges...)
		}

		currentList = nextList
	}

	return currentList
}

func (r *Range) mapToRanges(maps []RangeMap) []Range {
	var nextRanges []Range

	// sort maps by srcStart
	srcCmp := func(a, b RangeMap) int {
		return cmp.Compare(a.srcStart, b.srcStart)
	}
	slices.SortFunc(maps, srcCmp)

	currentRange := *r
	for _, m := range maps {
		// case 0: map does not overlap with range
		// ...

		if m.srcStart <= currentRange.low && m.srcStart+m.length > currentRange.low {
			// case 1: map starts below or at bottom of the range
			//           |---cr----|
			//         |--m--|
			//
			//           |---cr----|
			//         |------m------|
			// result: MAP cr.low - min(m.high, cr.high)
			//         LOOP remainder
			overlapTop := min(m.srcStart+m.length, currentRange.high)
			mappedRange := Range{
				currentRange.low - m.srcStart + m.destStart,
				overlapTop - m.srcStart + m.destStart,
			}

			nextRanges = append(nextRanges, mappedRange)

			currentRange = Range{overlapTop, currentRange.high}
			if currentRange.low >= currentRange.high {
				break
			}
		} else if m.srcStart <= currentRange.high && m.srcStart+m.length > currentRange.high {

			// case 2: map starts above bottom of the range
			//           |---------|
			//              |--m--|
			//
			//           |---------|
			//                  |--m--|
			// result: NOMAP cr.low - m.low
			//         MAP m.low - min(m.high, cr.high)
			//         LOOP remainder
			overlapTop := min(m.srcStart+m.length, currentRange.high)
			unmapped := Range{currentRange.low, m.srcStart}
			mapped := Range{
				m.destStart,
				overlapTop - m.srcStart + m.destStart,
			}

			nextRanges = append(nextRanges, unmapped, mapped)

			currentRange = Range{overlapTop, currentRange.high}
			if currentRange.low >= currentRange.high {
				break
			}
		}
	}

	// if there's any remaining range, add that to nextRanges
	if currentRange.low < currentRange.high {
		nextRanges = append(nextRanges, currentRange)
	}

	// fmt.Printf("done mapping. mapped ranges: %+v\n", nextRanges)

	return nextRanges
}
