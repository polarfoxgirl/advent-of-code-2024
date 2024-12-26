package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	// "slices"
	// "strconv"
	// "iter"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	letters := parseInput(data)

	var sum int
	n := len(letters)
	for i := 0; i < n-2; i++ {
		for j := 0; j < n-2; j++ {
			if checkXmas(letters, i, j) {
				sum++
			}
		}
	}

	fmt.Printf("Result: %d", sum)
}

func parseInput(data []uint8) (letters [][]rune) {
	lines := strings.Split(string(data), "\n")
	letters = make([][]rune, len(lines))

	for i, line := range lines {
		letters[i] = make([]rune, len(line))
		for j, rune := range line {
			letters[i][j] = rune
		}
	}
	return
}

func checkXmas(letters [][]rune, x int, y int) bool {
	if letters[x+1][y+1] != 'A' {
		return false
	}

	if letters[x][y] == 'M' && letters[x+2][y+2] == 'S' {
		if letters[x][y+2] == 'M' {
			return letters[x+2][y] == 'S'
		} else if letters[x+2][y] == 'M' {
			return letters[x][y+2] == 'S'
		}
	} else if letters[x][y] == 'S' && letters[x+2][y+2] == 'M' {
		if letters[x][y+2] == 'M' {
			return letters[x+2][y] == 'S'
		} else if letters[x+2][y] == 'M' {
			return letters[x][y+2] == 'S'
		}
	}

	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
