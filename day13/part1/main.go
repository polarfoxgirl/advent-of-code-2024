package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	// "slices"
	// "iter"
	"regexp"
	"strconv"
	"strings"
)

type ClawMachine struct {
	a     [2]int
	b     [2]int
	prize [2]int
}

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	machines := parseInput(data)
	fmt.Printf("Input: %d machines\n", len(machines))

	sum := 0
	for _, machine := range machines {
		sum += getMinTokens(machine)
	}

	fmt.Printf("Result: %d", sum)
}

func parseInput(data []uint8) (machines []ClawMachine) {
	sections := strings.Split(string(data), "\n\n")
	machines = make([]ClawMachine, len(sections))

	re := regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)
Button B: X\+(\d+), Y\+(\d+)
Prize: X=(\d+), Y=(\d+)`)

	for i, section := range sections {
		matches := re.FindStringSubmatch(section)

		machines[i] = ClawMachine{
			a:     [2]int{parseInt(matches[1]), parseInt(matches[2])},
			b:     [2]int{parseInt(matches[3]), parseInt(matches[4])},
			prize: [2]int{parseInt(matches[5]), parseInt(matches[6])},
		}
	}
	return
}

func getMinTokens(cm ClawMachine) (result int) {
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			if cm.a[0]*i+cm.b[0]*j == cm.prize[0] && cm.a[1]*i+cm.b[1]*j == cm.prize[1] {
				price := 3*i + j
				if result == 0 || price < result {
					result = price
				}
			}
		}
	}
	return
}

func parseInt(s string) int {
	result, e := strconv.ParseInt(s, 0, 64)
	check(e)
	return int(result)
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
