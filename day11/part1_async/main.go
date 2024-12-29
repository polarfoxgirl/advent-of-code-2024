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
	"sync"
	sa "sync/atomic"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	stones := parseInput(data)
	fmt.Printf("Input: %d stones\n", len(stones))

	var wg sync.WaitGroup
	var sum sa.Int32
	for _, stone := range stones {
		wg.Add(1)
		blink(stone, 25, &sum, &wg)
	}

	wg.Wait()
	fmt.Printf("Result: %d\n", sum.Load())
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

func blink(stone int64, blinks int, sum *sa.Int32, wg *sync.WaitGroup) {
	defer wg.Done()
	if blinks == 0 {
		sum.Add(1)
		return
	}

	if stone == 0 {
		wg.Add(1)
		go blink(1, blinks-1, sum, wg)
	} else if ok, left, right := shouldBeSplit2(stone); ok {
		wg.Add(2)
		go blink(left, blinks-1, sum, wg)
		go blink(right, blinks-1, sum, wg)
	} else {
		wg.Add(1)
		go blink(stone*2024, blinks-1, sum, wg)
	}
}

// func blink(stones []int64) (newStones []int64) {
// 	newStones = make([]int64, 0)
// 	for _, stone := range stones {
// 		if stone == 0 {
// 			newStones = append(newStones, 1)
// 		} else if ok, left, right := shouldBeSplit2(stone); ok {
// 			newStones = append(newStones, left)
// 			newStones = append(newStones, right)
// 		} else {
// 			newStones = append(newStones, stone*2024)
// 		}
// 	}
// 	return
// }

func shouldBeSplit(stone int64) (ok bool, left int64, right int64) {
	str := strconv.FormatInt(stone, 10)
	if len(str)%2 == 0 {
		ok = true
		var e error

		left, e = strconv.ParseInt(str[:len(str)/2], 0, 64)
		check(e)

		trimmed := strings.TrimLeft(str[len(str)/2:], "0")
		if len(trimmed) == 0 {
			right = 0
		} else {
			right, e = strconv.ParseInt(trimmed, 0, 64)
			check(e)
		}
	}
	return
}

func shouldBeSplit2(stone int64) (ok bool, left int64, right int64) {
	log := int(math.Log10(float64(stone)))
	if log%2 == 1 {
		ok = true

		left = stone / int64(math.Pow10(log/2+1))
		right = stone % int64(math.Pow10(log/2+1))
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
