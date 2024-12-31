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

type State struct {
	pos     [2]int
	bearing Bearing
}

type QueueItem struct {
	state State
	score int
}

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	walls, start, end, n := parseInput(data)
	fmt.Printf("Input: %d walls, %v start, %v end, n = %d\n", len(walls), start, end, n)

	ok, result := findMinPath(walls, end, State{start, RIGHT})

	if !ok {
		fmt.Println("Did not find a path")
	}

	fmt.Printf("Result: %d", result)
}

func findMinPath(walls map[[2]int]struct{}, end [2]int, start State) (ok bool, score int) {
	scores := make(map[State]int)

	queue := make([]QueueItem, 1)
	queue[0] = QueueItem{start, 0}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if knownScore, known := scores[current.state]; known && knownScore <= current.score {
			continue
		}

		scores[current.state] = current.score
		if current.state.pos == end {
			continue
		}

		forward := move(current.state)
		if _, isWall := walls[forward.pos]; !isWall {
			queue = append(queue, QueueItem{forward, current.score + 1})
		}

		right := turnRight(current.state)
		queue = append(queue, QueueItem{right, current.score + 1000})

		left := turnLeft(current.state)
		queue = append(queue, QueueItem{left, current.score + 1000})
	}

	// Only UP and RIGHT bearings are relevant
	scoreUp, hasUp := scores[State{end, UP}]
	scoreRight, hasRight := scores[State{end, RIGHT}]
	ok = hasUp || hasRight
	if hasUp {
		if hasRight {
			score = min(scoreUp, scoreRight)
		} else {
			score = scoreUp
		}
	} else if hasRight {
		score = scoreRight
	}
	return
}

func move(state State) State {
	switch state.bearing {
	case UP:
		return State{[2]int{state.pos[0] - 1, state.pos[1]}, state.bearing}
	case RIGHT:
		return State{[2]int{state.pos[0], state.pos[1] + 1}, state.bearing}
	case DOWN:
		return State{[2]int{state.pos[0] + 1, state.pos[1]}, state.bearing}
	case LEFT:
		return State{[2]int{state.pos[0], state.pos[1] - 1}, state.bearing}
	}
	panic("Invalid bearing")
}

func turnRight(state State) State {
	return State{state.pos, (state.bearing + 1) % 4}
}

func turnLeft(state State) State {
	return State{state.pos, (state.bearing + 3) % 4}
}

func parseInput(data []uint8) (walls map[[2]int]struct{}, start [2]int, end [2]int, n int) {
	walls = make(map[[2]int]struct{})
	for i, line := range strings.Split(string(data), "\n") {
		n = len(line)
		for j, rune := range line {
			switch rune {
			case 'S':
				start = [2]int{i, j}
			case 'E':
				end = [2]int{i, j}
			case '#':
				walls[[2]int{i, j}] = struct{}{}
			case '.':
				// no-op
			default:
				panic("Unexpected rune on the map: " + string(rune))
			}
		}
	}
	return
}

func min(x int, y int) int {
	if x <= y {
		return x
	}
	return y
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
