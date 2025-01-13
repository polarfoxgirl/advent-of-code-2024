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

	prices := make([][2001]int, len(secrets))
	diffs := make([][2000]int, len(secrets))
	for b, secret := range secrets {
		prices[b][0] = secret % 10
		random := secret
		for n := 1; n <= 2000; n++ {
			random = next(random)
			prices[b][n] = random % 10
			diffs[b][n-1] = prices[b][n] - prices[b][n-1]
		}
	}

	// Get all relevant sequences
	seqs := make(map[[4]int]int)
	for b := 0; b < len(secrets); b++ {
		seenSeqForBuyer := make(map[[4]int]struct{})
		for i := 0; i < 1997; i++ {
			seq := [4]int{diffs[b][i], diffs[b][i+1], diffs[b][i+2], diffs[b][i+3]}

			// If first encounter for this buyer, record price
			if _, seen := seenSeqForBuyer[seq]; !seen {
				seenSeqForBuyer[seq] = struct{}{}
				// Prices has an extra value in front
				seqs[seq] += prices[b][i+4]
			}
		}
	}
	fmt.Printf("Considering %d sequences\n", len(seqs))

	result := 0
	finalSeq := [4]int{}
	for seq, total := range seqs {
		if total > result {
			result = total
			finalSeq = seq
		}
	}

	fmt.Printf("Result: %d for sequence %v\n", result, finalSeq)
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
