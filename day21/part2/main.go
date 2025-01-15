package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// NOTE: for [2]int coordinates {0, 0} is bottom left corner
type state struct {
	keyPad   [2]int
	progress int
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

const ITERATIONS = 25

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	codes := parseInput(data)
	fmt.Printf("Input: %d codes\n", len(codes))

	conMap := getConMap()
	reverseDirButtons := getReverseDirButtons()

	// Only leave paths with the most efficient 25th "derivatives"
	conMap = pruneMap(conMap, reverseDirButtons)
	lenMap := calcLenMap(conMap, reverseDirButtons)

	result := 0
	for _, code := range codes {
		result += processNumCode(lenMap, reverseDirButtons, code)
	}

	fmt.Printf("Result: %d\n", result)
}

func getConMap() (conMap map[[2][2]int][][]dirButton) {
	conMap = make(map[[2][2]int][][]dirButton)
	for x1 := 0; x1 < 3; x1++ {
		for y1 := 0; y1 < 2; y1++ {
			pos1 := [2]int{x1, y1}
			if _, ok1 := dirButtons[pos1]; ok1 {
				for x2 := 0; x2 < 3; x2++ {
					for y2 := 0; y2 < 2; y2++ {
						pos2 := [2]int{x2, y2}
						if target, ok2 := dirButtons[pos2]; ok2 {
							conMap[[2][2]int{pos1, pos2}] = runDijkstra(1, pos1, getUpdateFn([]dirButton{target}, dirButtons))
						}
					}
				}
			}
		}
	}
	return
}

func pruneMap(conMap map[[2][2]int][][]dirButton, reverseDirButtons map[dirButton][2]int) map[[2][2]int][][]dirButton {
	pruned := conMap
	cache := make(map[string][][]dirButton)
	for _, paths := range conMap {
		for _, path := range paths {
			cache[printCode(path)] = [][]dirButton{path}
		}
	}

	// Can combine pruning and length calculation
	shouldPrune := true
	for k := 0; k < 25 && shouldPrune; k++ {
		current := pruned
		pruned = make(map[[2][2]int][][]dirButton)
		shouldPrune = false

		for pair, paths := range current {
			if len(paths) > 1 {
				minLen := 0
				goodPaths := make([][]dirButton, 0)
				for _, path := range paths {
					codStr := printCode(path)
					results := processDirPad(current, reverseDirButtons, cache[codStr])
					cache[codStr] = results

					pathMinLen := len(results[0])
					if minLen == 0 {
						// Init
						goodPaths = append(goodPaths, path)
						minLen = pathMinLen
					} else if minLen > pathMinLen {
						// Replace
						goodPaths = [][]dirButton{path}
						minLen = pathMinLen
					} else if minLen == pathMinLen {
						// Append
						goodPaths = append(goodPaths, path)
					} else {
						// Ignore
					}
				}

				pruned[pair] = goodPaths
				shouldPrune = true
			} else {
				pruned[pair] = paths
			}
		}
	}

	// After this we expect one path per pair
	for _, paths := range pruned {
		if len(paths) > 1 {
			panic("Unsuccessful pruning!")
		}
	}
	return pruned
}

func calcLenMap(pruned map[[2][2]int][][]dirButton, reverseDirButtons map[dirButton][2]int) (lenMap map[[2][2]int]int) {
	lenMap = make(map[[2][2]int]int)
	for pair, paths := range pruned {
		lenMap[pair] = len(paths[0])
	}

	// One less iteration since we already did one when populating the map
	for k := 0; k < ITERATIONS-1; k++ {
		nextLenMap := make(map[[2][2]int]int)

		for pair, paths := range pruned {
			// We always come back to A after button cycle
			current := reverseDirButtons[PRESS]
			pathTotal := 0
			for _, button := range paths[0] {
				target := reverseDirButtons[button]

				pathTotal += lenMap[[2][2]int{current, target}]
				current = target
			}
			nextLenMap[pair] = pathTotal
		}

		lenMap = nextLenMap
	}

	return
}

func getReverseDirButtons() (reverseDirButtons map[dirButton][2]int) {
	reverseDirButtons = make(map[dirButton][2]int)
	for key, value := range dirButtons {
		reverseDirButtons[value] = key
	}
	return
}

func processNumCode(lenMap map[[2][2]int]int, reverseDirButtons map[dirButton][2]int, code [4]numButton) int {
	nextCodes := runDijkstra(len(code), [2]int{2, 0}, getUpdateFn(code[:], numButtons))

	minTotal := 0
	for _, nextCode := range nextCodes {
		total := calcDirCodeLen(lenMap, reverseDirButtons, nextCode)

		if minTotal == 0 || minTotal > total {
			minTotal = total
		}
	}
	return getComplexity(code, minTotal)
}

func processDirPad(conMap map[[2][2]int][][]dirButton, reverseDirButtons map[dirButton][2]int, codes [][]dirButton) (allDirCodes [][]dirButton) {
	minLen := 0
	for _, code := range codes {
		dirCodes := processDirCode(conMap, reverseDirButtons, code)

		// We expect all codes for the same sequence to be the same
		codeMinLen := len(dirCodes[0])
		if minLen == 0 {
			// Init
			allDirCodes = dirCodes
			minLen = codeMinLen
		} else if minLen > codeMinLen {
			// Replace
			allDirCodes = dirCodes
			minLen = codeMinLen
		} else if minLen == codeMinLen {
			// Append
			allDirCodes = slices.Concat(allDirCodes, dirCodes)
		} else {
			// Ignore
		}
	}

	return allDirCodes
}

func processDirCode(conMap map[[2][2]int][][]dirButton, reverseDirButtons map[dirButton][2]int, code []dirButton) (dirCodes [][]dirButton) {
	dirCodes = [][]dirButton{{}}
	current := [2]int{2, 1}
	for _, button := range code {
		target := reverseDirButtons[button]
		extendedDirCodes := make([][]dirButton, 0)
		for _, suffix := range conMap[[2][2]int{current, target}] {
			extendedDirCodes = slices.Concat(extendedDirCodes, genPaths(dirCodes, suffix))
		}
		dirCodes = extendedDirCodes
		current = target
	}
	return
}

func calcDirCodeLen(lenMap map[[2][2]int]int, reverseDirButtons map[dirButton][2]int, code []dirButton) (total int) {
	current := reverseDirButtons[PRESS]
	for _, button := range code {
		target := reverseDirButtons[button]
		total += lenMap[[2][2]int{current, target}]
		current = target
	}
	return
}

func getUpdateFn[T dirButton | numButton](code []T, buttonMap map[[2]int]T) func(button dirButton, current state) (next state, ok bool) {
	return func(button dirButton, current state) (next state, ok bool) {
		if button == PRESS {
			targetButton, validButton := buttonMap[current.keyPad]
			if !validButton {
				panic(fmt.Sprintf("Attempting to press invalid button %v on the key pad", current.keyPad))
			}

			if code[current.progress] == targetButton {
				ok = true
				next = state{
					keyPad:   current.keyPad,
					progress: current.progress + 1,
				}
			}

		} else {
			// Moving num pad arm
			target := pressDirMove(button, current.keyPad)
			if _, validButton := buttonMap[target]; validButton {
				ok = true
				next = state{
					keyPad:   target,
					progress: current.progress,
				}
			}
		}
		return
	}
}

func pressDirMove(button dirButton, current [2]int) [2]int {
	shift, isMove := moves[button]
	if !isMove {
		panic("Trying to do a move with non-move dir button")
	}

	return [2]int{current[0] + shift[0], current[1] + shift[1]}
}

func runDijkstra(codeLen int, initPos [2]int, updateFn func(button dirButton, current state) (next state, ok bool)) [][]dirButton {
	priorityQueue := make(map[int][]state, 1)
	visited := make(map[state]struct{})
	paths := make(map[state][][]dirButton)
	pathsLen := make(map[state]int)

	start := state{
		keyPad:   initPos,
		progress: 0,
	}
	priorityQueue[0] = []state{start}
	paths[start] = [][]dirButton{{}}
	pathsLen[start] = 0
	scoreWatermark := 0
	var end *state

	for len(priorityQueue) > 0 {
		current := pop(priorityQueue, &scoreWatermark)

		if _, isVisited := visited[current]; isVisited {
			continue
		}
		visited[current] = struct{}{}
		currentPaths := paths[current]
		currentPathLens := pathsLen[current]

		if current.progress == codeLen {
			end = &current
		} else {
			for _, button := range [5]dirButton{UP, DOWN, RIGHT, LEFT, PRESS} {
				if next, ok := updateFn(button, current); ok {
					queueUp(priorityQueue, scoreWatermark+1, next)

					if knownPathLen, hasPaths := pathsLen[next]; hasPaths {
						if knownPathLen == currentPathLens+1 {
							// Append
							paths[next] = slices.Concat(paths[next], genPaths(currentPaths, []dirButton{button}))
						}
					} else {
						// Init
						paths[next] = genPaths(currentPaths, []dirButton{button})
						pathsLen[next] = currentPathLens + 1
					}
				}
			}
		}
	}

	if end == nil {
		panic("Min path not found!")
	}
	return paths[*end]
}

func genPaths(current [][]dirButton, suffix []dirButton) (result [][]dirButton) {
	result = make([][]dirButton, len(current))
	for i, path := range current {
		result[i] = slices.Concat(slices.Clone(path), suffix)
	}
	return
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
	codeValue := (int(code[0])*100 + int(code[1])*10 + int(code[2]))
	fmt.Printf("Complexity is %d * %d\n", codeValue, minLen)
	return codeValue * minLen
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

func printCode(code []dirButton) string {
	var b strings.Builder
	for _, button := range code {
		switch button {
		case UP:
			fmt.Fprint(&b, "^")
		case DOWN:
			fmt.Fprint(&b, "v")
		case LEFT:
			fmt.Fprint(&b, "<")
		case RIGHT:
			fmt.Fprint(&b, ">")
		case PRESS:
			fmt.Fprint(&b, "A")
		default:
			panic("Can't print invalid button")
		}
	}
	return b.String()
}

func parseCode(text string) (code []dirButton) {
	code = make([]dirButton, len(text))
	for i, rune := range text {
		switch rune {
		case '^':
			code[i] = UP
		case 'v':
			code[i] = DOWN
		case '<':
			code[i] = LEFT
		case '>':
			code[i] = RIGHT
		case 'A':
			code[i] = PRESS
		default:
			panic("Can't parse code")
		}
	}
	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
