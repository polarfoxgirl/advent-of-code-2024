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

	frequencies, n := parseInput(data)
	fmt.Printf("Input: %d frequencies, n = %d\n", len(frequencies), n)

	antinodes := make(map[[2]int]struct{})
	for _, antennas := range frequencies {
		for i, a1 := range antennas {
			for j := i + 1; j < len(antennas); j++ {
				a2 := antennas[j]
				for _, antinode := range getPairAntinodes(a1, a2, n) {
					antinodes[antinode] = struct{}{}
				}
			}
		}
	}

	fmt.Printf("Result: %d", len(antinodes))
}

func parseInput(data []uint8) (frequencies map[rune][][2]int, n int) {
	lines := strings.Split(string(data), "\n")
	n = len(lines)
	frequencies = make(map[rune][][2]int)

	for i, line := range lines {
		for j, rune := range line {
			if rune != '.' {
				if locations, ok := frequencies[rune]; ok {
					frequencies[rune] = append(locations, [2]int{i, j})
				} else {
					frequencies[rune] = [][2]int{{i, j}}
				}
			}
		}
	}
	return
}

func getPairAntinodes(a1, a2 [2]int, n int) (antinodes [][2]int) {
	antinodes = make([][2]int, 0)

	first := [2]int{2*a1[0] - a2[0], 2*a1[1] - a2[1]}
	if inBounds(first, n) {
		antinodes = append(antinodes, first)
	}

	second := [2]int{2*a2[0] - a1[0], 2*a2[1] - a1[1]}
	if inBounds(second, n) {
		antinodes = append(antinodes, second)
	}

	return
}

func inBounds(point [2]int, n int) bool {
	return point[0] >= 0 && point[0] < n && point[1] >= 0 && point[1] < n
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
