package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	// "reflect"
	// "math"
	// "slices"
	// "strconv"
	// "iter"
	"strings"
)

const CHEAT_LEN = 6
const MIN_SAVINGS = 50

type cheat struct {
	end [2]int
	len int
}

func main() {
	data, err := os.ReadFile("test.txt")
	check(err)

	walls, start, end, n := parseInput(data)
	fmt.Printf("Input: %d walls, %v start, %v end, n = %d\n", len(walls), start, end, n)

	minPath, scores := runDijkstra(walls, end, start)
	fmt.Printf("Min path is %d\n", minPath)

	result := runDijkstraWithCheats(walls, end, start, n, minPath, scores)
	fmt.Printf("Result: %d", result)
}

func runDijkstra(walls map[[2]int]struct{}, end [2]int, start [2]int) (minPath int, scores map[[2]int]int) {
	priorityQueue := make(map[int][][2]int, 1)
	priorityQueue[0] = [][2]int{start}
	scoreWatermark := 0

	scores = make(map[[2]int]int)

	for len(priorityQueue) > 0 {
		current := pop(priorityQueue, &scoreWatermark)

		if _, isVisited := scores[current]; isVisited {
			continue
		}
		scores[current] = scoreWatermark

		if current == end {
			minPath = scoreWatermark
			return
		} else {
			for _, next := range getMoves(current) {
				if _, isWall := walls[next]; !isWall {
					queueUp(priorityQueue, scoreWatermark+1, next)
				}
			}
		}
	}

	panic("Min path not found!")
}

func runDijkstraWithCheats(walls map[[2]int]struct{}, end [2]int, start [2]int, n int, minPath int, scores map[[2]int]int) (cheatCount int) {
	cheatMap := make(map[int]int, 0)

	priorityQueue := make(map[int][][2]int, 1)
	priorityQueue[0] = [][2]int{start}
	scoreWatermark := 0

	visited := make(map[[2]int]struct{})

	for len(priorityQueue) > 0 && scoreWatermark <= minPath {
		current := pop(priorityQueue, &scoreWatermark)

		if _, isVisited := visited[current]; isVisited {
			continue
		}
		visited[current] = struct{}{}
		// fmt.Printf("Processing state %v with score %d\n", current, scoreWatermark)

		if current != end {
			for _, cheat := range getCheats(current, n) {
				if _, isVisited := visited[cheat.end]; !isVisited {
					for _, next := range getMoves(cheat.end) {
						if _, isWall := walls[next]; !isWall {
							savings := scores[next] - scoreWatermark - 1 - cheat.len
							if savings >= MIN_SAVINGS {
								cheatMap[savings]++
							}
						}
					}
				}
			}

			for _, next := range getMoves(current) {
				if _, isWall := walls[next]; !isWall {
					queueUp(priorityQueue, scoreWatermark+1, next)
				}
			}
		}
	}

	for _, count := range cheatMap {
		cheatCount += count
	}

	for _, score := range slices.Sorted(maps.Keys(cheatMap)) {
		fmt.Printf("Got %d cheats that save %d\n", cheatMap[score], score)
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

func queueUp(priorityQueue map[int][][2]int, score int, item [2]int) {
	subQueue, ok := priorityQueue[score]
	if !ok {
		priorityQueue[score] = [][2]int{item}
	} else {
		priorityQueue[score] = append(subQueue, item)
	}
}

func getMoves(pos [2]int) [4][2]int {
	return [4][2]int{
		{pos[0] - 1, pos[1]},
		{pos[0] + 1, pos[1]},
		{pos[0], pos[1] - 1},
		{pos[0], pos[1] + 1},
	}
}

func getCheats(pos [2]int, n int) (cheats []cheat) {
	cheats = make([]cheat, 0)
	for x := 0; x <= CHEAT_LEN; x++ {
		for y := 0; y <= CHEAT_LEN-x; y++ {
			if x+y > 0 {
				appendIfInBounds(&cheats, n, pos[0]+x, pos[1]+y, x+y)
				appendIfInBounds(&cheats, n, pos[0]+x, pos[1]-y, x+y)
				appendIfInBounds(&cheats, n, pos[0]-x, pos[1]+y, x+y)
				appendIfInBounds(&cheats, n, pos[0]-x, pos[1]-y, x+y)
			}
		}
	}
	return
}

func appendIfInBounds(cheats *[]cheat, n int, x int, y int, d int) {
	if x > 0 || x < n-1 || y > 0 || y < n-1 {
		*cheats = append(*cheats, cheat{[2]int{x, y}, d})
	}
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