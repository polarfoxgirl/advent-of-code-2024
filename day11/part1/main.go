package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	// "slices"
	"strconv"
	// "iter"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	stones := parseInput(data)
	fmt.Printf("Input: %d stones\n", len(stones))

	for n := 0; n < 25; n++ {
		stones = blink(stones)
	}

	fmt.Printf("Result: %d\n", len(stones))
}

func parseInput(data []uint8) (stones []int64) {
	fields := strings.Fields(string(data))
	stones = make([]int64, len(fields))

	for i, field := range fields {
		value, e := strconv.ParseInt(field, 0, 64)
		check(e)
		stones[i] = value
	}
	return
}

func blink(stones []int64) (newStones []int64) {
	newStones = make([]int64, 0)
	for _, stone := range stones {
		if stone == 0 {
			newStones = append(newStones, 1)
		} else if str := strconv.FormatInt(stone, 10); len(str)%2 == 0 {
			left, e := strconv.ParseInt(str[:len(str)/2], 0, 64)
			check(e)
			newStones = append(newStones, left)

			trimmed := strings.TrimLeft(str[len(str)/2:], "0")
			if len(trimmed) == 0 {
				newStones = append(newStones, 0)
			} else {
				right, e := strconv.ParseInt(trimmed, 0, 64)
				check(e)
				newStones = append(newStones, right)
			}
		} else {
			newStones = append(newStones, stone*2024)
		}
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
