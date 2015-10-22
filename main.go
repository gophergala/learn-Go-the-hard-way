package main

import (
	"math/rand"
	"time"
)

const (
	ROCK int = iota
	PAPER
	SCISSORS
)

type Choice struct {
	Who   int //0 you 1 your opponent
	Guess int
}

//Win returns true if you win.
func Win(you, he int) bool {
	if you == ROCK && he == SCISSORS {
		return true
	}
	if you == PAPER && he == ROCK {
		return true
	}
	if you == SCISSORS && he == PAPER {
		return true
	}
	return false
}

func Opponent(guess chan Choice, please chan struct{}) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 3; i++ {
		<-please
		choice := r.Intn(3)
		who := 1
		guess <- Choice{who, choice}
		please <- struct{}{}
	}
}

var Cheat func(guess chan Choice) chan Choice

func Me(guess chan Choice, please chan struct{}) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 3; i++ {
		<-please
		choice := r.Intn(3)
		who := 0
		guess <- Choice{who, choice}
		please <- struct{}{}
	}
}

func Game() []bool {
	guess := make(chan Choice)
	//please sync 2 goroutines.
	please := make(chan struct{})
	go func() { please <- struct{}{} }()
	go Opponent(guess, please)
	go Me(guess, please)
	guess = Cheat(guess)
	var wins []bool

	for i := 0; i < 3; i++ {
		g1 := <-guess
		g2 := <-guess
		win := false
		if g1.Who == 0 {
			win = Win(g1.Guess, g2.Guess)
		} else {
			win = Win(g2.Guess, g1.Guess)
		}
		wins = append(wins, win)
	}

	return wins
}

func main() {
	println("Now let's play a game 'rock-paper-scissors',there are 2 players-you and a goroutine!\nTo be bound to win,you should call a goroutine to help you to peer what your opponent choose.\nTwo out of three.\nPlease edit main.go to complete func 'Cheat' to win!")
}
