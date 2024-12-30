package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	"slices"
	// "strconv"
	// "iter"
	"strings"
)

type BoxMove struct {
	from [2]int
	to   [2]int
}

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
				robot = [2]int{i, 2 * j}
			case '#':
				walls[[2]int{i, 2 * j}] = struct{}{}
				walls[[2]int{i, 2*j + 1}] = struct{}{}
			case 'O':
				boxes[[2]int{i, 2 * j}] = struct{}{}
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

	if isBox, box := isBox(boxes, target); isBox {
		ok, boxStack := getBoxMovePlan(walls, boxes, box, move)
		if !ok {
			// Can't move
			return
		}

		// Move all boxes
		for _, shift := range boxStack {
			delete(boxes, shift.from)
		}
		for _, shift := range boxStack {
			boxes[shift.to] = struct{}{}
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

func isBox(boxes map[[2]int]struct{}, pos [2]int) (ok bool, box [2]int) {
	if _, hasBox := boxes[pos]; hasBox {
		ok = true
		box = pos
		return
	}

	shifted := [2]int{pos[0], pos[1] - 1}
	if _, hasBox := boxes[shifted]; hasBox {
		ok = true
		box = shifted
		return
	}

	ok = false
	return
}

func getBoxMovePlan(walls map[[2]int]struct{}, boxes map[[2]int]struct{}, box [2]int, move rune) (ok bool, boxStack []BoxMove) {
	boxStack = make([]BoxMove, 0)
	nextLeft := doMove(box, move)
	nextRight := [2]int{nextLeft[0], nextLeft[1] + 1}

	if _, isWall := walls[nextLeft]; isWall {
		ok = false
		return
	}
	if _, isWall := walls[nextRight]; isWall {
		ok = false
		return
	}

	checkLeft, boxLeft := isBox(boxes, nextLeft)
	if checkLeft && boxLeft != box {
		canMove, boxStackLeft := getBoxMovePlan(walls, boxes, boxLeft, move)
		if !canMove {
			ok = false
			return
		}

		boxStack = boxStackLeft
	}

	checkRight, boxRight := isBox(boxes, nextRight)
	if checkRight && boxRight != box && boxRight != boxLeft {
		canMove, boxStackRight := getBoxMovePlan(walls, boxes, boxRight, move)
		if !canMove {
			ok = false
			return
		}

		boxStack = slices.Concat(boxStack, boxStackRight)
	}

	ok = true
	boxStack = append(boxStack, BoxMove{from: box, to: nextLeft})
	return
}

func print(walls map[[2]int]struct{}, boxes map[[2]int]struct{}, robot [2]int, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < 2*n; j++ {
			pos := [2]int{i, j}
			if pos == robot {
				fmt.Print("@")
			} else if _, isWall := walls[pos]; isWall {
				fmt.Print("#")
			} else if _, isBox := boxes[pos]; isBox {
				fmt.Print("[]")
				j++
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
