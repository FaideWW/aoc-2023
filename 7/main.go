package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2023/lib"
)

type HandType int

const (
	FiveOfAKind HandType = iota
	FourOfAKind
	FullHouse
	ThreeOfAKind
	TwoPair
	OnePair
	HighCard
)

const CardOrder = "AKQT98765432J"

type Hand struct {
	handType HandType
	cards    string
	wager    int
}

func main() {
	input := io.ReadInputFile(os.Args[1])
	hands := parseInput(input)
	slices.SortFunc(hands, sortHands)
	slices.Reverse(hands)

	winnings := 0
	for i, hand := range hands {
		rank := i + 1
		winnings += hand.wager * rank
	}

	fmt.Printf("hands: %+v\n", hands)
	fmt.Printf("winnings: %d\n", winnings)
}

func sortHands(a, b Hand) int {
	if a.handType != b.handType {
		return int(a.handType) - int(b.handType)
	}

	for i := 0; i < len(a.cards); i++ {
		aIdx := strings.Index(CardOrder, a.cards[i:i+1])
		bIdx := strings.Index(CardOrder, b.cards[i:i+1])
		if aIdx != bIdx {
			return aIdx - bIdx
		}
	}

	return 0
}

func parseInput(input string) []Hand {
	lines := strings.Split(strings.ReplaceAll(strings.TrimSpace(input), "\r\n", "\n"), "\n")

	var hands []Hand
	for _, line := range lines {
		components := strings.Fields(line)
		cards, wagerStr := components[0], components[1]
		wager, _ := strconv.Atoi(wagerStr)
		handType := getHandType(cards)

		hand := Hand{
			handType,
			cards,
			wager,
		}

		hands = append(hands, hand)
	}

	return hands
}

func getHandType(cards string) HandType {
	seenCards := make(map[rune]int)
	jokers := 0

	for _, char := range cards {
		if char == 'J' {
			jokers++
		} else {
			if _, ok := seenCards[char]; !ok {
				seenCards[char] = 1
			} else {
				seenCards[char]++
			}
		}
	}

	// edge case: all 5 cards are jokers
	if len(seenCards) == 0 {
		return FiveOfAKind
	}

	var sortedSeenCounts []int
	for _, v := range seenCards {
		sortedSeenCounts = append(sortedSeenCounts, v)
	}

	slices.Sort(sortedSeenCounts)

	highCount := sortedSeenCounts[len(sortedSeenCounts)-1] + jokers
	if highCount == 5 {
		return FiveOfAKind
	}
	if highCount == 4 {
		return FourOfAKind
	}
	secondCount := sortedSeenCounts[len(sortedSeenCounts)-2]
	if highCount == 3 && secondCount == 2 {
		return FullHouse
	}
	if highCount == 3 && secondCount == 1 {
		return ThreeOfAKind
	}
	if highCount == 2 && secondCount == 2 {
		return TwoPair
	}
	if highCount == 2 && secondCount == 1 {
		return OnePair
	}

	return HighCard
}
