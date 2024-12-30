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

	walls, boxes, robot, moves, n := parseInput(data)
	fmt.Printf("Input: %d walls, %d boxes, %v start, %d moves\n", len(walls), len(boxes), robot, len(moves))

	for _, move := range moves {
		robot = tryMove(walls, boxes, robot, move)
	}
	print(walls, boxes, robot, n)

	sum := 0
	for box := range boxes {
		sum += 100*box[0] + box[1]
	}

	fmt.Printf("Result: %d", sum)
}

func parseInput(data []uint8) (walls map[[2]int]struct{}, boxes map[[2]int]struct{}, robot [2]int, moves []rune, n int) {
	sections := strings.Split(string(data), "\n\n")
	if len(sections) != 2 {
		panic("Unexpected number of input sections")
	}

	walls = make(map[[2]int]struct{})
	boxes = make(map[[2]int]struct{})
	for i, line := range strings.Split(string(sections[0]), "\n") {
		n = len(line)
		for j, rune := range line {
			switch rune {
			case '@':
				robot = [2]int{i, j}
			case '#':
				walls[[2]int{i, j}] = struct{}{}
			case 'O':
				boxes[[2]int{i, j}] = struct{}{}
			case '.':
				// no-op
			default:
				panic("Unexpected rune on the map: " + string(rune))
			}
		}
	}

	moves = make([]rune, 0)
	for _, rune := range string(sections[1]) {
		if rune != '\n' {
			moves = append(moves, rune)
		}
	}
	return
}

func tryMove(walls map[[2]int]struct{}, boxes map[[2]int]struct{}, current [2]int, move rune) (next [2]int) {
	next = current

	target := doMove(current, move)
	if _, isWall := walls[target]; isWall {
		// Can't move
		return
	}

	if _, isBox := boxes[target]; isBox {
		pos := target
		for {
			pos = doMove(pos, move)
			if _, isWall := walls[pos]; isWall {
				// Boxes all the way to a wall
				return
			}
			if _, stillBox := boxes[pos]; stillBox {
				// More boxes
				continue
			}

			// Found an empty space! Shift boxes here
			boxes[pos] = struct{}{}
			delete(boxes, target)
			break
		}
	}

	// Move robot
	next = target
	return
}

func doMove(pos [2]int, move rune) [2]int {
	switch move {
	case '<':
		return [2]int{pos[0], pos[1] - 1}
	case '>':
		return [2]int{pos[0], pos[1] + 1}
	case '^':
		return [2]int{pos[0] - 1, pos[1]}
	case 'v':
		return [2]int{pos[0] + 1, pos[1]}
	default:
		panic("Unexpected move: " + string(move))
	}
}

func print(walls map[[2]int]struct{}, boxes map[[2]int]struct{}, robot [2]int, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			pos := [2]int{i, j}
			if pos == robot {
				fmt.Print("@")
			} else if _, isWall := walls[pos]; isWall {
				fmt.Print("#")
			} else if _, isBox := boxes[pos]; isBox {
				fmt.Print("O")
			} else {
				fmt.Print(".")

			}
		}
		fmt.Println()
	}
	fmt.Println()
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
