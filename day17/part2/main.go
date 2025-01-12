package main

import (
	"fmt"
	"math"
	"os"

	// "reflect"
	// "math"
	// "slices"
	// "iter"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	_, instructions := parseInput(data)
	fmt.Printf("Input: %d instructions\n", len(instructions))

	// fmt.Println(exec(117440, instructions))

	a := 0
	for ; ; a++ {
		// if a%1000 == 0 {
		// 	fmt.Println(a)
		// }
		if success := exec(a, instructions); success {
			break
		}

		if a == math.MaxInt {
			panic("Reached max int")
		}
	}

	fmt.Printf("Result: %d\n", a)
}

func exec(a int, instructions []int) bool {
	registers := [3]int{a, 0, 0}
	pointer := 0
	out := make([]int, 0)

	for pointer < len(instructions) {
		opcode := instructions[pointer]
		operand := instructions[pointer+1]
		pointer += 2

		switch opcode {
		case 0: // adv
			registers[0] = registers[0] / (1 << combo(&registers, operand))
		case 1: //bxl
			registers[1] = registers[1] ^ operand
		case 2: // bst
			registers[1] = combo(&registers, operand) % 8
		case 3: // jnz
			if registers[0] != 0 {
				pointer = operand
			}
		case 4: //bxc
			registers[1] = registers[1] ^ registers[2]
		case 5: // out
			i := len(out)
			value := combo(&registers, operand) % 8
			if i >= len(instructions) || instructions[i] != value {
				return false
			}
			out = append(out, value)
		case 6: // bdv
			registers[1] = registers[0] / (1 << combo(&registers, operand))
		case 7: // cdv
			registers[2] = registers[0] / (1 << combo(&registers, operand))
		default:
			panic("Unknown opcode")
		}
	}

	return len(out) == len(instructions)
}

func combo(registers *[3]int, operand int) int {
	if operand < 0 || operand >= 7 {
		panic("Invalid combo operand")
	}

	if operand <= 3 {
		return operand
	}

	return registers[operand-4]
}

func parseInput(data []uint8) (registers [3]int, instructions []int) {
	sections := strings.Split(string(data), "\n\n")
	if len(sections) != 2 {
		panic("Invalid number of input sections")
	}

	re := regexp.MustCompile(`Register A: (\d+)
Register B: (\d+)
Register C: (\d+)`)

	registerMatches := re.FindStringSubmatch(sections[0])
	if len(registerMatches) != 4 {
		panic("Invalid register matches")
	}
	registers[0] = parseInt(registerMatches[1])
	registers[1] = parseInt(registerMatches[2])
	registers[2] = parseInt(registerMatches[3])

	chunks := strings.Split(strings.TrimPrefix(sections[1], "Program: "), ",")
	instructions = make([]int, len(chunks))
	for i, chunk := range chunks {
		instructions[i] = parseInt(chunk)
	}
	return
}

func parseInt(s string) int {
	result, e := strconv.ParseInt(s, 0, 64)
	check(e)
	return int(result)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
