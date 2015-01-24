package main

import (
	"fmt"
)

const (
	n  = 5 //scale of the tree.
	vp = 2 //number of chunks on one demension,totally vp(vp+1)/2
)

var (
	cost   [n + 1][n + 1]float64 //cost for tree(i-1,j)
	root   [n + 1][n + 1]int     //root for tree(i-1,j)
	prob   [n]float64            //probability for n nodes.
	h, v   [vp][vp]chan struct{} //horizontal channel and vertical channel.
	finish chan struct{}         //if reach the chunk(0,vp-1),send finish signal.
)

//prints from (il,jl) to (ih,jh)
func P(il, ih, jl, jh int) string {
	var out string
	for ii := il; ii <= ih; ii++ {
		for jj := jl; jj <= jh; jj++ {
			out += fmt.Sprint(cost[ii][jj], " ")
		}
		out += fmt.Sprintln()
	}
	return out
}

//get sequential mean search time.
func MstSeq(i, j int) {
	var bestCost float64 = 1e9 + 0.0
	var bestRoot int = -1
	switch {
	case i >= j: //empty tree
		cost[i][j] = 0.0
		root[i][j] = -1
	case i+1 == j: //single tree
		cost[i][j] = prob[i]
		root[i][j] = i + 1
	case i+1 < j:
		psum := 0.0
		for k := i; k <= j-1; k++ {
			psum += prob[k]
		}
		for r := i; r <= j-1; r++ {
			rcost := cost[i][r] + cost[r+1][j]
			if rcost < bestCost {
				bestCost = rcost
				bestRoot = r + 1
			}
			cost[i][j] = bestCost + psum
			root[i][j] = bestRoot
		}
	}
}

//compute chunk(i,j)
var chunk func(i, j int)

func init() {
	for i := 0; i < vp; i++ {
		for j := i; j < vp; j++ {
			if j < vp-1 {
				h[i][j] = make(chan struct{})
			}
			if i > 0 {
				v[i][j] = make(chan struct{})
			}
		}
	}
	finish = make(chan struct{})
}

func main() {
	prob = [n]float64{5, 10, 2, 3, 4}
	for i := n; i >= 0; i-- {
		for j := i; j <= n; j++ {
			MstSeq(i, j)
		}
	}
	if P(0, n, 0, n) == `0 5 20 24 32 44 
0 0 10 14 22 34 
0 0 0 2 7 15 
0 0 0 0 3 10 
0 0 0 0 0 4 
0 0 0 0 0 0 
` {
		println(`Again,Go is designed for concurrency,it makes it easy to parrallel.This exercise is about parrallel dynamic programming.
We use the optimal binary search tree problem to show the multi-core parallel code in Go.
The mean search time for tree(nodes form i+1 to j) is stored in cost[i][j],similarily root of the search tree is stored in root[i][j]
And a tree mean search time is the total probability of all nodes plus the minium mean time of 2 possible sub trees.
The sequential MST method has been given in the main.go
In fact the every cost[i][j] depencies on their left and bottom neighbor.
So we can divide the computation into chunks like that (assume n is 5):
  
  * * *    | * * *
    * * -> | * * *
      *    | * * *
           + - - -
               ^
               |
             * * *
               * *
                 *

denotes left-top chunk(0,0),right-top chunk(0,1),right-bottom chunk(1,1).(0,1)dependencies (0,1) and (1,1),(0,1) and (1,1) can compute imeadiately.You should use channels to syncronize computation,please completes the func 'chunk' to pass the test"
`)
	}
}
