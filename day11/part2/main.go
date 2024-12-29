package main

import (
	"fmt"
	"os"
	// "reflect"
	"math"
	// "slices"
	"strconv"
	// "iter"
	"strings"
)

type Outcome struct {
	stone  int64
	blinks int
}

type Split struct {
	stone int64
	ok    bool
	left  int64
	right int64
}

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	stones := parseInput(data)
	fmt.Printf("Input: %d stones\n", len(stones))

	sum := 0
	memo := make(map[Outcome]int)
	splits := make(map[int64]Split)
	for n, stone := range stones {
		sum += blink(Outcome{stone, 75}, memo, splits)
		fmt.Printf("Progress: %d/%d (memo %d, splits %d)\n", n+1, len(stones), len(memo), len(splits))
	}

	fmt.Printf("Result: %d\n", sum)
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

func blink(outcome Outcome, memo map[Outcome]int, splits map[int64]Split) (count int) {
	if outcome.blinks == 0 {
		return 1
	}

	if result, isKnown := memo[outcome]; isKnown {
		return result
	}

	if outcome.stone == 0 {
		count = blink(Outcome{1, outcome.blinks - 1}, memo, splits)
	} else if split := shouldBeSplit(outcome.stone, splits); split.ok {
		count = blink(Outcome{split.left, outcome.blinks - 1}, memo, splits) + blink(Outcome{split.right, outcome.blinks - 1}, memo, splits)
	} else {
		count = blink(Outcome{outcome.stone * 2024, outcome.blinks - 1}, memo, splits)
	}
	memo[outcome] = count
	return
}

func shouldBeSplit(stone int64, splits map[int64]Split) (result Split) {
	if split, isKnown := splits[stone]; isKnown {
		return split
	}

	log := int(math.Log10(float64(stone)))
	if log%2 == 1 {
		result.ok = true

		result.left = stone / int64(math.Pow10(log/2+1))
		result.right = stone % int64(math.Pow10(log/2+1))
	} else {
		result.ok = false
	}

	splits[stone] = result
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
