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

type node map[color]node

const END = "."

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	towels, designs := parseInput(data)
	fmt.Printf("Input: %d towels, %d designs\n", len(towels), len(designs))

	root := getTowelTree(towels)

	count := 0
	for _, design := range designs {
		if canMatch(design, 0, root, root) {
			count++
		}
	}

	fmt.Printf("Result: %d\n", count)
}

func getTowelTree(towels [][]color) (root node) {
	root = make(node)

	for _, towel := range towels {
		pointer := root
		for _, stripe := range towel {
			if _, ok := pointer[stripe]; !ok {
				pointer[stripe] = make(node)
			}
			pointer, _ = pointer[stripe]
		}
		pointer[END] = make(node)
	}
	return
}

func canMatch(design []color, pos int, root node, child node) bool {
	if pos >= len(design) {
		_, hasEnd := child[END]
		return hasEnd
	}
	stripe := design[pos]

	if next, canContinue := child[stripe]; canContinue {
		if canMatch(design, pos+1, root, next) {
			return true
		}
	}

	if _, canEnd := child[END]; canEnd {
		if canMatch(design, pos, root, root) {
			return true
		}
	}

	return false
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
