package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2023/lib"
)

func main() {
	input := io.ReadInputFile(os.Args[1])
	sum := getValuesSum(input)
	fmt.Println("sum of values", sum)
}

func getValuesSum(input string) int {
	lines := strings.Split(input, "\n")

	sum := 0

	for _, line := range lines {
		nums := findDigits(line)

		if len(nums) < 1 {
			err := errors.New(fmt.Sprintf("did not find any numbers in %s\n", line))
			panic(err)
		}
		v, err := strconv.Atoi(fmt.Sprintf("%s%s", nums[0], nums[len(nums)-1]))
		if err != nil {
			panic(err)
		}
		sum += v
		fmt.Printf("%s: %d\n", line, v)
	}

	return sum
}

func findDigits(in string) []string {
	var digits []string

	for cursor := 0; cursor < len(in); cursor++ {
		ok, digit := strStartsWithDigit(in[cursor:])
		if ok {
			digits = append(digits, digit)
		}
	}

	return digits
}

func strStartsWithDigit(in string) (bool, string) {
	if _, err := strconv.Atoi(in[0:1]); err == nil {
		return true, in[0:1]
	}
	if i := strings.Index(in, "one"); i == 0 {
		return true, "1"
	} else if i := strings.Index(in, "two"); i == 0 {
		return true, "2"
	} else if i := strings.Index(in, "three"); i == 0 {
		return true, "3"
	} else if i := strings.Index(in, "four"); i == 0 {
		return true, "4"
	} else if i := strings.Index(in, "five"); i == 0 {
		return true, "5"
	} else if i := strings.Index(in, "six"); i == 0 {
		return true, "6"
	} else if i := strings.Index(in, "seven"); i == 0 {
		return true, "7"
	} else if i := strings.Index(in, "eight"); i == 0 {
		return true, "8"
	} else if i := strings.Index(in, "nine"); i == 0 {
		return true, "9"
	}

	return false, ""
}
