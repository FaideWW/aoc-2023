package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2023/lib"
)

const MAX_RED = 12
const MAX_GREEN = 13
const MAX_BLUE = 14

func main() {
	input := io.ReadInputFile(os.Args[1])
	validGames, gamePowers := getGameData(input)

	validSum := 0
	for _, id := range validGames {
		fmt.Printf("game %d is valid\n", id)
		validSum += id
	}

	powerSum := 0
	for _, power := range gamePowers {
		powerSum += power
	}

	fmt.Printf("sum of games: %d\n", validSum)
	fmt.Printf("power of games: %d\n", powerSum)
}

func getGameData(input string) ([]int, []int) {
	lines := strings.Split(input, "\n")

	var validGames []int
	var gamePowers []int
	for _, line := range lines {
		gameId := getGameId(line)
		r, g, b := getMaxCubes(line)

		if r <= MAX_RED && g <= MAX_GREEN && b <= MAX_BLUE {
			validGames = append(validGames, gameId)
		}

		power := r * g * b
		gamePowers = append(gamePowers, power)

	}
	return validGames, gamePowers
}

func getGameId(line string) int {
	re := regexp.MustCompile("Game (?P<gameid>[0-9]+):")
	matches := re.FindStringSubmatch(line)
	matchIndex := re.SubexpIndex("gameid")
	gameIdStr := matches[matchIndex]
	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		panic(err)
	}
	return gameId
}

func getMaxCubes(line string) (int, int, int) {
	segments := strings.Split(line, ":")
	game := segments[1]

	pulls := strings.Split(game, ";")

	maxR := 0
	maxG := 0
	maxB := 0
	for _, pull := range pulls {
		cubes := strings.Split(pull, ",")
		for _, cube := range cubes {
			s := strings.Split(strings.TrimSpace(cube), " ")
			numCubes, err := strconv.Atoi(s[0])
			if err != nil {
				panic(err)
			}
			cubeColor := s[1]

			switch cubeColor {
			case "red":
				if maxR < numCubes {
					maxR = numCubes
				}
			case "green":
				if maxG < numCubes {
					maxG = numCubes
				}
			case "blue":
				if maxB < numCubes {
					maxB = numCubes
				}
			default:
				panic(errors.New(fmt.Sprintf("unrecognized color %s", cubeColor)))
			}
		}
	}

	return maxR, maxG, maxB
}
