package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	secrets := parseInput(data)
	fmt.Printf("Input: %d secrets\n", len(secrets))

	result := 0
	for _, secret := range secrets {
		result += repeatNext(secret, 2000)
	}

	fmt.Printf("Result: %d\n", result)
}

func repeatNext(secret int, times int) (result int) {
	result = secret
	for i := 0; i < times; i++ {
		result = next(result)
	}
	return
}

func next(secret int) (result int) {
	result = prune(mix(secret*64, secret))
	result = prune(mix(result/32, result))
	result = prune(mix(result*2048, result))
	return
}

func mix(n int, secret int) int {
	return n ^ secret
}

func prune(n int) int {
	return n % 16777216
}

func parseInput(data []uint8) (secrets []int) {
	lines := strings.Split(string(data), "\n")
	secrets = make([]int, len(lines))

	for i, line := range lines {
		secrets[i] = parseInt(line)
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
