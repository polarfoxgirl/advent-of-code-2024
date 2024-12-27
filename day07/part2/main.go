package main

import (
	"fmt"
	"os"
	// "reflect"
	"math"
	// "slices"
	// "iter"
	"strconv"
	"strings"
)

type Equation struct {
	result int64
	args   []int64
}

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	equations := parseInput(data)
	fmt.Printf("Input: %d equations\n", len(equations))

	sum := 0
	for _, equation := range equations {
		if isSolvable(equation.result, equation.args[0], equation.args[1:]) {
			sum += int(equation.result)
		}
	}

	fmt.Printf("Result: %d", sum)
}

func parseInput(data []uint8) (equations []Equation) {
	lines := strings.Split(string(data), "\n")
	equations = make([]Equation, len(lines))

	for i, line := range lines {
		pair := strings.Split(line, ": ")
		if len(pair) != 2 {
			panic("Invalid equation")
		}

		result, e := strconv.ParseInt(pair[0], 0, 64)
		check(e)

		argStrs := strings.Split(pair[1], " ")
		equation := Equation{result: result, args: make([]int64, len(argStrs))}
		for j, argStr := range argStrs {
			arg, e := strconv.ParseInt(argStr, 0, 64)
			check(e)
			equation.args[j] = arg
		}

		equations[i] = equation
	}
	return
}

func isSolvable(target int64, left int64, args []int64) bool {
	if len(args) == 0 {
		return target == left
	}

	return isSolvable(target, left+args[0], args[1:]) || isSolvable(target, left*args[0], args[1:]) || isSolvable(target, concat(left, args[0]), args[1:])
}

func concat(x int64, y int64) int64 {
	shift := int(math.Log10(float64(y))) + 1
	return x*int64(math.Pow10(shift)) + y
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
