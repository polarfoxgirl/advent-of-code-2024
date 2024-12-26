package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	// "slices"
	// "strconv"
	"iter"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	letters := parseInput(data)

	var sum int
	sum += findXmas(getHorizontalIter(letters))
	sum += findXmas(getVerticalIter(letters))
	sum += findXmas(getRightDiagIter(letters))
	sum += findXmas(getLeftDiagIter(letters))

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

func findXmas(provider iter.Seq[[4]rune]) (wordCount int) {
	for word := range provider {
		if word[0] == 'X' && word[1] == 'M' && word[2] == 'A' && word[3] == 'S' {
			wordCount++
		}
		if word[0] == 'S' && word[1] == 'A' && word[2] == 'M' && word[3] == 'X' {
			wordCount++
		}
	}
	return
}

func getHorizontalIter(letters [][]rune) iter.Seq[[4]rune] {
	return func(yield func([4]rune) bool) {
		n := len(letters)
		for i := 0; i < n; i++ {
			for j := 0; j < n-3; j++ {
				word := [4]rune{letters[i][j], letters[i][j+1], letters[i][j+2], letters[i][j+3]}
				if !yield(word) {
					return
				}
			}
		}
	}
}

func getVerticalIter(letters [][]rune) iter.Seq[[4]rune] {
	return func(yield func([4]rune) bool) {
		n := len(letters)
		for j := 0; j < n; j++ {
			for i := 0; i < n-3; i++ {
				word := [4]rune{letters[i][j], letters[i+1][j], letters[i+2][j], letters[i+3][j]}
				if !yield(word) {
					return
				}
			}
		}
	}
}

func getRightDiagIter(letters [][]rune) iter.Seq[[4]rune] {
	return func(yield func([4]rune) bool) {
		n := len(letters)
		for i := 0; i < n-3; i++ {
			for d := 0; d < n-i-3; d++ {
				word := [4]rune{letters[i+d][d], letters[i+d+1][d+1], letters[i+d+2][d+2], letters[i+d+3][d+3]}
				if !yield(word) {
					return
				}
			}
		}

		for j := 1; j < n; j++ {
			for d := 0; d < n-j-3; d++ {
				word := [4]rune{letters[d][j+d], letters[d+1][j+d+1], letters[d+2][j+d+2], letters[d+3][j+d+3]}
				if !yield(word) {
					return
				}
			}
		}
	}
}

func getLeftDiagIter(letters [][]rune) iter.Seq[[4]rune] {
	return func(yield func([4]rune) bool) {
		n := len(letters)
		for i := 0; i < n-3; i++ {
			for d := 0; d < n-i-3; d++ {
				word := [4]rune{letters[i+d][n-1-d], letters[i+d+1][n-2-d], letters[i+d+2][n-3-d], letters[i+d+3][n-4-d]}
				if !yield(word) {
					return
				}
			}
		}

		for j := n - 2; j > 2; j-- {
			for d := 0; d < j-2; d++ {
				word := [4]rune{letters[d][j-d], letters[d+1][j-1-d], letters[d+2][j-2-d], letters[d+3][j-3-d]}
				if !yield(word) {
					return
				}
			}
		}
	}
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
