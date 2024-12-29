package main

import (
	"fmt"
	"math/big"
	"os"

	// "reflect"
	// "math"
	// "slices"
	// "iter"
	// "big"
	// "cmp"
	"regexp"
	"strconv"
	"strings"
)

type ClawMachine struct {
	a     [2]big.Int
	b     [2]big.Int
	prize [2]big.Int
}

// Could have used big.rat, but welp
type Fraction struct {
	num big.Int
	den big.Int
}

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	machines := parseInput(data)
	fmt.Printf("Input: %d machines\n", len(machines))

	var sum big.Int
	for _, machine := range machines {
		price := getMinTokens(machine)
		sum.Add(&sum, &price)
	}

	fmt.Printf("Result: %s", sum.String())
}

func parseInput(data []uint8) (machines []ClawMachine) {
	sections := strings.Split(string(data), "\n\n")
	machines = make([]ClawMachine, len(sections))

	re := regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)
Button B: X\+(\d+), Y\+(\d+)
Prize: X=(\d+), Y=(\d+)`)

	for i, section := range sections {
		matches := re.FindStringSubmatch(section)

		prize := [2]big.Int{parseInt(matches[5]), parseInt(matches[6])}
		prize[0].Add(&prize[0], big.NewInt(10000000000000))
		prize[1].Add(&prize[1], big.NewInt(10000000000000))

		machines[i] = ClawMachine{
			a:     [2]big.Int{parseInt(matches[1]), parseInt(matches[2])},
			b:     [2]big.Int{parseInt(matches[3]), parseInt(matches[4])},
			prize: prize,
		}
	}
	return
}

func getMinTokens(cm ClawMachine) (result big.Int) {
	mA := Fraction{cm.a[1], cm.a[0]} // Slope of vector A
	mB := Fraction{cm.b[1], cm.b[0]} // Slope of vector B

	// We assume that no machine has vertical slope
	if cm.a[0].Sign() == 0 || cm.b[0].Sign() == 0 {
		fmt.Println("Found vertical slope :(")
		return
	}

	comparison := compare(mA, mB)
	if comparison == 0 {
		// Cheating here since no machines satisfy this condition
		fmt.Println("Found a claw machine with aligned vectors :(")
		return
	} else {
		if comparison < 0 {
			// Const for line parallel to A that crosses through prize
			c := sub(Fraction{cm.prize[1], *big.NewInt(1)}, mul(Fraction{cm.prize[0], *big.NewInt(1)}, mA))

			// Coordinates of "turn"
			x := div(c, sub(mB, mA))
			y := mul(mB, x)

			if xInt, ok := tryReduce(x); ok {
				if _, ok := tryReduce(y); ok {
					bRem := new(big.Int)
					bRem.Rem(&xInt, &cm.b[0])

					aRem := new(big.Int)
					aRem.Sub(&cm.prize[0], &xInt)
					aRem.Rem(aRem, &cm.a[0])

					if bRem.Sign() == 0 && aRem.Sign() == 0 {
						bCount := new(big.Int)
						bCount.Div(&xInt, &cm.b[0])

						aCount := new(big.Int)
						aCount.Sub(&cm.prize[0], &xInt)
						aCount.Div(aCount, &cm.a[0])
						aCount.Mul(aCount, big.NewInt(3))

						result.Add(aCount, bCount)
					}
				}
			}
		} else {
			// Const for line parallel to B that crosses through prize
			c := sub(Fraction{cm.prize[1], *big.NewInt(1)}, mul(Fraction{cm.prize[0], *big.NewInt(1)}, mB))

			// Coordinates of "turn"
			x := div(c, sub(mA, mB))
			y := mul(mA, x)

			if xInt, ok := tryReduce(x); ok {
				if _, ok := tryReduce(y); ok {
					aRem := new(big.Int)
					aRem.Rem(&xInt, &cm.a[0])

					bRem := new(big.Int)
					bRem.Sub(&cm.prize[0], &xInt)
					bRem.Rem(bRem, &cm.b[0])

					if bRem.Sign() == 0 && aRem.Sign() == 0 {
						aCount := new(big.Int)
						aCount.Div(&xInt, &cm.a[0])
						aCount.Mul(aCount, big.NewInt(3))

						bCount := new(big.Int)
						bCount.Sub(&cm.prize[0], &xInt)
						bCount.Div(bCount, &cm.b[0])

						result.Add(aCount, bCount)
					}
				}
			}
		}
	}

	return
}

func parseInt(s string) big.Int {
	result, e := strconv.ParseInt(s, 0, 64)
	check(e)
	return *big.NewInt(result)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func add(x Fraction, y Fraction) Fraction {
	result := Fraction{}

	left := new(big.Int)
	left.Mul(&x.num, &y.den)
	right := new(big.Int)
	right.Mul(&y.num, &x.den)

	result.num.Add(left, right)
	result.den.Mul(&x.den, &y.den)

	return result
}

func sub(x Fraction, y Fraction) Fraction {
	result := Fraction{}

	left := new(big.Int)
	left.Mul(&x.num, &y.den)
	right := new(big.Int)
	right.Mul(&y.num, &x.den)

	result.num.Sub(left, right)
	result.den.Mul(&x.den, &y.den)

	return result
}

func mul(x Fraction, y Fraction) Fraction {
	result := Fraction{}

	result.num.Mul(&x.num, &y.num)
	result.den.Mul(&x.den, &y.den)

	return result
}

func div(x Fraction, y Fraction) Fraction {
	result := Fraction{}

	result.num.Mul(&x.num, &y.den)
	result.den.Mul(&x.den, &y.num)

	return result
}

func tryReduce(x Fraction) (result big.Int, ok bool) {
	temp := new(big.Int)
	temp.Rem(&x.num, &x.den)
	if temp.Sign() == 0 {
		ok = true
		result.Div(&x.num, &x.den)
	}

	return
}

func compare(x Fraction, y Fraction) int {
	left := new(big.Int)
	left.Mul(&x.num, &y.den)
	right := new(big.Int)
	right.Mul(&y.num, &x.den)
	return left.Cmp(right)
}
