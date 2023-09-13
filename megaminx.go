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
type State [12][10]int

type Megaminx struct {
	faces [12][5]int
	state State
}

// NewState generates the start state
func NewState() State {
	var s State
	for i := 0; i < 12; i++ {
		for j := 0; j < 10; j++ {
			s[i][j] = i
		}
	}
	return s
}

func NewMegaminx() Megaminx {

	s := NewState()
	m := Megaminx{state: s}

	m.faces[0] = [5]int{1, 2, 3, 4, 5}   // white
	m.faces[1] = [5]int{0, 5, 10, 9, 2}  // blue
	m.faces[2] = [5]int{0, 1, 9, 8, 3}   // yellow
	m.faces[3] = [5]int{0, 2, 8, 7, 4}   // purple
	m.faces[4] = [5]int{0, 3, 7, 11, 5}  // green
	m.faces[5] = [5]int{0, 4, 11, 10, 1} // red
	m.faces[6] = [5]int{7, 8, 9, 10, 11} // gray
	m.faces[7] = [5]int{6, 11, 4, 3, 8}  // cyan
	m.faces[8] = [5]int{6, 7, 3, 2, 9}   // orange
	m.faces[9] = [5]int{6, 8, 2, 1, 10}  // lime green
	m.faces[10] = [5]int{6, 9, 1, 5, 11} // pink
	m.faces[11] = [5]int{6, 10, 5, 4, 7} // vanilla

	return m
}

// for a better commented and similar explanation, look at CW
func (m *Megaminx) CCW(face int) {
	// shift tiles on face
	for i := 0; i < 2; i++ {
		m.state.shiftTilesRight(face)
	}

	adjFace := m.faces[face][4] // get last adjacent face
	adjRow := m.outerIndex(face, adjFace)

	endTileColors := make([]int, 3)
	for i := 0; i < 3; i++ {
		endTileColors[i] = m.state[adjFace][((adjRow*2)+i)%10]
	}
	var endColors [3]int
	copy(endColors[:], endTileColors)

	for i := 4; i > 0; i-- {
		prevFace := m.faces[face][i-1]
		prevAdjRow := m.outerIndex(face, prevFace)
		prevTileColors := make([]int, 3)
		for j := 0; j < 3; j++ {
			prevTileColors[j] = m.state[prevFace][((prevAdjRow*2)+j)%10]
		}

		curFace := m.faces[face][i]
		curAdjRow := m.outerIndex(face, curFace)
		for j := 0; j < 3; j++ {
			m.state[curFace][((curAdjRow*2)+j)%10] = prevTileColors[j]
		}
	}
	startFace := m.faces[face][0]
	startAdjRow := m.outerIndex(face, startFace)
	for i := 0; i < 3; i++ {
		m.state[startFace][((startAdjRow*2)+i)%10] = endColors[i]
	}
}

func (m *Megaminx) CW(face int) {
	// shift tiles on face
	for i := 0; i < 2; i++ {
		m.state.shiftTilesLeft(face)
	}

	// get first adjacent face color
	adjFace := m.faces[face][0]

	// find where the current face is in the first adjacent face's adj array
	adjRow := m.outerIndex(face, adjFace)

	// create a copy of them
	endTileColors := make([]int, 3)
	for i := 0; i < 3; i++ {
		endTileColors[i] = m.state[adjFace][((adjRow*2)+i)%10]
	}
	var endColors [3]int
	copy(endColors[:], endTileColors)

	for i := 0; i < 4; i++ {
		// get the next face color in the adjacency array of current face
		prevFace := m.faces[face][i+1]
		// find which edge the current face is adjacent to
		prevAdjRow := m.outerIndex(face, prevFace)
		prevTileColors := make([]int, 3)
		for j := 0; j < 3; j++ {
			prevTileColors[j] = m.state[prevFace][((prevAdjRow*2)+j)%10]
		}

		curFace := m.faces[face][i]
		curAdjRow := m.outerIndex(face, curFace)
		for j := 0; j < 3; j++ {
			m.state[curFace][((curAdjRow*2)+j)%10] = prevTileColors[j]
		}
	}
	endFace := m.faces[face][4]
	endAdjRow := m.outerIndex(face, endFace)
	for i := 0; i < 3; i++ {
		m.state[endFace][((endAdjRow*2)+i)%10] = endColors[i]
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
func (m *Megaminx) outerIndex(f1, f2 int) int {
	for i, e := range m.faces[f2] {
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

	tiles := megaminx.state[face]
	// paint the first five vertices in vs with tiles[0], next five with tiles[1], ...
	for j := 0; j < len(tiles); j++ {
		for k := 0; k < 5; k++ {
			idx := (i + (5 * j) + k) % 56
			vs[idx].ColorR = float32(numToColor[tiles[j]].R) / 255
			vs[idx].ColorG = float32(numToColor[tiles[j]].G) / 255
			vs[idx].ColorB = float32(numToColor[tiles[j]].B) / 255
		}
	}
}

// randomize makes 20 random clockwise turns on the faces of the megaminx
func randomize(moves int) {
	for i := 0; i < moves; i++ {
		face := rand.Intn(12)
		megaminx.CW(face)
	}
}
