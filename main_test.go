package main

import (
	"runtime"
	"testing"
)

func TestParallelMST(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	prob = [n]float64{5, 10, 2, 3, 4}
	for d := 0; d < vp; d++ { //sub-diagonal of j=i+d
		for i := 0; i+d < vp; i++ {
			//runs chunk (i,i+d)
			go chunk(i, i+d)
		}
	}
	<-finish
	if P(0, 5, 0, 5) != `0 5 20 24 32 44 
0 0 10 14 22 34 
0 0 0 2 7 15 
0 0 0 0 3 10 
0 0 0 0 0 4 
0 0 0 0 0 0 
` {
		t.Fail()
	}
}
