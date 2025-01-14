package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

func main() {
	// data, err := os.ReadFile("test.txt")
	data, err := os.ReadFile("input.txt")
	check(err)

	edges := parseInput(data)
	fmt.Printf("Input: %d edges\n", len(edges))

	vertexMap := make(map[string]map[string]struct{})
	for _, edge := range edges {
		appendConnection(vertexMap, edge[0], edge[1])
		appendConnection(vertexMap, edge[1], edge[0])
	}
	fmt.Printf("Got %d vertexes\n", len(vertexMap))

	seen := make(map[[3]string]struct{})
	cycles := make(map[[3]string]struct{})
	for _, edge := range edges {
		for third := range vertexMap[edge[0]] {
			canonical := getCanonical(edge, third)
			if _, wasSeen := seen[canonical]; !wasSeen {
				seen[canonical] = struct{}{}

				if containsT(canonical) {
					if _, hasCycle := vertexMap[edge[1]][third]; hasCycle {
						cycles[canonical] = struct{}{}
					}
				}
			}
		}
	}

	fmt.Printf("Result: %d\n", len(cycles))
}

func getCanonical(edge [2]string, node string) [3]string {
	temp := []string{edge[0], edge[1], node}
	slices.Sort(temp)
	return [3]string{temp[0], temp[1], temp[2]}
}

func containsT(cycle [3]string) (result bool) {
	for k := 0; k < 3; k++ {
		result = result || strings.HasPrefix(cycle[k], "t")
	}
	return
}

func appendConnection(vertexMap map[string]map[string]struct{}, from string, to string) {
	if _, ok := vertexMap[from]; !ok {
		vertexMap[from] = make(map[string]struct{})
	}

	vertexMap[from][to] = struct{}{}
}

func parseInput(data []uint8) (edges [][2]string) {
	lines := strings.Split(string(data), "\n")
	edges = make([][2]string, len(lines))

	re := regexp.MustCompile(`(\w\w)\-(\w\w)`)

	for i, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) != 3 {
			panic("Invalid regexp matches")
		}

		edges[i] = [2]string{matches[1], matches[2]}
	}
	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
