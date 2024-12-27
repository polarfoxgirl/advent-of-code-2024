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

type Bearing int

const (
	UP Bearing = iota
	RIGHT
	DOWN
	LEFT
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	blocked, startPos, n := parseInput(data)
	fmt.Printf("Input: %d blocked, (%d, %d) start position, n = %d\n", len(blocked), startPos[0], startPos[1], n)

	sum := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			extra := [2]int{i, j}
			if _, isBlocked := blocked[extra]; !isBlocked && extra != startPos {
				if hasLoop(blocked, startPos, n, extra) {
					sum++
				}
			}
		}
	}

	fmt.Printf("Result: %d", sum)
}

func parseInput(data []uint8) (blocked map[[2]int]struct{}, startPos [2]int, n int) {
	lines := strings.Split(string(data), "\n")
	n = len(lines)
	blocked = make(map[[2]int]struct{})

	for i, line := range lines {
		for j, rune := range line {
			switch rune {
			case '^':
				startPos = [2]int{i, j}
			case '#':
				blocked[[2]int{i, j}] = struct{}{}
			}
		}
	}
	return
}

func move(pos [2]int, bearing Bearing) [2]int {
	switch bearing {
	case UP:
		return [2]int{pos[0] - 1, pos[1]}
	case RIGHT:
		return [2]int{pos[0], pos[1] + 1}
	case DOWN:
		return [2]int{pos[0] + 1, pos[1]}
	case LEFT:
		return [2]int{pos[0], pos[1] - 1}
	}
	panic("Invalid bearing")
}

func hasLoop(blocked map[[2]int]struct{}, startPos [2]int, n int, extraBlock [2]int) bool {
	pos := startPos
	bearing := UP
	visited := make(map[[2]int]map[Bearing]struct{})
	for {
		next := move(pos, bearing)
		if _, isBlocked := blocked[next]; isBlocked || next == extraBlock {
			bearing = turnRight(bearing)
		} else {
			if bearings, isVisited := visited[pos]; isVisited {
				if _, hasBearing := bearings[bearing]; hasBearing {
					return true
				}

				bearings[bearing] = struct{}{}
			} else {
				visited[pos] = make(map[Bearing]struct{})
			}
			pos = next
		}

		if isOutOfBounds(pos, n) {
			break
		}
	}
	return false
}

func isOutOfBounds(pos [2]int, n int) bool {
	return pos[0] < 0 || pos[0] >= n || pos[1] < 0 || pos[1] >= n
}

func turnRight(b Bearing) Bearing {
	return (b + 1) % 4
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
