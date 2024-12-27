package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	"slices"
	// "iter"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	rules, updates := parseInput(data)
	fmt.Printf("Input: %d rules, %d updates\n", len(rules), len(updates))

	ruleMap := make(map[int]map[int]struct{})
	for _, rule := range rules {
		if _, ok := ruleMap[rule[0]]; !ok {
			ruleMap[rule[0]] = make(map[int]struct{})
		}
		ruleMap[rule[0]][rule[1]] = struct{}{}
	}

	compFunc := getCompFunc(ruleMap)
	var sum int
	for _, update := range updates {
		if !isInOrder(ruleMap, update) {
			slices.SortStableFunc(update, compFunc)
			sum += update[len(update)/2]
		}
	}

	fmt.Printf("Result: %d", sum)
}

func parseInput(data []uint8) (rules [][2]int, updates [][]int) {
	sections := strings.Split(string(data), "\n\n")
	if len(sections) > 2 {
		panic("Too many sections")
	}

	ruleLines := strings.Split(sections[0], "\n")
	rules = make([][2]int, len(ruleLines))

	for i, line := range ruleLines {
		pair := strings.Split(line, "|")
		if len(pair) != 2 {
			panic("Invalid rule")
		}

		n1, e1 := strconv.ParseInt(pair[0], 0, 32)
		check(e1)
		n2, e2 := strconv.ParseInt(pair[1], 0, 32)
		check(e2)
		rules[i] = [2]int{int(n1), int(n2)}
	}

	updateLines := strings.Split(sections[1], "\n")
	updates = make([][]int, len(updateLines))

	for i, line := range updateLines {
		numbers := strings.Split(line, ",")
		updates[i] = make([]int, len(numbers))
		for j, number := range numbers {
			n, e := strconv.ParseInt(number, 0, 32)
			check(e)
			updates[i][j] = int(n)
		}
	}
	return
}

func isInOrder(ruleMap map[int]map[int]struct{}, update []int) bool {
	seen := make(map[int]struct{})
	for _, n := range update {
		followers, hasRules := ruleMap[n]
		if hasRules {
			for follower := range followers {
				if _, inSeen := seen[follower]; inSeen {
					return false
				}
			}
		}

		seen[n] = struct{}{}
	}
	return true
}

func getCompFunc(ruleMap map[int]map[int]struct{}) func(int, int) int {
	return func(x int, y int) int {
		xFollowers, hasXRules := ruleMap[x]
		if hasXRules {
			if _, hasY := xFollowers[y]; hasY {
				return -1
			}
		}

		yFollowers, hasYRules := ruleMap[y]
		if hasYRules {
			if _, hasX := yFollowers[y]; hasX {
				return 1
			}
		}

		return 0
	}
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
