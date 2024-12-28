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

	frequencies, n := parseInput(data)
	fmt.Printf("Input: %d frequencies, n = %d\n", len(frequencies), n)

	antinodes := make(map[[2]int]struct{})
	for _, antennas := range frequencies {
		for i, a1 := range antennas {
			for j := i + 1; j < len(antennas); j++ {
				a2 := antennas[j]
				for antinode := range getAllAntinodes(a1, a2, n) {
					antinodes[antinode] = struct{}{}
				}
			}
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if _, ok := antinodes[[2]int{i, j}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
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

func getAllAntinodes(a1, a2 [2]int, n int) (antinodes iter.Seq[[2]int]) {
	return func(yield func([2]int) bool) {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if (a2[1]-a1[1])*(i-a1[0]) == (a2[0]-a1[0])*(j-a1[1]) {
					if !yield([2]int{i, j}) {
						return
					}
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
