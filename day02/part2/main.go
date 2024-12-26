package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	// "slices"
	"iter"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	reports := parseInput(data)

	var sum int
	for _, report := range reports {
		if isSafeWithDampener(report) {
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

func isSafeWithDampener(report []int) bool {

	// Since -1 is not in range, this does not skip any elements
	if isSafe(getSkipIter(report, -1)) {
		return true
	}

	for i := 0; i < len(report); i++ {
		if isSafe(getSkipIter(report, i)) {
			return true
		}
	}

	return false
}

func getSkipIter(report []int, iSkip int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i, x := range report {
			if i != iSkip {
				if !yield(x) {
					return
				}
			}
		}
	}
}

func isSafe(report iter.Seq[int]) bool {
	prev := -1
	order := 0
	for x := range report {
		if prev >= 0 {
			d := x - prev
			if d == 0 || abs(d) > 3 {
				return false
			}

			if order == 0 {
				order = d
			} else if order*d < 0 {
				return false
			}
		}
		prev = x
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
