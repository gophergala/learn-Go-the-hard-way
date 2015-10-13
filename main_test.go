package main

import (
	"testing"
)

type myType int

func TestReverse(t *testing.T) {
	var slice = []myType{0, 1, 2}
	Reverse(&slice)
	if slice[0] != myType(2) || slice[1] != myType(1) || slice[2] != myType(0) || len(slice) != 3 {
		t.Fail()
	}
}
