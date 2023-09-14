package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math/rand"
)

var numToColor = map[int]color.RGBA{
	// top half
	0: {0xff, 0xff, 0xff, 0xff}, // white
	1: {0, 0, 0xff, 0xff},       // blue
	2: {0xff, 0xff, 0, 0xff},    // yellow
	3: {80, 00, 80, 0xff},       // purple
	4: {0, 0x64, 0, 0xff},       // green
	5: {0xff, 0, 0, 0xff},       // red

	// bottom half
	6:  {80, 80, 80, 0xff},       // gray
	7:  {00, 0xff, 0xff, 0xff},   // cyan
	8:  {0xff, 0xa5, 0, 0xff},    // orange
	9:  {0x65, 0xfe, 0x08, 0xff}, // lime green
	10: {0xfc, 0x6c, 0x85, 0xff}, // pink
	11: {0xf3, 0xe5, 0xab, 0xff}, // vanilla
}

// State contains all state needed to specify a Megaminx; the tile colors of each of the 12 faces
type State [12][10]byte

// NewState generates the start state
func NewState() State {
	var s State
	for i := 0; i < 12; i++ {
		for j := 0; j < 10; j++ {
			s[i][j] = byte(i)
		}
	}
	return s
}

func CopyState(s State) State {
	var n State
	for i := 0; i < 12; i++ {
		for j := 0; j < 10; j++ {
			n[i][j] = s[i][j]
		}
	}
	return n
}

// 2D adjacency array for megaminx
var m = [12][5]int{
	{1, 2, 3, 4, 5},   // white
	{0, 5, 10, 9, 2},  // blue
	{0, 1, 9, 8, 3},   // yellow
	{0, 2, 8, 7, 4},   // purple
	{0, 3, 7, 11, 5},  // green
	{0, 4, 11, 10, 1}, // red
	{7, 8, 9, 10, 11}, // gray
	{6, 11, 4, 3, 8},  // cyan
	{6, 7, 3, 2, 9},   // orange
	{6, 8, 2, 1, 10},  // lime green
	{6, 9, 1, 5, 11},  // pink
	{6, 10, 5, 4, 7},  // vanilla
}

// for a better commented and similar explanation, look at CW
func (s *State) CCW(face int) {
	// shift tiles on face
	for i := 0; i < 2; i++ {
		s.shiftTilesRight(face)
	}

	adjFace := m[face][4] // get last adjacent face
	adjRow := outerIndex(face, adjFace)

	endTileColors := make([]byte, 3)
	for i := 0; i < 3; i++ {
		endTileColors[i] = s[adjFace][((adjRow*2)+i)%10]
	}
	var endColors [3]byte
	copy(endColors[:], endTileColors)

	for i := 4; i > 0; i-- {
		prevFace := m[face][i-1]
		prevAdjRow := outerIndex(face, prevFace)
		prevTileColors := make([]byte, 3)
		for j := 0; j < 3; j++ {
			prevTileColors[j] = s[prevFace][((prevAdjRow*2)+j)%10]
		}

		curFace := m[face][i]
		curAdjRow := outerIndex(face, curFace)
		for j := 0; j < 3; j++ {
			s[curFace][((curAdjRow*2)+j)%10] = prevTileColors[j]
		}
	}
	startFace := m[face][0]
	startAdjRow := outerIndex(face, startFace)
	for i := 0; i < 3; i++ {
		s[startFace][((startAdjRow*2)+i)%10] = endColors[i]
	}
}
func (s *State) CW(face int) {
	// shift tiles on face
	for i := 0; i < 2; i++ {
		s.shiftTilesLeft(face)
	}

	// get first adjacent face color
	adjFace := m[face][0]

	// find where the current face is in the first adjacent face's adj array
	adjRow := outerIndex(face, adjFace)

	// create a copy of them
	endTileColors := make([]byte, 3)
	for i := 0; i < 3; i++ {
		endTileColors[i] = s[adjFace][((adjRow*2)+i)%10]
	}
	var endColors [3]byte
	copy(endColors[:], endTileColors)

	for i := 0; i < 4; i++ {
		// get the next face color in the adjacency array of current face
		prevFace := m[face][i+1]
		// find which edge the current face is adjacent to
		prevAdjRow := outerIndex(face, prevFace)
		prevTileColors := make([]byte, 3)
		for j := 0; j < 3; j++ {
			prevTileColors[j] = s[prevFace][((prevAdjRow*2)+j)%10]
		}

		curFace := m[face][i]
		curAdjRow := outerIndex(face, curFace)
		for j := 0; j < 3; j++ {
			s[curFace][((curAdjRow*2)+j)%10] = prevTileColors[j]
		}
	}
	endFace := m[face][4]
	endAdjRow := outerIndex(face, endFace)
	for i := 0; i < 3; i++ {
		s[endFace][((endAdjRow*2)+i)%10] = endColors[i]
	}
}

func (s *State) shiftTilesRight(face int) {
	end := s[face][9]
	for i := 9; i > 0; i-- {
		s[face][i] = s[face][i-1]
	}
	s[face][0] = end
}

func (s *State) shiftTilesLeft(face int) {
	start := s[face][0]
	for i := 0; i < 9; i++ {
		s[face][i] = s[face][i+1]
	}
	s[face][9] = start
}

// outerIndex returns the index of face f1 in the adjacency array of f2, an operation used in rotation
func outerIndex(f1, f2 int) int {
	for i, e := range m[f2] {
		if e == f1 {
			return i
		}
	}
	panic("f1 not in f2's adjacency array")
}

func PaintFace(vs []ebiten.Vertex, face int) {
	// first six vertices of vs are painted the color of the face
	for i := 0; i < 6; i++ {
		vs[i].ColorR = float32(numToColor[face].R) / 255
		vs[i].ColorG = float32(numToColor[face].G) / 255
		vs[i].ColorB = float32(numToColor[face].B) / 255
	}
	i := 6 // index into vs

	tiles := state[face]
	// paint the first five vertices in vs with tiles[0], next five with tiles[1], ...
	for j := 0; j < len(tiles); j++ {
		for k := 0; k < 5; k++ {
			idx := (i + (5 * j) + k) % 56
			vs[idx].ColorR = float32(numToColor[int(tiles[j])].R) / 255
			vs[idx].ColorG = float32(numToColor[int(tiles[j])].G) / 255
			vs[idx].ColorB = float32(numToColor[int(tiles[j])].B) / 255
		}
	}
}

// randomize makes 20 random clockwise turns on the faces of the megaminx
func (s *State) randomize(moves int) {
	for i := 0; i < moves; i++ {
		face := rand.Intn(12)
		s.CW(face)
	}
}
