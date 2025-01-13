package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	// "slices"
	// "iter"
	// "regexp"
	// "strconv"
	"strings"
)

type color string

type node struct {
	children map[color]node
	prefix   string
}

const END = "."

type state struct {
	pos    int
	prefix string
}

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	towels, designs := parseInput(data)
	fmt.Printf("Input: %d towels, %d designs\n", len(towels), len(designs))

	towelRoot, _ := getTree(towels)
	designRoot, designIndex := getTree(designs)

	count := runDijkstra(designRoot, designIndex, towelRoot)

	fmt.Printf("Result: %d\n", count)
}

func getTree(towels [][]color) (root node, index map[string]node) {
	root.children = make(map[color]node)
	root.prefix = ""

	index = map[string]node{root.prefix: root}

	for _, towel := range towels {
		pointer := root
		for _, stripe := range towel {
			if _, ok := pointer.children[stripe]; !ok {
				prefix := pointer.prefix + string(stripe)
				childNode := node{make(map[color]node), prefix}
				pointer.children[stripe] = childNode
				index[prefix] = childNode
			}
			pointer, _ = pointer.children[stripe]
		}
		pointer.children[END] = node{make(map[color]node), pointer.prefix + END}
	}
	return
}

func runDijkstra(designRoot node, designIndex map[string]node, towelRoot node) (count int) {
	paths := map[string]int{"": 1}

	priorityQueue := make(map[int][]string)
	priorityQueue[0] = []string{designRoot.prefix}
	visited := make(map[string]struct{})
	scoreWatermark := 0

	for len(priorityQueue) > 0 {
		designNodePrefix := pop(priorityQueue, &scoreWatermark)
		designNode, ok := designIndex[designNodePrefix]
		if !ok {
			panic(fmt.Sprintf("Unknown prefix %s", designNodePrefix))
		}

		if _, isVisited := visited[designNodePrefix]; isVisited {
			continue
		}
		visited[designNodePrefix] = struct{}{}

		nodePaths, _ := paths[designNodePrefix]
		if _, isDesignEnd := designNode.children[END]; isDesignEnd {
			count += nodePaths
		}

		queueChild := func(prefix string) {
			paths[prefix] += nodePaths
			queueUp(priorityQueue, prefix)
		}

		queueChildren(designNode, towelRoot, queueChild)
	}

	return
}

func queueChildren(designChild node, towelChild node, queueChild func(string)) {
	if _, hasTowelEnd := towelChild.children[END]; hasTowelEnd {
		queueChild(designChild.prefix)
	}

	for stripe, designNext := range designChild.children {
		if stripe == END {
			continue
		}

		if towelNext, canContinueTowel := towelChild.children[stripe]; canContinueTowel {
			queueChildren(designNext, towelNext, queueChild)
		}
	}
}

func pop(priorityQueue map[int][]string, scoreWatermark *int) (result string) {
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

func queueUp(priorityQueue map[int][]string, prefix string) {
	subQueue, ok := priorityQueue[len(prefix)]
	if !ok {
		priorityQueue[len(prefix)] = []string{prefix}
	} else {
		priorityQueue[len(prefix)] = append(subQueue, prefix)
	}
}

func parseInput(data []uint8) (towels [][]color, designs [][]color) {
	sections := strings.Split(string(data), "\n\n")
	if len(sections) != 2 {
		panic("Invalid number of input sections")
	}

	towelStrs := strings.Split(sections[0], ", ")
	towels = make([][]color, len(towelStrs))
	for i, towelStr := range towelStrs {
		towels[i] = parseStripes(towelStr)
	}

	designStrs := strings.Split(sections[1], "\n")
	designs = make([][]color, len(designStrs))
	for i, designStr := range designStrs {
		designs[i] = parseStripes(designStr)
	}

	return
}

func parseStripes(s string) (stripes []color) {
	stripes = make([]color, 0)
	for _, rune := range s {
		stripes = append(stripes, parseColor(rune))
	}
	return
}

func parseColor(r rune) color {
	switch r {
	case 'w':
		fallthrough
	case 'u':
		fallthrough
	case 'b':
		fallthrough
	case 'r':
		fallthrough
	case 'g':
		return color(r)
	default:
		panic("Unknown color!")
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
