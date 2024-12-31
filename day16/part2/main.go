package main

import (
	"fmt"
	// "maps"
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

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	walls, start, end, n := parseInput(data)
	fmt.Printf("Input: %d walls, %v start, %v end, n = %d\n", len(walls), start, end, n)

	scores := findMinPaths(walls, end, State{start, RIGHT})

	ok, score, endStates := getFinalScore(scores, end)

	if !ok {
		fmt.Println("Did not find a path")
	} else if len(endStates) > 1 {
		fmt.Println("Multiple end states, unlucky!")
	}

	// Cheat and assume there is only one end state
	seats := findSeats(walls, endStates[0], scores, score)

	print(walls, seats, n)

	fmt.Printf("Result: %d, end state(s): %v, seats: %d", score, endStates, len(seats))
}

func findMinPaths(walls map[[2]int]struct{}, end [2]int, start State) (scores map[State]int) {
	scores = make(map[State]int)

	// Dijkstra FTW!
	priorityQueue := make(map[int][]State, 1)
	priorityQueue[0] = []State{start}
	scoreWatermark := 0

	for len(priorityQueue) > 0 {
		currentState := pop(priorityQueue, &scoreWatermark)

		if knownScore, known := scores[currentState]; known {
			if knownScore > scoreWatermark {
				panic(fmt.Sprintf("Somehow arrived at known score %d with lower score %d", knownScore, scoreWatermark))
			}
			// There is already a cheaper way to arrive at this state
			continue
		}

		scores[currentState] = scoreWatermark

		if currentState.pos == end {
			// Finish line!
			continue
		}

		forward := move(currentState)
		if _, isWall := walls[forward.pos]; !isWall {
			queueUp(priorityQueue, scoreWatermark+1, forward)
		}

		turnScore := scoreWatermark + 1000
		right := turnRight(currentState)
		queueUp(priorityQueue, turnScore, right)

		left := turnLeft(currentState)
		queueUp(priorityQueue, turnScore, left)
	}

	return
}

func getFinalScore(scores map[State]int, end [2]int) (ok bool, score int, endStates []State) {
	// Only UP and RIGHT bearings are relevant
	endUp := State{end, UP}
	endRight := State{end, RIGHT}

	scoreUp, hasUp := scores[endUp]
	scoreRight, hasRight := scores[endRight]
	ok = hasUp || hasRight

	if hasUp {
		if hasRight {
			if scoreUp == scoreRight {
				score = scoreUp
				endStates = []State{endUp, endRight}
			} else if scoreUp < scoreRight {
				score = scoreUp
				endStates = []State{endUp}
			} else {
				score = scoreRight
				endStates = []State{endRight}
			}
		} else {
			score = scoreUp
			endStates = []State{endUp}
		}
	} else if hasRight {
		score = scoreRight
		endStates = []State{endRight}
	}
	return
}

func findSeats(walls map[[2]int]struct{}, end State, scores map[State]int, targetScore int) (seats map[[2]int]struct{}) {
	seats = make(map[[2]int]struct{})

	priorityQueue := make(map[int][]State, 1)
	priorityQueue[0] = []State{end}
	diffWatermark := 0

	for len(priorityQueue) > 0 {
		currentState := pop(priorityQueue, &diffWatermark)

		score, hasScore := scores[currentState]
		if !hasScore {
			// Unexplored seats are not on an optimum path
			continue
		}

		if score+diffWatermark != targetScore {
			// Not on an optimum path
			continue
		}

		seats[currentState.pos] = struct{}{}

		forward := reverse(currentState)
		if _, isWall := walls[forward.pos]; !isWall {
			queueUp(priorityQueue, diffWatermark+1, forward)
		}

		turnScore := diffWatermark + 1000
		right := turnRight(currentState)
		queueUp(priorityQueue, turnScore, right)

		left := turnLeft(currentState)
		queueUp(priorityQueue, turnScore, left)
	}

	return
}

func pop(priorityQueue map[int][]State, scoreWatermark *int) (result State) {
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

func queueUp(priorityQueue map[int][]State, score int, state State) {
	subQueue, ok := priorityQueue[score]
	if !ok {
		priorityQueue[score] = []State{state}
	} else {
		priorityQueue[score] = append(subQueue, state)
	}
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

func reverse(state State) State {
	fakeMove := move(State{state.pos, (state.bearing + 2) % 4})
	return State{fakeMove.pos, state.bearing}
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

func print(walls map[[2]int]struct{}, seats map[[2]int]struct{}, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			pos := [2]int{i, j}
			if _, isWall := walls[pos]; isWall {
				fmt.Print("#")
			} else if _, isSeat := seats[pos]; isSeat {
				fmt.Print("O")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
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
