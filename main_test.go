package main

import (
	"testing"
)

func TestMap(t *testing.T) {
	var f func(func(e int) int, []int) []int
	MakeMap(&f)

	var a = []int{1, 2, 3}

	b := f(func(e int) int {
		return e + 1
	}, a)

	if b[0] != 2 || b[1] != 3 || b[2] != 4 || len(b) != 3 {
		t.Fail()
	}

	var f2 func(func(e int) int, map[int]int) map[int]int

	MakeMap(&f2)

	var c = map[int]int{0: 0, 1: 1, 2: 3}

	d := f2(func(e int) int {
		return e + 1
	}, c)

	if d[0] != 1 || d[1] != 2 || d[2] != 4 || len(d) != 3 {
		t.Fail()
	}
}
