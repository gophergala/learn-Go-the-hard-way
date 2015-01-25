package main

import (
	"testing"
)

func TestCheetSheet(t *testing.T) {
	out := CheatSheet("animation")
	if out != `//glasses
(•_•)
( •_•)>⌐■-■
(⌐■_■)

` {
		t.Fail()
	}
}
