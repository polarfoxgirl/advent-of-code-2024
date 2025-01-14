package main

import (
	"fmt"
	"maps"
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

	// Start with clusters of 2 (edges)
	knownClusters := make(map[string]map[string]struct{})
	for _, edge := range edges {
		cluster := map[string]struct{}{edge[0]: {}, edge[1]: {}}
		pwd := getPassword(cluster)
		knownClusters[pwd] = cluster
	}

	// Try to increase cluster rank
	for {
		nextRankClusters := make(map[string]map[string]struct{})
		for _, cluster := range knownClusters {
			seen := make(map[string]struct{})
			for clusterNode := range cluster {
				for newNode := range vertexMap[clusterNode] {
					if _, isSeen := seen[newNode]; !isSeen {
						seen[newNode] = struct{}{}
						if newCluster, ok := expandCluster(vertexMap, newNode, cluster); ok {
							pwd := getPassword(newCluster)
							nextRankClusters[pwd] = newCluster
						}
					}
				}
			}
		}

		if len(nextRankClusters) == 0 {
			for pwd, _ := range knownClusters {
				fmt.Println(pwd)
			}
			break
		} else {
			knownClusters = nextRankClusters
		}
	}
}

func expandCluster(vertexMap map[string]map[string]struct{}, newNode string, cluster map[string]struct{}) (result map[string]struct{}, ok bool) {
	if _, inCluster := cluster[newNode]; inCluster {
		return
	}

	connections, _ := vertexMap[newNode]
	for node := range cluster {
		if _, connected := connections[node]; !connected {
			return
		}
	}

	result = map[string]struct{}{newNode: {}}
	maps.Copy(result, cluster)
	ok = true
	return
}

func getPassword(cluster map[string]struct{}) string {
	temp := slices.Sorted(maps.Keys(cluster))
	return strings.Join(temp, ",")
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
