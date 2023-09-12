package main

import "math"

type Node struct {
	s State
	g int
}

// H returns the heuristic value for a node

func G(n Node) int {
	return n.g
}

func H(n Node) int {
	wrong := 0 // count of stickers on the wrong face
	for i := 0; i < 12; i++ {
		for j := 0; j < 10; j++ {
			if n.s[i][j] != i {
				wrong++
			}
		}
	}
	return int(math.Ceil(float64(wrong) / 15.0)) // ceil(wrong / 15)
}

func Expand(n Node)

// G returns the cost so far for a node
func (n Node) G() int {
	return n.g
}
