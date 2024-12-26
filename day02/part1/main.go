package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	// "slices"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	reports := parseInput(data)

	var sum int
	for _, report := range reports {
		if isSafe(report) {
			sum++
		}
	}

	fmt.Printf("Result: %d", sum)
}

func parseInput(data []uint8) (reports [][]int) {
	lines := strings.Split(string(data), "\n")
	reports = make([][]int, len(lines))

	for i, line := range lines {
		numbers := strings.Split(line, " ")
		levels := make([]int, len(numbers))
		for j, s := range numbers {
			n, e := strconv.ParseInt(s, 0, 32)
			check(e)
			levels[j] = int(n)
		}
		reports[i] = levels
	}
	return
}

func isSafe(report []int) bool {
	order := 0
	for i, x := range report[1:] {
		d := x - report[i]
		if d == 0 || abs(d) > 3 {
			return false
		}

		if order == 0 {
			order = d
		} else if order*d < 0 {
			return false
		}
	}
	return true
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
