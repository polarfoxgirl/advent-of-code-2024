package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// NOTE: for [2]int coordinates {0, 0} is bottom left corner
type state struct {
	numKeyPad  [2]int
	dirKeyPads [2][2]int
	progress   int
}

type dirButton int

const (
	UP dirButton = iota
	DOWN
	RIGHT
	LEFT
	PRESS
)

var moves = map[dirButton][2]int{
	UP:    {0, 1},
	DOWN:  {0, -1},
	RIGHT: {1, 0},
	LEFT:  {-1, 0},
}

var dirButtons = map[[2]int]dirButton{
	{1, 1}: UP,
	{1, 0}: DOWN,
	{2, 0}: RIGHT,
	{0, 0}: LEFT,
	{2, 1}: PRESS,
}

type numButton int

const (
	SEVEN  numButton = 7
	EIGHT  numButton = 8
	NINE   numButton = 9
	FOUR   numButton = 4
	FIVE   numButton = 5
	SIX    numButton = 6
	ONE    numButton = 1
	TWO    numButton = 2
	THREE  numButton = 3
	ZERO   numButton = 0
	ACTION numButton = -1
)

var numButtons = map[[2]int]numButton{
	{0, 3}: SEVEN,
	{1, 3}: EIGHT,
	{2, 3}: NINE,
	{0, 2}: FOUR,
	{1, 2}: FIVE,
	{2, 2}: SIX,
	{0, 1}: ONE,
	{1, 1}: TWO,
	{2, 1}: THREE,
	{1, 0}: ZERO,
	{2, 0}: ACTION,
}

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	codes := parseInput(data)
	fmt.Printf("Input: %d codes\n", len(codes))

	result := 0
	for _, code := range codes {
		minPresses := runDijkstra(code)
		fmt.Printf("Need %d presses for code %v\n", minPresses, code)
		result += getComplexity(code, minPresses)
	}

	fmt.Printf("Result: %d\n", result)
}

func press(button dirButton, current state, code [4]numButton) (next state, ok bool) {
	if button == PRESS {
		secondDirButton, validSecondDirButton := dirButtons[current.dirKeyPads[1]]
		if !validSecondDirButton {
			panic(fmt.Sprintf("Attempting to press invalid button %v on the second dir pad", current.dirKeyPads[1]))
		}

		if secondDirButton == PRESS {
			firstDirButton, validFirstDirButton := dirButtons[current.dirKeyPads[0]]
			if !validFirstDirButton {
				panic(fmt.Sprintf("Attempting to press invalid button %v on the first dir pad", current.dirKeyPads[0]))
			}

			if firstDirButton == PRESS {
				numButton, validNumButton := numButtons[current.numKeyPad]
				if !validNumButton {
					panic(fmt.Sprintf("Attempting to press invalid button %v on the num pad", current.numKeyPad))
				}

				if code[current.progress] == numButton {
					ok = true
					next = state{
						numKeyPad:  current.numKeyPad,
						dirKeyPads: current.dirKeyPads,
						progress:   current.progress + 1,
					}
				}

			} else {
				// Moving num pad arm
				num := pressDirMove(firstDirButton, current.numKeyPad)
				if _, validNum := numButtons[num]; validNum {
					ok = true
					next = state{
						numKeyPad:  num,
						dirKeyPads: current.dirKeyPads,
						progress:   current.progress,
					}
				}
			}
		} else {
			// Moving first dir pad arm
			firstDir := pressDirMove(secondDirButton, current.dirKeyPads[0])
			if _, valid := dirButtons[firstDir]; valid {
				ok = true
				next = state{
					numKeyPad:  current.numKeyPad,
					dirKeyPads: [2][2]int{firstDir, current.dirKeyPads[1]},
					progress:   current.progress,
				}
			}
		}
	} else {
		// Moving second dir pad arm
		secondDir := pressDirMove(button, current.dirKeyPads[1])
		if _, valid := dirButtons[secondDir]; valid {
			ok = true
			next = state{
				numKeyPad:  current.numKeyPad,
				dirKeyPads: [2][2]int{current.dirKeyPads[0], secondDir},
				progress:   current.progress,
			}
		}
	}
	return
}

func pressDirMove(button dirButton, current [2]int) [2]int {
	shift, isMove := moves[button]
	if !isMove {
		panic("Trying to do a move with non-move dir button")
	}

	return [2]int{current[0] + shift[0], current[1] + shift[1]}
}

func runDijkstra(code [4]numButton) int {
	priorityQueue := make(map[int][]state, 1)
	start := state{
		numKeyPad:  [2]int{2, 0},
		dirKeyPads: [2][2]int{{2, 1}, {2, 1}},
		progress:   0,
	}
	priorityQueue[0] = []state{start}
	scoreWatermark := 0

	visited := make(map[state]struct{})

	for len(priorityQueue) > 0 {
		current := pop(priorityQueue, &scoreWatermark)

		if _, isVisited := visited[current]; isVisited {
			continue
		}
		visited[current] = struct{}{}

		if current.progress == len(code) {
			// fmt.Println(currentPath)
			return scoreWatermark
		} else {
			for _, button := range [5]dirButton{UP, DOWN, RIGHT, LEFT, PRESS} {
				if next, ok := press(button, current, code); ok {
					queueUp(priorityQueue, scoreWatermark+1, next)
				}
			}
		}
	}

	panic("Min path not found!")
}

func pop(priorityQueue map[int][]state, scoreWatermark *int) (result state) {
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

func queueUp(priorityQueue map[int][]state, score int, item state) {
	subQueue, ok := priorityQueue[score]
	if !ok {
		priorityQueue[score] = []state{item}
	} else {
		priorityQueue[score] = append(subQueue, item)
	}
}

func getComplexity(code [4]numButton, minLen int) int {
	return (int(code[0])*100 + int(code[1])*10 + int(code[2])) * minLen
}

func parseInput(data []uint8) (codes [][4]numButton) {
	lines := strings.Split(string(data), "\n")
	codes = make([][4]numButton, len(lines))

	re := regexp.MustCompile(`(\d)(\d)(\d)A`)

	for i, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) != 4 {
			panic("Invalid regexp matches")
		}

		codes[i] = [4]numButton{
			parseNumButton(matches[1]),
			parseNumButton(matches[2]),
			parseNumButton(matches[3]),
			ACTION,
		}
	}
	return
}

func parseNumButton(s string) numButton {
	result, e := strconv.ParseInt(s, 0, 64)
	check(e)
	return numButton(int(result))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
