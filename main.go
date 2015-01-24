package main

import (
	"sync"
)

//Parallelsum does parallel vector sum,
//in each loop,buffered input capacity will be cut half
//for the next loop to Sum goroutine to consume.
//it terminates untill the input remains one.
func ParallelSum(slcs ...[]int) []int {
	input := make(chan []int, len(slcs))
	output := make(chan []int)
	var result []int
	go func(input chan []int) {
		for _, slc := range slcs {
			input <- slc
		}
		close(input)
	}(input)

	for {
		var wg sync.WaitGroup
		wg.Add(cap(input) / 2)
		for i := 0; i < cap(input)/2; i++ {
			out := Sum(input)
			go func() {
				defer wg.Done()
				for o := range out {
					output <- o
				}
			}()
		}
		go func(output chan []int) {
			wg.Wait()
			close(output)
		}(output)

		input = make(chan []int, cap(input)/2)
		if cap(input) < 2 {
			result = <-output
			break
		}
		for o := range output {
			input <- o
		}
		output = make(chan []int)
		close(input)
	}
	return result
}

//TODO:complete the Sum for the parallel sum function.
var Sum func(sum chan []int) (output chan []int)

func main() {
	println("Please edit main.go,and complete the 'Sum' function for the parallel sum to pass the test.\nConcurrency is the most important feature of Go,and the principle is\n'Do not communicate by sharing memory; instead, share memory by communicating.'\nIn this exercise you need to catch many features of channels.This is a tour for you to figure out!")
}
