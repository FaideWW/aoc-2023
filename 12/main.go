package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2023/lib"
)

type Row struct {
	layout string
	groups []int
}

type MapKey struct {
	s string
	n int
}

func main() {
	f, err := os.Create("profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	input := io.ReadInputFile(os.Args[1])
	rows := parseInput(input)
	sum := 0
	for _, r := range rows {
		fits := findPossibleFits(r.layout, r.groups, make(map[MapKey]int))
		sum += fits

	}

	fmt.Printf("possible arrangements: %d\n", sum)
}

func parseInput(input string) []Row {
	lines := io.TrimAndSplit(input)
	rows := make([]Row, len(lines))

	for i, line := range lines {
		line = expandRowString(line)
		components := strings.Fields(line)
		readout, groupsStr := components[0], components[1]

		var groups []int
		for _, v := range strings.Split(groupsStr, ",") {
			n, _ := strconv.Atoi(v)
			groups = append(groups, n)
		}

		rows[i] = Row{
			readout,
			groups,
		}
	}
	return rows
}

func expandRowString(row string) string {
	rowComponents := strings.Fields(row)
	layout, groups := rowComponents[0], rowComponents[1]

	expandedLayout := make([]string, 5)
	expandedGroups := make([]string, 5)
	for i := 0; i < 5; i++ {
		expandedLayout[i] = layout
		expandedGroups[i] = groups
	}

	return strings.Join(expandedLayout, "?") + " " + strings.Join(expandedGroups, ",")
}

// The general idea: recursively try to fit groups into the remaining layout
// string. base case: if there is one group left, find all locations it can fit
// into the remaining string otherwise: take the leftmost group, identify where
// it can fit into the string, and recurse with the remaining groups and the
// substring after the group has been "set"
func findPossibleFits(layout string, groups []int, memo map[MapKey]int) int {
	groupSize := groups[0]

	ranges := findRanges(layout, groupSize)

	if len(groups) == 1 {
		validArrangements := len(ranges)
		for _, r := range ranges {
			// check the consumed string; if there are any remaining '#'s, this is an
			// invalid arrangement
			didNotSkip := true
			for _, c := range layout[r-1:] {
				if c == '#' {
					didNotSkip = false
					break
				}
			}
			if !didNotSkip {
				validArrangements--
			}
		}
		return validArrangements
	}

	sum := 0
	rest := groups[1:]
	for _, r := range ranges {
		if r <= len(layout) {
			substr := layout[r:]
			mapKey := MapKey{substr, len(rest)}

			// memoize recursion cases
			var subRanges int
			if memoizedCount, ok := memo[mapKey]; ok {
				subRanges = memoizedCount
			} else {
				subRanges = findPossibleFits(substr, rest, memo)
				memo[mapKey] = subRanges
			}

			sum += subRanges
		}
	}

	return sum
}

// rules for a range
// - must be contiguous run of '?' or '#' runes
// - must be surrounded by '?' or '.' runes, or be at the beginning or end of the string
// - must not leave any unconsumed '#' runes
func findRanges(layout string, size int) []int {
	var ranges []int

	// quick optimization - we only need to search up to the first '#' + the size of the group
	searchStr := layout
	firstHash := strings.Index(layout, "#")
	if firstHash != -1 && firstHash+size+2 < len(layout) {
		searchStr = layout[:firstHash+size+2]
	}

	for i := 0; i < len(searchStr)-size+1; i++ {
		r := searchStr[i : i+size]
		validLeft := i == 0 || searchStr[i-1] == '?' || searchStr[i-1] == '.'
		validRight := i+size+1 >= len(searchStr) || (searchStr[i+size] == '?' || searchStr[i+size] == '.')
		validRange := true
		for _, c := range r {
			if c != '?' && c != '#' {
				validRange = false
				break
			}
		}

		// ensure we didn't skip any '#' runes
		didNotSkip := true
		for _, c := range searchStr[:i] {
			if c == '#' {
				didNotSkip = false
				break
			}
		}

		// consume to the left of the range, the group itself, and 1 "buffer" rune
		// (groups must be separated by 1 '.' character
		if validLeft && validRight && validRange && didNotSkip {
			ranges = append(ranges, i+size+1)
		}
	}

	return ranges
}
