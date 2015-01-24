package main

import (
	"testing"
)

func TestGame(t *testing.T) {
	win := Game()
	for _, w := range win {
		if !w {
			t.Fail()
		}
	}
}
