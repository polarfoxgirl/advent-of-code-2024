package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
Context: looking for cycles in numbers that produce prefixes

Step 1: at least 7

Register 7171773 produced [2 4 1 1 7 5 1]
256
Register 7172029 produced [2 4 1 1 7 5 1]
Register 13793717 produced [2 4 1 1 7 5 1]
Register 15560381 produced [2 4 1 1 7 5 1]
Register 15560637 produced [2 4 1 1 7 5 1]
Register 32337597 produced [2 4 1 1 7 5 1]
Register 32337853 produced [2 4 1 1 7 5 1]
... pairs down the line

Step 2: at least 9, start 7171773, step 256
Register 158166717 produced [2 4 1 1 7 5 1 5 4]
134217728
Register 292384445 produced [2 4 1 1 7 5 1 5 4]
134217728
Register 426602173 produced [2 4 1 1 7 5 1 5 4]
Register 560819901 produced [2 4 1 1 7 5 1 5 4]
Register 695037629 produced [2 4 1 1 7 5 1 5 4]
Register 963473085 produced [2 4 1 1 7 5 1 5 4]
...

Step 3: start 158166717, step 134217728
PROFIT!!
*/
func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	_, instructions := parseInput(data)
	fmt.Printf("Input: %d instructions\n", len(instructions))
	fmt.Println(instructions)

	// a := 7171773
	// for ; ; a += 256 {
	a := 158166717
	for ; ; a += 134217728 {
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
				// if len(out) > 9 {
				// 	fmt.Printf("Register %d produced %v\n", a, out)
				// }
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
