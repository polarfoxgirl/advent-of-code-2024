package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
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
	data, err := os.ReadFile("input_fixed.txt")
	check(err)

	values, gates := parseInput(data)
	fmt.Printf("Input: %d inputs, %d gates\n", len(values), len(gates))

	inputAnds := make(map[string]string)
	inputXors := make(map[string]string)
	for _, gate := range gates {
		if strings.HasPrefix(gate.left, "x") || strings.HasPrefix(gate.left, "y") {
			code := gate.left[1:]
			if gate.op == AND {
				inputAnds[code] = gate.out
			} else if gate.op == XOR {
				inputXors[code] = gate.out
			} else {
				panic(fmt.Sprintf("Invalid input gate: %v", gate))
			}
		}
	}

	// 45 inputs from x00 to x44
	if len(inputAnds) != 45 || len(inputAnds) != len(inputXors) {
		fmt.Printf("Incomplete input ANDs (%d) or input XORs (%d)\n", len(inputAnds), len(inputXors))
	}

	if z00Out, ok := inputXors["00"]; !ok || z00Out != "z00" {
		fmt.Println("Missing z00 gate")
	}

	orOuts := make(map[string]string)
	andOuts := make(map[string]string)
	for code, inputXor := range inputXors {
		// x00 AND y00 -> hmc
		// x00 XOR y00 -> z00

		if code != "00" {
			for _, gate := range gates {
				if gate.left == inputXor || gate.right == inputXor {
					if gate.op == XOR {
						if strings.HasPrefix(gate.out, "z") {
							if gate.out[1:] != code {
								fmt.Printf("Invalid z gate %v (for code %s and input %s)\n", gate, code, inputXor)
							}

							prevCode := fmt.Sprintf("%d", parseInt(code)-1)
							if len(prevCode) == 1 {
								prevCode = "0" + prevCode
							}
							if gate.left == inputXor {
								orOuts[prevCode] = gate.right
							} else {
								orOuts[prevCode] = gate.left
							}
						} else {
							// Gate {nfh nqq 2 qgd} should be a z-gate for code 18
							fmt.Printf("Gate %v should be a z-gate for code %s\n", gate, code)
						}
					} else if gate.op == AND {
						if strings.HasPrefix(gate.out, "z") {
							// Invalid shift AND gate {cvt wwp 0 z33} (for code 33 and input wwp)
							// Swap: z33 and ?gqp
							// wwp XOR cvt -> gqp
							// cvt AND wwp -> z33
							fmt.Printf("Invalid shift AND gate %v (for code %s and input %s)\n", gate, code, inputXor)
						} else {
							andOuts[code] = gate.out
						}
					} else {
						// Invalid OR gate {bmv hsw 1 mtn} (for code 24 and input hsw)
						// Swap: hsw and ?jmh
						// y24 XOR x24 -> hsw
						// y24 AND x24 -> jmh
						// jmh AND mrs -> bmv
						// jmh XOR mrs -> z24
						fmt.Printf("Invalid OR gate %v (for code %s and input %s)\n", gate, code, inputXor)
					}
				}
			}
		}
	}

	for code, inputAnd := range inputAnds {
		// Technically need to validate z00 transfer, but it's fine
		if code != "00" {
			if andOut, ok := andOuts[code]; ok {
				for _, gate := range gates {
					if gate.left == inputAnd || gate.right == inputAnd {
						if gate.left != andOut && gate.right != andOut {
							fmt.Printf("Gate %v has inputAnd %s but not andOut %s for code %s\n", gate, inputAnd, andOut, code)
						}

						if gate.op != OR {
							fmt.Printf("Non-OR gate %v for code %s\n", gate, code)
						}

						// Gate {kgd kqf 1 z10} doesn't have expected out mwk for code 10
						// Swap: z10 and ?mwk
						// y10 XOR x10 -> sqr
						// mwq AND sqr -> kgd
						// sqr XOR mwq -> mwk
						//
						// Gate {ndb qqc 1 nfh} doesn't have expected out  for code 17
						// x17 XOR y17 -> thq +
						// x17 AND y17 -> qqc +
						// mfd AND thq -> ndb +
						// thq XOR mfd -> z17 +
						// ndb OR qqc -> nfh ...
						// nfh XOR nqq -> qgd ?
						// nfh AND nqq -> ggn +
						// y18 AND x18 -> z18 ?
						// qgd OR ggn -> cvr
						// Swap z18 and qgd
						if code != "44" {
							if orOut, known := orOuts[code]; !known || orOut != gate.out {
								fmt.Printf("Gate %v doesn't have expected out %s for code %s\n", gate, orOut, code)
							}
						}
					}
				}
			} else {
				fmt.Printf("Missing and out for code %s\n", code)
			}
		}
	}

	// Technically need to validate z45 gate, but it's fine

	unprocessed := make(map[string]struct{})
	for wire := range values {
		unprocessed[wire] = struct{}{}
	}
	for len(unprocessed) > 0 {
		usedGates := make([]gate, 0)
		nextWave := make(map[string]struct{})
		for _, gate := range gates {
			if _, knownOut := values[gate.out]; !knownOut {
				if _, leftUnprocessed := unprocessed[gate.left]; leftUnprocessed {
					if rightValue, rightPresent := values[gate.right]; rightPresent {
						if _, inNextWave := nextWave[gate.right]; !inNextWave {
							values[gate.out] = exec(values[gate.left], rightValue, gate.op)
							nextWave[gate.out] = struct{}{}
							usedGates = append(usedGates, gate)
						}
					}
				} else if _, rightUnprocessed := unprocessed[gate.right]; rightUnprocessed {
					if leftValue, leftPresent := values[gate.left]; leftPresent {
						if _, inNextWave := nextWave[gate.left]; !inNextWave {
							values[gate.out] = exec(leftValue, values[gate.right], gate.op)
							nextWave[gate.out] = struct{}{}
							usedGates = append(usedGates, gate)
						}
					}
				}
			}
		}

		// fmt.Println(usedGates)
		unprocessed = nextWave
	}

	fmt.Printf("Desired: %b\nActual:: %b\n", getResult(values, "z"), getResult(values, "x")+getResult(values, "y"))

	answer := []string{"z33", "gqp", "hsw", "jmh", "z10", "mwk", "z18", "qgd"}
	slices.Sort(answer)
	fmt.Println(strings.Join(answer, ","))
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

func isValidSum(values map[string]bool) bool {
	return getResult(values, "z") == getResult(values, "x")+getResult(values, "y")
}

func getResult(values map[string]bool, prefix string) (result int) {
	re := regexp.MustCompile(prefix + `(\d+)`)
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
