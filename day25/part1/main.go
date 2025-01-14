package main

import (
	"fmt"
	"os"
	"strings"
)

const EMPTY = "....."
const FULL = "#####"

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	locks, keys := parseInput(data)
	fmt.Printf("Input: %d locks, %d keys\n", len(locks), len(keys))

	result := 0
	for _, lock := range locks {
		for _, key := range keys {
			if fit(lock, key) {
				result++
			}
		}
	}

	fmt.Printf("Result: %d", result)
}

func fit(lock [5]int, key [5]int) bool {
	for i := 0; i < 5; i++ {
		if lock[i]+key[i] > 5 {
			return false
		}
	}
	return true
}

func parseInput(data []uint8) (locks [][5]int, keys [][5]int) {
	locks = make([][5]int, 0)
	keys = make([][5]int, 0)

	for _, section := range strings.Split(string(data), "\n\n") {
		code, isLock := parseSchematic(section)
		if isLock {
			locks = append(locks, code)
		} else {
			keys = append(keys, code)
		}
	}
	return
}

func parseSchematic(text string) (code [5]int, isLock bool) {
	lines := strings.Split(text, "\n")
	if len(lines) != 7 {
		panic("Invalid schematic length")
	}

	if lines[0] == FULL && lines[6] == EMPTY {
		isLock = true
	} else if lines[0] != EMPTY || lines[6] != FULL {
		panic("Schematic is neither lock nor key")
	}

	for _, line := range lines[1:6] {
		for i, rune := range line {
			switch rune {
			case '#':
				code[i]++
			case '.':
				// no-op
			default:
				panic("Unexpected rune in schematic: " + string(rune))
			}
		}
	}
	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
