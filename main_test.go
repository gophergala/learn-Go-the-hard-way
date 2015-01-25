package main

import (
	"testing"
)

func TestLex(t *testing.T) {
	var tokens []token
	input := "100.11 - 10.e2"
	l := newLexer(input)
	for t := l.token(); t.typ != tEOF; t = l.token() {
		tokens = append(tokens, t)
	}
	if len(tokens) != 3 {
		t.Fail()
	}
	if tokens[0].typ != tNUM && tokens[0].lit != "100.11" {
		t.Fail()
	}
	if tokens[1].typ != tMUNIS && tokens[1].lit != "-" {
		t.Fail()
	}
	if tokens[2].typ != tNUM && tokens[2].lit != "10.e2" {
		t.Fail()
	}
	input = ".11+22.1e2"
	l = newLexer(input)
	tokens = tokens[:0]
	for t := l.token(); t.typ != tEOF; t = l.token() {
		tokens = append(tokens, t)
	}
	if tokens[0].typ != tNUM && tokens[0].lit != ".11" {
		t.Fail()
	}
	if tokens[1].typ != tMUNIS && tokens[1].lit != "+" {
		t.Fail()
	}
	if tokens[2].typ != tNUM && tokens[2].lit != "22.1e2" {
		t.Fail()
	}
}
