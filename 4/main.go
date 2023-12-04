package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2023/lib"
)

type Card struct {
	id             int
	winningNumbers []int
	myNumbers      []int
	score          int
	nextCards      []int
}

func main() {
	input := io.ReadInputFile(os.Args[1])
	cards := parseInput(input)
	cardsSeen := calculateTotalCards(cards)

	sumOfScores := 0
	for _, g := range cards {
		sumOfScores += g.score
	}

	fmt.Printf("cards: %+v\n", cards)
	fmt.Printf("sum of scores: %d\n", sumOfScores)
	fmt.Printf("cards seen: %d\n", cardsSeen)
}

func parseInput(input string) []Card {
	var games []Card

	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		var err error
		components := strings.Split(line, ":")
		gameId, err := strconv.Atoi(strings.TrimSpace(components[0][4:]))
		if err != nil {
			panic(err)
		}

		numStrs := strings.Split(components[1], "|")

		winNumStrs := strings.Fields(numStrs[0])
		myNumStrs := strings.Fields(numStrs[1])

		winningNums := make([]int, len(winNumStrs))
		myNums := make([]int, len(myNumStrs))

		winnersMap := make(map[int]bool)
		score := 0
		matchingCards := 0

		for i := 0; i < len(winNumStrs); i++ {
			winningNums[i], err = strconv.Atoi(winNumStrs[i])
			if err != nil {
				panic(err)
			}
			winnersMap[winningNums[i]] = true
		}
		for i := 0; i < len(myNumStrs); i++ {
			myNums[i], err = strconv.Atoi(myNumStrs[i])
			if err != nil {
				panic(err)
			}

			if _, ok := winnersMap[myNums[i]]; ok {
				if score == 0 {
					score = 1
				} else {
					score *= 2
				}

				matchingCards++

			}
		}

		nextCards := make([]int, matchingCards)
		for i := 0; i < matchingCards; i++ {
			nextCards[i] = gameId + i + 1
		}

		g := Card{
			gameId,
			winningNums,
			myNums,
			score,
			nextCards,
		}

		games = append(games, g)
	}

	return games
}

func calculateTotalCards(cards []Card) int {
	cardsMap := make(map[int]Card, len(cards))
	queue := make([]int, 0)

	for _, c := range cards {
		cardsMap[c.id] = c
		queue = append(queue, c.id)
	}

	totalCardsSeen := 0

	for len(queue) > 0 {
		var id int
		id, queue = queue[0], queue[1:]
		totalCardsSeen++

		for _, next := range cardsMap[id].nextCards {
			// ensure next is a valid card id
			if next <= len(cards) {
				queue = append(queue, next)
			}
		}
	}

	return totalCardsSeen
}
