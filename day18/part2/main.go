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
)

func main() {
	// n := 7
	// knownSteps := 12
	// data, err := os.ReadFile("test.txt")
	n := 71
	knownSteps := 1024
	data, err := os.ReadFile("input.txt")
	check(err)

	bytes := parseInput(data)
	fmt.Printf("Input: %d bytes, n = %d\n", len(bytes), n)

	for steps := knownSteps + 1; steps <= len(bytes); steps++ {
		corrupted := calcCorrupted(bytes, steps)

		if hasPath := runDijkstra(corrupted, n); !hasPath {
			fmt.Printf("Result: %v\n", bytes[steps-1])
			break
		}
	}
}

func calcCorrupted(bytes [][2]int, steps int) (corrupted map[[2]int]struct{}) {
	corrupted = make(map[[2]int]struct{})

	for i := 0; i < steps; i++ {
		corrupted[bytes[i]] = struct{}{}
	}

	return
}

func runDijkstra(corrupted map[[2]int]struct{}, n int) bool {
	visited := make(map[[2]int]struct{})
	finish := [2]int{n - 1, n - 1}

	priorityQueue := make(map[int][][2]int, 1)
	priorityQueue[0] = [][2]int{{0, 0}}
	scoreWatermark := 0

	for len(priorityQueue) > 0 {
		current := pop(priorityQueue, &scoreWatermark)
		if _, isVisited := visited[current]; isVisited {
			continue
		}
		visited[current] = struct{}{}

		if current == finish {
			return true
		}

		for _, next := range getMoves(current, n) {
			if _, isCorrupted := corrupted[next]; !isCorrupted {
				queueUp(priorityQueue, scoreWatermark+1, next)
			}
		}
	}

	return false
}

func getMoves(pos [2]int, n int) (moves [][2]int) {
	moves = make([][2]int, 0)

	if pos[0] > 0 {
		moves = append(moves, [2]int{pos[0] - 1, pos[1]})
	}
	if pos[0] < n-1 {
		moves = append(moves, [2]int{pos[0] + 1, pos[1]})
	}
	if pos[1] > 0 {
		moves = append(moves, [2]int{pos[0], pos[1] - 1})
	}
	if pos[1] < n-1 {
		moves = append(moves, [2]int{pos[0], pos[1] + 1})
	}

	return
}

func pop(priorityQueue map[int][][2]int, scoreWatermark *int) (result [2]int) {
	for len(priorityQueue) > 0 {
		// Find the next score-level queue that is not empty
		if subQueue, ok := priorityQueue[*scoreWatermark]; ok {
			result = subQueue[0]
			if len(subQueue) > 1 {
				priorityQueue[*scoreWatermark] = subQueue[1:]
			} else {
				delete(priorityQueue, *scoreWatermark)
			}
			return
		}
		*scoreWatermark++
	}
	panic("The queue is empty")
}

func queueSize(priorityQueue map[int][][2]int) (size int) {
	for _, subQueue := range priorityQueue {
		size += len(subQueue)
	}

	return
}

func queueUp(priorityQueue map[int][][2]int, score int, state [2]int) {
	subQueue, ok := priorityQueue[score]
	if !ok {
		priorityQueue[score] = [][2]int{state}
	} else {
		priorityQueue[score] = append(subQueue, state)
	}
}

func parseInput(data []uint8) (bytes [][2]int) {
	lines := strings.Split(string(data), "\n")
	bytes = make([][2]int, len(lines))

	re := regexp.MustCompile(`(\d+),(\d+)`)

	for i, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) != 3 {
			panic("Invalid regexp matches")
		}

		bytes[i] = [2]int{parseInt(matches[1]), parseInt(matches[2])}
	}
	return
}

func parseInt(s string) int {
	result, e := strconv.ParseInt(s, 0, 64)
	check(e)
	return int(result)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
