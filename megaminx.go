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

type Megaminx struct {
	faces       [12]Face
	top, bottom int
}

type Face struct {
	adj   [5]Edge // array of edges, where the first element is the "bottom" edge
	tiles [10]int // hold the actual color of the tiles
}

type Edge struct {
	color  int    // adjacent face's color
	colors [3]int // array of indices into tiles in Face
}

func NewMegaminx() Megaminx {
	m := Megaminx{}

	m.top = 0
	m.bottom = 6

	// white face
	m.faces[0] = Face{
		adj: [5]Edge{
			{1, [3]int{0, 1, 2}},
			{2, [3]int{2, 3, 4}},
			{3, [3]int{4, 5, 6}},
			{4, [3]int{6, 7, 8}},
			{5, [3]int{8, 9, 0}},
		},
		tiles: [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	// blue face
	m.faces[1] = Face{
		adj: [5]Edge{
			{0, [3]int{0, 1, 2}},
			{5, [3]int{2, 3, 4}},
			{10, [3]int{4, 5, 6}},
			{9, [3]int{6, 7, 8}},
			{2, [3]int{8, 9, 0}},
		},
		tiles: [10]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	// yellow face
	m.faces[2] = Face{
		adj: [5]Edge{
			{0, [3]int{0, 1, 2}},
			{1, [3]int{2, 3, 4}},
			{9, [3]int{4, 5, 6}},
			{8, [3]int{6, 7, 8}},
			{3, [3]int{8, 9, 0}},
		},
		tiles: [10]int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
	}

	// purple face
	m.faces[3] = Face{
		adj: [5]Edge{
			{0, [3]int{0, 1, 2}},
			{2, [3]int{2, 3, 4}},
			{8, [3]int{4, 5, 6}},
			{7, [3]int{6, 7, 8}},
			{4, [3]int{8, 9, 0}},
		},
		tiles: [10]int{3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
	}

	// green face
	m.faces[4] = Face{
		adj: [5]Edge{
			{0, [3]int{0, 1, 2}},
			{3, [3]int{2, 3, 4}},
			{7, [3]int{4, 5, 6}},
			{11, [3]int{6, 7, 8}},
			{5, [3]int{8, 9, 0}},
		},
		tiles: [10]int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
	}

	// red face
	m.faces[5] = Face{
		adj: [5]Edge{
			{0, [3]int{0, 1, 2}},
			{4, [3]int{2, 3, 4}},
			{11, [3]int{4, 5, 6}},
			{10, [3]int{6, 7, 8}},
			{1, [3]int{8, 9, 0}},
		},
		tiles: [10]int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5},
	}
	// gray face
	m.faces[6] = Face{
		adj: [5]Edge{
			{7, [3]int{0, 1, 2}},
			{8, [3]int{2, 3, 4}},
			{9, [3]int{4, 5, 6}},
			{10, [3]int{6, 7, 8}},
			{11, [3]int{8, 9, 0}},
		},
		tiles: [10]int{6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
	}

	// cyan face
	m.faces[7] = Face{
		adj: [5]Edge{
			{6, [3]int{0, 1, 2}},
			{11, [3]int{2, 3, 4}},
			{4, [3]int{4, 5, 6}},
			{3, [3]int{6, 7, 8}},
			{8, [3]int{8, 9, 0}},
		},
		tiles: [10]int{7, 7, 7, 7, 7, 7, 7, 7, 7, 7},
	}

	// orange face
	m.faces[8] = Face{
		adj: [5]Edge{
			{6, [3]int{0, 1, 2}},
			{7, [3]int{2, 3, 4}},
			{3, [3]int{4, 5, 6}},
			{2, [3]int{6, 7, 8}},
			{9, [3]int{8, 9, 0}},
		},
		tiles: [10]int{8, 8, 8, 8, 8, 8, 8, 8, 8, 8},
	}

	// lime green face
	m.faces[9] = Face{
		adj: [5]Edge{
			{6, [3]int{0, 1, 2}},
			{8, [3]int{2, 3, 4}},
			{2, [3]int{4, 5, 6}},
			{1, [3]int{6, 7, 8}},
			{10, [3]int{8, 9, 0}},
		},
		tiles: [10]int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
	}

	// pink face
	m.faces[10] = Face{
		adj: [5]Edge{
			{6, [3]int{0, 1, 2}},
			{9, [3]int{2, 3, 4}},
			{1, [3]int{4, 5, 6}},
			{5, [3]int{6, 7, 8}},
			{11, [3]int{8, 9, 0}},
		},
		tiles: [10]int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
	}

	// vanilla face
	m.faces[11] = Face{
		adj: [5]Edge{
			{6, [3]int{0, 1, 2}},
			{10, [3]int{2, 3, 4}},
			{5, [3]int{4, 5, 6}},
			{4, [3]int{6, 7, 8}},
			{7, [3]int{8, 9, 0}},
		},
		tiles: [10]int{11, 11, 11, 11, 11, 11, 11, 11, 11, 11},
	}
	return m
}

// for a better commented and similar explanation, look at CW
func (m *Megaminx) CCW(face int) {
	// shift tiles on face
	for i := 0; i < 2; i++ {
		m.faces[face].shiftTilesRight()
	}

	adjFace := m.faces[face].adj[4].color // get last adjacent face
	adjRow := m.outerIndex(face, adjFace)
	adjTiles := m.faces[adjFace].adj[adjRow].colors
	endTileColors := make([]int, 3)
	for i := range adjTiles {
		endTileColors[i] = m.faces[adjFace].tiles[adjTiles[i]]
	}
	var endColors [3]int
	copy(endColors[:], endTileColors)

	for i := 4; i > 0; i-- {
		prevFace := m.faces[face].adj[i-1].color
		prevAdjRow := m.outerIndex(face, prevFace)
		prevAdjTiles := m.faces[prevFace].adj[prevAdjRow].colors
		prevTileColors := make([]int, 3)
		for j := range prevAdjTiles {
			prevTileColors[j] = m.faces[prevFace].tiles[prevAdjTiles[j]]
		}

		curFace := m.faces[face].adj[i].color
		curAdjRow := m.outerIndex(face, curFace)
		curAdjTiles := m.faces[curFace].adj[curAdjRow].colors

		for k, j := range curAdjTiles {
			m.faces[curFace].tiles[j] = prevTileColors[k]
		}
	}
	startFace := m.faces[face].adj[0].color
	startAdjRow := m.outerIndex(face, startFace)
	startAdjTiles := m.faces[startFace].adj[startAdjRow].colors
	for k, j := range startAdjTiles {
		m.faces[startFace].tiles[j] = endColors[k]
	}

}

func (m *Megaminx) CW(face int) {
	// shift tiles on face
	for i := 0; i < 2; i++ {
		m.faces[face].shiftTilesLeft()
	}

	// get first adjacent face color
	adjFace := m.faces[face].adj[0].color

	// find where the current face is in the first adjacent face's adj array
	adjRow := m.outerIndex(face, adjFace)

	// get the adjacent tile colors
	adjTiles := m.faces[adjFace].adj[adjRow].colors

	// create a copy of them
	endTileColors := make([]int, 3)
	for i := range adjTiles {
		endTileColors[i] = m.faces[adjFace].tiles[adjTiles[i]]
	}
	var endColors [3]int
	copy(endColors[:], endTileColors)

	for i := 0; i < 4; i++ {
		// get the next face color in the adjacency array of current face
		prevFace := m.faces[face].adj[i+1].color
		// find which edge the current face is adjacent to
		prevAdjRow := m.outerIndex(face, prevFace)
		prevAdjTiles := m.faces[prevFace].adj[prevAdjRow].colors

		prevTileColors := make([]int, 3)
		for j := range prevAdjTiles {
			prevTileColors[j] = m.faces[prevFace].tiles[prevAdjTiles[j]]
		}

		curFace := m.faces[face].adj[i].color
		curAdjRow := m.outerIndex(face, curFace)
		curAdjTiles := m.faces[curFace].adj[curAdjRow].colors

		for k, j := range curAdjTiles {
			m.faces[curFace].tiles[j] = prevTileColors[k]
		}
	}
	endFace := m.faces[face].adj[4].color
	endAdjRow := m.outerIndex(face, endFace)
	endAdjTiles := m.faces[endFace].adj[endAdjRow].colors
	for k, j := range endAdjTiles {
		m.faces[endFace].tiles[j] = endColors[k]
	}
}

func (f *Face) shiftTilesRight() {
	end := f.tiles[len(f.tiles)-1]
	for i := len(f.tiles) - 1; i > 0; i-- {
		f.tiles[i] = f.tiles[i-1]
	}
	f.tiles[0] = end
}

func (f *Face) shiftTilesLeft() {
	start := f.tiles[0]
	for i := 0; i < len(f.tiles)-1; i++ {
		f.tiles[i] = f.tiles[i+1]
	}
	f.tiles[len(f.tiles)-1] = start
}

// outerIndex returns the index of face f1 in the adjacency array of f2, an operation used in rotation
func (m *Megaminx) outerIndex(f1, f2 int) int {
	for i, e := range m.faces[f2].adj {
		if e.color == f1 {
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

	tiles := megaminx.faces[face].tiles
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
