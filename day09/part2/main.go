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

type chunk struct {
	pos  int
	size int
}

func main() {
	data, err := os.ReadFile("input.txt")
	check(err)

	files, free := parseInput(data)
	fmt.Printf("Input: %d files, %d total free space\n", len(files), len(free))

	// print(files, free)

	compact(files, free)

	// print(files, free)

	sum := 0
	for id, file := range files {
		for p := file.pos; p < file.pos+file.size; p++ {
			sum += id * p
		}
	}

	fmt.Printf("Result: %d", sum)
}

func parseInput(data []uint8) (files []chunk, free []chunk) {
	text := string(data)
	files = make([]chunk, 0, len(text)/2+1)
	free = make([]chunk, 0, len(text)/2+1)
	pointer := 0
	for i, rune := range text {
		size64, e := strconv.ParseInt(string(rune), 0, 8)
		check(e)
		size := int(size64)
		next := chunk{pointer, size}

		if i%2 == 0 {
			files = append(files, next)
		} else {
			free = append(free, next)
		}
		pointer += size
	}
	return
}

func compact(files []chunk, free []chunk) {
	for id := len(files) - 1; id >= 0; id-- {
		file := &files[id]
		for freeId := 0; freeId < len(free); freeId++ {
			freeChunk := &free[freeId]
			if freeChunk.pos >= file.pos {
				break
			}
			if freeChunk.size >= file.size {
				file.pos = freeChunk.pos
				freeChunk.pos += file.size
				freeChunk.size -= file.size
				break
			}
		}
	}
}

func print(files []chunk, free []chunk) {
	p := 0
	filesPrinted := 0
	for filesPrinted < len(files) {
		found := false
		for id, ch := range files {
			if ch.pos == p {
				for c := 0; c < ch.size; c++ {
					fmt.Printf("%d", id)
				}
				p += ch.size
				filesPrinted++
				found = true
				break
			}
		}
		if found {
			continue
		}

		for _, ch := range free {
			if ch.pos == p {
				for c := 0; c < ch.size; c++ {
					fmt.Print(".")
				}
				p += ch.size
				if ch.size > 0 {
					found = true
				}
				break
			}
		}

		if !found {
			fmt.Print(".")
			p++
		}

		if p > 100 {
			break
		}
	}
	fmt.Println()
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
