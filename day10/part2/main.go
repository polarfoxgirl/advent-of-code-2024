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

	topo, trailheads := parseInput(data)
	fmt.Printf("Input: %d trailheads\n", len(trailheads))

	sum := 0
	for trailhead := range trailheads {
		sum += searchForTrails(topo, trailhead[0], trailhead[1], 0)
	}

	fmt.Printf("Result: %d", sum)
}

func parseInput(data []uint8) (topo [][]int, trailheads map[[2]int]struct{}) {
	lines := strings.Split(string(data), "\n")
	topo = make([][]int, len(lines))
	trailheads = make(map[[2]int]struct{})

	for i, line := range lines {
		topo[i] = make([]int, len(lines))
		for j, rune := range line {
			value, e := strconv.ParseInt(string(rune), 0, 32)
			check(e)
			topo[i][j] = int(value)

			if value == 0 {
				trailheads[[2]int{i, j}] = struct{}{}
			}
		}
	}
	return
}

func searchForTrails(topo [][]int, x int, y int, target int) (count int) {
	elevation := topo[x][y]
	if elevation != target {
		return
	}

	if elevation == 9 {
		count = 1
		return
	}

	if x > 0 {
		count += searchForTrails(topo, x-1, y, elevation+1)
	}
	if y > 0 {
		count += searchForTrails(topo, x, y-1, elevation+1)
	}
	if x < len(topo)-1 {
		count += searchForTrails(topo, x+1, y, elevation+1)
	}
	if y < len(topo)-1 {
		count += searchForTrails(topo, x, y+1, elevation+1)
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
