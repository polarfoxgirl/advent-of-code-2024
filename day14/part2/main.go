package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	// "slices"
	// "iter"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	// lenX := 11
	// lenY := 7
	lenX := 101
	lenY := 103

	robots, velocities := parseInput(data)
	fmt.Printf("Input: %d robots, %d velocities\n", len(robots), len(velocities))

	for s := 1; ; s++ {
		move(robots, velocities, lenX, lenY)
		if isWeird(robots) {
			fmt.Printf("After %d seconds:\n", s)
			print(robots, lenX, lenY)
			time.Sleep(time.Second)
		}
	}
}

func move(robots [][2]int, velocities [][2]int, lenX int, lenY int) {
	for i := range robots {
		robots[i][0] = mod(robots[i][0]+velocities[i][0], lenX)
		robots[i][1] = mod(robots[i][1]+velocities[i][1], lenY)
	}
}

func isWeird(robots [][2]int) bool {
	xCount := make(map[int]int)
	yCount := make(map[int]int)

	for _, robot := range robots {
		xCount[robot[0]]++
		yCount[robot[1]]++
	}

	for _, value := range xCount {
		if value > 20 {
			return true
		}
	}

	for _, value := range yCount {
		if value > 20 {
			return true
		}
	}

	return false
}

func print(robots [][2]int, lenX int, lenY int) {
	for y := 0; y < lenY; y++ {
		for x := 0; x < lenX; x++ {
			count := 0
			for _, robot := range robots {
				if robot[0] == x && robot[1] == y {
					count++
				}
			}

			if count > 9 {
				fmt.Print("@")
			} else if count > 0 {
				fmt.Print(count)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func parseInput(data []uint8) (robots [][2]int, velocities [][2]int) {
	lines := strings.Split(string(data), "\n")
	robots = make([][2]int, len(lines))
	velocities = make([][2]int, len(lines))

	re := regexp.MustCompile(`p=(\d+),(\d+) v=(\-?\d+),(\-?\d+)`)

	for i, line := range lines {
		matches := re.FindStringSubmatch(line)

		robots[i] = [2]int{parseInt(matches[1]), parseInt(matches[2])}
		velocities[i] = [2]int{parseInt(matches[3]), parseInt(matches[4])}
	}
	return
}

func mod(x int, m int) (result int) {
	result = x % m
	if result < 0 {
		result += m
	}
	return
}

func parseInt(s string) int {
	result, e := strconv.ParseInt(s, 0, 64)
	check(e)
	return int(result)
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
