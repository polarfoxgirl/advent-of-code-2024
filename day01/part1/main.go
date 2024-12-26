package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	"slices"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	left, right := parseInput(data)
	slices.Sort(left)
	slices.Sort(right)

	var sum int
	for i := 0; i < len(left); i++ {
		sum += abs(left[i] - right[i])
	}

	fmt.Printf("Result: %d", sum)
}

func parseInput(data []uint8) (left []int, right []int) {
	lines := strings.Split(string(data), "\n")
	left = make([]int, len(lines))
	right = make([]int, len(lines))

	for i, line := range lines {
		numbers := strings.Split(line, "   ")
		n0, e0 := strconv.ParseInt(numbers[0], 0, 32)
		check(e0)
		left[i] = int(n0)
		n1, e1 := strconv.ParseInt(numbers[1], 0, 32)
		check(e1)
		right[i] = int(n1)
	}
	return
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
