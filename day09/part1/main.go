package main

import (
	"fmt"
	"os"
	// "reflect"
	// "math"
	// "slices"
	"strconv"
	// "iter"
	// "strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	files, free := parseInput(data)
	fmt.Printf("Input: %d files, %d total free space\n", len(files), len(free))

	compact(files, free)

	sum := 0
	for id, file := range files {
		for _, p := range file {
			sum += id * p
		}
	}

	fmt.Printf("Result: %d", sum)
}

func parseInput(data []uint8) (files map[int][]int, free []int) {
	// Can use a slice or array instead of a map
	files = make(map[int][]int)
	free = make([]int, 0)
	pointer := 0
	for i, rune := range string(data) {
		size64, e := strconv.ParseInt(string(rune), 0, 8)
		check(e)
		size := int(size64)

		if i%2 == 0 {
			id := i / 2
			files[id] = make([]int, size)
			for chunk := 0; chunk < size; chunk++ {
				files[id][chunk] = pointer
				pointer++
			}
		} else {
			for chunk := 0; chunk < size; chunk++ {
				free = append(free, pointer)
				pointer++
			}
		}
	}
	return
}

func compact(files map[int][]int, free []int) {
	iFree := 0
	for id := len(files) - 1; id >= 0; id-- {
		file := files[id]
		for chunk := len(files[id]) - 1; chunk >= 0; chunk-- {
			if iFree == len(free) || free[iFree] >= file[chunk] {
				return
			}

			file[chunk] = free[iFree]
			iFree++
		}
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
