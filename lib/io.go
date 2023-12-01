package io

import (
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return strings.TrimSpace(string(dat))
}
