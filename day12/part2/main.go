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

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	plots := parseInput(data)
	fmt.Printf("Input: %d x %d plots\n", len(plots), len(plots))

	visited := make(map[[2]int]struct{})
	sum := 0
	for i := 0; i < len(plots); i++ {
		for j := 0; j < len(plots); j++ {
			if _, isVisited := visited[[2]int{i, j}]; !isVisited {
				sum += crawlRegion(plots, i, j, visited)
			}
		}
	}

	fmt.Printf("Result: %d", sum)
}

func parseInput(data []uint8) (plots [][]rune) {
	lines := strings.Split(string(data), "\n")
	plots = make([][]rune, len(lines))

	for i, line := range lines {
		plots[i] = make([]rune, len(lines))
		for j, rune := range line {
			plots[i][j] = rune
		}
	}
	return
}

func crawlRegion(plots [][]rune, x int, y int, visited map[[2]int]struct{}) int {
	crop := plots[x][y]
	queue := [][2]int{{x, y}}
	region := make(map[[2]int]struct{})
	for len(queue) > 0 {
		plot := queue[0]
		queue = queue[1:]

		if plots[plot[0]][plot[1]] != crop {
			continue
		}
		if _, isVisited := visited[plot]; isVisited {
			continue
		}

		visited[plot] = struct{}{}
		region[plot] = struct{}{}

		if plot[0] > 0 {
			queue = append(queue, [2]int{plot[0] - 1, plot[1]})
		}
		if plot[0] < len(plots)-1 {
			queue = append(queue, [2]int{plot[0] + 1, plot[1]})
		}
		if plot[1] > 0 {
			queue = append(queue, [2]int{plot[0], plot[1] - 1})
		}
		if plot[1] < len(plots)-1 {
			queue = append(queue, [2]int{plot[0], plot[1] + 1})
		}
	}

	sides := countSides(region)
	return len(region) * sides
}

func countSides(region map[[2]int]struct{}) (sides int) {
	left := make(map[int][]int)
	right := make(map[int][]int)
	up := make(map[int][]int)
	down := make(map[int][]int)

	for plot := range region {
		if _, ok := region[[2]int{plot[0] - 1, plot[1]}]; !ok {
			if _, ok := left[plot[0]]; !ok {
				left[plot[0]] = make([]int, 0)
			}
			left[plot[0]] = append(left[plot[0]], plot[1])
		}
		if _, ok := region[[2]int{plot[0] + 1, plot[1]}]; !ok {
			if _, ok := right[plot[0]]; !ok {
				right[plot[0]] = make([]int, 0)
			}
			right[plot[0]] = append(right[plot[0]], plot[1])
		}
		if _, ok := region[[2]int{plot[0], plot[1] - 1}]; !ok {
			if _, ok := up[plot[1]]; !ok {
				up[plot[1]] = make([]int, 0)
			}
			up[plot[1]] = append(up[plot[1]], plot[0])
		}
		if _, ok := region[[2]int{plot[0], plot[1] + 1}]; !ok {
			if _, ok := down[plot[1]]; !ok {
				down[plot[1]] = make([]int, 0)
			}
			down[plot[1]] = append(down[plot[1]], plot[0])
		}
	}

	return countDistinctSides(left) + countDistinctSides(right) + countDistinctSides(up) + countDistinctSides(down)
}

func countDistinctSides(borders map[int][]int) (sides int) {
	for _, border := range borders {
		slices.Sort(border)

		sides++
		for i := 0; i < len(border)-1; i++ {
			if border[i]+1 < border[i+1] {
				sides++
			}
		}
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
