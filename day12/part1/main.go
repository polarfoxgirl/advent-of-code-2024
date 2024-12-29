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

	perimeter := calcPerimeter(region)
	return len(region) * perimeter
}

func calcPerimeter(region map[[2]int]struct{}) (perimeter int) {
	for plot := range region {
		if _, ok := region[[2]int{plot[0] - 1, plot[1]}]; !ok {
			perimeter++
		}
		if _, ok := region[[2]int{plot[0] + 1, plot[1]}]; !ok {
			perimeter++
		}
		if _, ok := region[[2]int{plot[0], plot[1] - 1}]; !ok {
			perimeter++
		}
		if _, ok := region[[2]int{plot[0], plot[1] + 1}]; !ok {
			perimeter++
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
