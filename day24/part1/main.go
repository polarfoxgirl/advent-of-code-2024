package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type operation int

const (
	AND operation = iota
	OR
	XOR
)

type gate struct {
	left  string
	right string
	op    operation
	out   string
}

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	values, gates := parseInput(data)
	fmt.Printf("Input: %d inputs, %d gates\n", len(values), len(gates))

	unprocessed := make(map[string]struct{})
	for wire, _ := range values {
		unprocessed[wire] = struct{}{}
	}
	for len(unprocessed) > 0 {
		nextWave := make(map[string]struct{})
		for _, gate := range gates {
			if _, knownOut := values[gate.out]; !knownOut {
				if _, leftUnprocessed := unprocessed[gate.left]; leftUnprocessed {
					if rightValue, rightPresent := values[gate.right]; rightPresent {
						values[gate.out] = exec(values[gate.left], rightValue, gate.op)
						nextWave[gate.out] = struct{}{}
					}
				} else if _, rightUnprocessed := unprocessed[gate.right]; rightUnprocessed {
					if leftValue, leftPresent := values[gate.left]; leftPresent {
						values[gate.out] = exec(leftValue, values[gate.right], gate.op)
						nextWave[gate.out] = struct{}{}
					}
				}
			}
		}

		unprocessed = nextWave
	}

	fmt.Printf("Result: %d", getResult(values))

}

func exec(left bool, right bool, op operation) bool {
	switch op {
	case AND:
		return left && right
	case OR:
		return left || right
	case XOR:
		return left != right
	}
	panic("Unknown op")
}

func getResult(values map[string]bool) (result int) {
	re := regexp.MustCompile(`z(\d+)`)
	for wire, value := range values {
		if value {
			matches := re.FindStringSubmatch(wire)
			if matches != nil {
				result += 1 << parseInt(matches[1])
			}
		}
	}
	return
}

func parseInput(data []uint8) (inputs map[string]bool, gates []gate) {
	sections := strings.Split(string(data), "\n\n")
	if len(sections) != 2 {
		panic("Invalid number of input sections")
	}

	re1 := regexp.MustCompile(`(\w+): (\d)`)
	inputs = make(map[string]bool)
	for _, line := range strings.Split(sections[0], "\n") {
		matches := re1.FindStringSubmatch(line)
		if len(matches) != 3 {
			panic("Invalid input matches")
		}

		inputs[matches[1]] = parseBool(matches[2])
	}

	re2 := regexp.MustCompile(`(\w+) (\w+) (\w+) -> (\w+)`)
	gateLines := strings.Split(sections[1], "\n")
	gates = make([]gate, len(gateLines))
	for i, line := range gateLines {
		matches := re2.FindStringSubmatch(line)
		if len(matches) != 5 {
			panic("Invalid gate matches")
		}

		gates[i] = gate{left: matches[1], op: parseOp(matches[2]), right: matches[3], out: matches[4]}
	}

	return
}

func parseBool(s string) bool {
	if s == "1" {
		return true
	}
	if s == "0" {
		return false
	}
	panic("Invalid boolean")
}

func parseOp(s string) operation {
	switch s {
	case "AND":
		return AND
	case "OR":
		return OR
	case "XOR":
		return XOR
	}
	panic("Invalid gate op")
}

func parseInt(s string) int {
	trimmed := strings.TrimLeft(s, "0")
	if len(trimmed) == 0 {
		return 0
	}
	result, e := strconv.ParseInt(trimmed, 0, 64)
	check(e)
	return int(result)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
