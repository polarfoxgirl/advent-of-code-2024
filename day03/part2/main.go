package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	// "slices"
	"regexp"
	"strconv"
	// "strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	instructions := string(data)

	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|(do\(\))|(don't\(\))`)

	var sum int
	enabled := true
	for _, matches := range re.FindAllStringSubmatch(instructions, -1) {
		if len(matches[3]) > 0 {
			enabled = true
			continue
		}

		if len(matches[4]) > 0 {
			enabled = false
			continue
		}

		if enabled {
			left, e1 := strconv.ParseInt(matches[1], 0, 32)
			check(e1)
			right, e2 := strconv.ParseInt(matches[2], 0, 32)
			check(e2)
			sum += int(left) * int(right)
		}
	}

	fmt.Printf("Result: %d", sum)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
