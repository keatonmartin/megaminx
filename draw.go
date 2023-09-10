package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math"
)

func degToR(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// rotateVertices rotates the vertices in vs with respect to the origin
// by r radians. Equivalent to multiplication with rotation matrix.
func rotateVertices(vs []ebiten.Vertex, r float64) {
	for i := range vs {
		x := vs[i].DstX
		y := vs[i].DstY
		vs[i].DstX = x*float32(math.Cos(r)) - y*float32(math.Sin(r))
		vs[i].DstY = x*float32(math.Sin(r)) + y*float32(math.Cos(r))
	}
}

// scaleVertices scales the vertices in vs by the specified scale factor
func scaleVertices(vs []ebiten.Vertex, scale float32) {
	for i := range vs {
		vs[i].DstX *= scale
		vs[i].DstY *= scale
	}
}

func translateVertices(vs []ebiten.Vertex, x, y float32) {
	for i := range vs {
		vs[i].DstX += x
		vs[i].DstY += y
	}
}

// reflectVertices reflects the vertics in vs over y = -x
func reflectVertices(vs []ebiten.Vertex) {
	for i := range vs {
		x := vs[i].DstX
		y := vs[i].DstY
		vs[i].DstX = -y
		vs[i].DstY = -x
	}
}

func reflectVerticesOverY(vs []ebiten.Vertex) {
	for i := range vs {
		vs[i].DstX = -vs[i].DstX
	}
}

func drawSelectors(screen *ebiten.Image) {
	var path vector.Path

	path.MoveTo(0, 0)
	path.LineTo(20, 0)
	path.LineTo(20, 20)
	path.LineTo(0, 20)
	path.Close()

	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	translateVertices(vs, 100, 10)

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true

	for i := 0; i < 12; i++ {
		newVs := make([]ebiten.Vertex, len(vs))
		copy(newVs, vs)
		translateVertices(newVs, float32(i*20), 0)
		for j := range newVs {
			newVs[j].ColorR = float32(numToColor[i].R) / 255
			newVs[j].ColorG = float32(numToColor[i].G) / 255
			newVs[j].ColorB = float32(numToColor[i].B) / 255
		}
		selectors[i] = newVs
		screen.DrawTriangles(newVs, is, whiteImage, op)
		translateVertices(vs, 5, 0)
	}
}

func drawMarker(screen *ebiten.Image, s int) {
	var path vector.Path

	path.MoveTo(selectors[s][0].DstX, selectors[s][2].DstY+5)
	path.LineTo(selectors[s][1].DstX, selectors[s][2].DstY+5)
	path.LineTo(selectors[s][1].DstX, selectors[s][2].DstY+10)
	path.LineTo(selectors[s][0].DstX, selectors[s][2].DstY+10)
	path.Close()

	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)

	for i := range vs {
		vs[i].ColorR = 1
		vs[i].ColorG = 1
		vs[i].ColorB = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true
	screen.DrawTriangles(vs, is, whiteImage, op)
}

func getFacePath() ([]ebiten.Vertex, []uint16) {
	var path vector.Path
	// https://mathworld.wolfram.com/RegularPentagon.html
	c1 := float32(math.Cos(2 * math.Pi / 5))
	c2 := float32(math.Cos(math.Pi / 5))
	s1 := float32(math.Sin(2 * math.Pi / 5))
	s2 := float32(math.Sin(4 * math.Pi / 5))

	// inner pentagon
	path.MoveTo(-s2, -c2) // bottom left
	path.LineTo(s2, -c2)  // bottom right
	path.LineTo(s1, c1)   // top right
	path.LineTo(0, 1)     // top middle
	path.LineTo(-s1, c1)  // top left
	path.Close()

	// corner points for the outer pentagon
	s1p := s1 * 2
	s2p := s2 * 2
	c1p := c1 * 2
	c2p := c2 * 2

	// bottom left corner
	path.MoveTo(-s2p, -c2p)
	path.LineTo(-s2p-(-s2p+s1p)*.3, -c2p+(c1p+c2p)*.3)
	path.LineTo(-s2*1.1, -c2*1.1)
	path.LineTo(-s2p+(s2p+s2p)*.3, -c2p)
	path.Close()

	// bottom edge
	path.MoveTo(-s2p+(s2p+s2p)*.34, -c2p)
	path.LineTo((-s2p+(s2p+s2p)*.34)*-1, -c2p)
	path.LineTo(s2-(s1-s2)*.07, -c2-(c1+c2)*.07)
	path.LineTo((s2-(s1-s2)*.07)*-1, -c2-(c1+c2)*.07)
	path.Close()

	// bottom right corner
	path.MoveTo(s2p, -c2p)
	path.LineTo((-s2p-(-s2p+s1p)*.3)*-1, -c2p+(c1p+c2p)*.3)
	path.LineTo(s2*1.1, -c2*1.1)
	path.LineTo((-s2p+(s2p+s2p)*.3)*-1, -c2p)
	path.Close()

	// right edge
	path.MoveTo((-s1p+(-s2p+s1p)*.34)*-1, c1p-(c1p+c2p)*.34)
	path.LineTo((-s2p-(-s2p+s1p)*.34)*-1, -c2p+(c1p+c2p)*.34)
	path.LineTo((-s2-(2*s2)*.07)*-1, -c2)
	path.LineTo((-s1-(s1)*.07)*-1, c1-(1-c1)*.07)
	path.Close()

	// top right corner
	path.MoveTo(s1p, c1p)
	path.LineTo((-s1p+(s1p)*.3)*-1, c1p+(2-c1p)*.3)
	path.LineTo((-s1*1.1)*-1, c1*1.1)
	path.LineTo((-s1p+(-s2p+s1p)*.3)*-1, c1p-(c1p+c2p)*.3)
	path.Close()

	// top right edge
	path.MoveTo((-s1p+(s1p)*.34)*-1, c1p+(2-c1p)*.34)
	path.LineTo((0+(s1p+0)*-.34)*-1, 2-(2-c1p)*.34)
	path.LineTo(s1*.07, 1+(1-c1)*.07)
	path.LineTo((-s1-(-s2+s1)*.07)*-1, c1+(c1+c2)*.07)
	path.Close()

	// top corner
	path.MoveTo(0, 2)
	path.LineTo(0+(s1p+0)*.3, 2-(2-c1p)*.3)
	path.LineTo(0, 1.1)
	path.LineTo(0+(s1p+0)*-.3, 2-(2-c1p)*.3)
	path.Close()

	// top left edge
	path.MoveTo(-s1p+(s1p)*.34, c1p+(2-c1p)*.34)
	path.LineTo(0+(s1p+0)*-.34, 2-(2-c1p)*.34)
	path.LineTo(s1*-.07, 1+(1-c1)*.07)
	path.LineTo(-s1-(-s2+s1)*.07, c1+(c1+c2)*.07)
	path.Close()

	// top left corner
	path.MoveTo(-s1p, c1p)
	path.LineTo(-s1p+(s1p)*.3, c1p+(2-c1p)*.3)
	path.LineTo(-s1*1.1, c1*1.1)
	path.LineTo(-s1p+(-s2p+s1p)*.3, c1p-(c1p+c2p)*.3)
	path.Close()

	// left edge
	path.MoveTo(-s1p+(-s2p+s1p)*.34, c1p-(c1p+c2p)*.34)
	path.LineTo(-s2p-(-s2p+s1p)*.34, -c2p+(c1p+c2p)*.34)
	path.LineTo(-s2-(2*s2)*.07, -c2)
	path.LineTo(-s1-(s1)*.07, c1-(1-c1)*.07)
	path.Close()

	return path.AppendVerticesAndIndicesForFilling(nil, nil)
}

func drawFaces(screen *ebiten.Image, scale float32) {
	drawTopHalf(screen, scale)
	drawBottomHalf(screen, scale)
}

func drawTopHalf(screen *ebiten.Image, scale float32) {
	vs, is := getFacePath()

	scaleVertices(vs, scale)
	rotateVertices(vs, math.Pi)

	n := len(vs)

	middleVs := make([]ebiten.Vertex, n)
	copy(middleVs, vs)
	reflectVerticesOverY(middleVs)
	translateVertices(middleVs, screenWidth/2, screenHeight/2)

	tLeftVs := make([]ebiten.Vertex, n)
	copy(tLeftVs, vs)
	rotateVertices(tLeftVs, math.Pi*1.8)
	translateVertices(tLeftVs, screenWidth/2, screenHeight/2)
	translateVertices(tLeftVs, scale*-1.95, scale*-2.7)

	tRightVs := make([]ebiten.Vertex, n)
	copy(tRightVs, vs)
	rotateVertices(tRightVs, math.Pi*2.2)
	translateVertices(tRightVs, screenWidth/2, screenHeight/2)
	translateVertices(tRightVs, scale*1.95, scale*-2.7)

	bLeftVs := make([]ebiten.Vertex, n)
	copy(bLeftVs, vs)
	rotateVertices(bLeftVs, math.Pi*3.1)
	reflectVertices(bLeftVs)
	translateVertices(bLeftVs, screenWidth/2, screenHeight/2)
	translateVertices(bLeftVs, scale*-3.18, scale)

	bRightVs := make([]ebiten.Vertex, n)
	copy(bRightVs, vs)
	rotateVertices(bRightVs, math.Pi*.6)
	translateVertices(bRightVs, screenWidth/2, screenHeight/2)
	translateVertices(bRightVs, scale*3.18, scale)

	bottomVs := make([]ebiten.Vertex, n)
	copy(bottomVs, vs)
	rotateVertices(bottomVs, math.Pi)
	reflectVerticesOverY(bottomVs)
	translateVertices(bottomVs, screenWidth/2, screenHeight/2)
	translateVertices(bottomVs, 0, scale*3.33)

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true

	PaintFace(middleVs, 0)
	PaintFace(bottomVs, 1)
	PaintFace(bRightVs, 2)
	PaintFace(tRightVs, 3)
	PaintFace(tLeftVs, 4)
	PaintFace(bLeftVs, 5)

	translateVertices(middleVs, scale*6, 0)
	translateVertices(tLeftVs, scale*6, 0)
	translateVertices(tRightVs, scale*6, 0)
	translateVertices(bottomVs, scale*6, 0)
	translateVertices(bRightVs, scale*6, 0)
	translateVertices(bLeftVs, scale*6, 0)

	screen.DrawTriangles(middleVs, is, whiteImage, op)
	screen.DrawTriangles(tLeftVs, is, whiteImage, op)
	screen.DrawTriangles(tRightVs, is, whiteImage, op)
	screen.DrawTriangles(bLeftVs, is, whiteImage, op)
	screen.DrawTriangles(bRightVs, is, whiteImage, op)
	screen.DrawTriangles(bottomVs, is, whiteImage, op)
}

func drawBottomHalf(screen *ebiten.Image, scale float32) {
	vs, is := getFacePath()

	scaleVertices(vs, scale)
	rotateVertices(vs, math.Pi)

	n := len(vs)

	middleVs := make([]ebiten.Vertex, n)
	copy(middleVs, vs)
	reflectVerticesOverY(middleVs)
	translateVertices(middleVs, screenWidth/2, screenHeight/2)

	tLeftVs := make([]ebiten.Vertex, n)
	copy(tLeftVs, vs)
	rotateVertices(tLeftVs, math.Pi*1.8)
	translateVertices(tLeftVs, screenWidth/2, screenHeight/2)
	translateVertices(tLeftVs, scale*-1.95, scale*-2.7)

	tRightVs := make([]ebiten.Vertex, n)
	copy(tRightVs, vs)
	rotateVertices(tRightVs, math.Pi*2.2)
	translateVertices(tRightVs, screenWidth/2, screenHeight/2)
	translateVertices(tRightVs, scale*1.95, scale*-2.7)

	bLeftVs := make([]ebiten.Vertex, n)
	copy(bLeftVs, vs)
	rotateVertices(bLeftVs, math.Pi*3.1)
	reflectVertices(bLeftVs)
	translateVertices(bLeftVs, screenWidth/2, screenHeight/2)
	translateVertices(bLeftVs, scale*-3.18, scale)

	bRightVs := make([]ebiten.Vertex, n)
	copy(bRightVs, vs)
	rotateVertices(bRightVs, math.Pi*.6)
	translateVertices(bRightVs, screenWidth/2, screenHeight/2)
	translateVertices(bRightVs, scale*3.18, scale)

	bottomVs := make([]ebiten.Vertex, n)
	copy(bottomVs, vs)
	rotateVertices(bottomVs, math.Pi)
	reflectVerticesOverY(bottomVs)
	translateVertices(bottomVs, screenWidth/2, screenHeight/2)
	translateVertices(bottomVs, 0, scale*3.33)

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true

	PaintFace(middleVs, 6)
	PaintFace(bottomVs, 7)
	PaintFace(bRightVs, 8)
	PaintFace(tRightVs, 9)
	PaintFace(tLeftVs, 10)
	PaintFace(bLeftVs, 11)

	translateVertices(middleVs, scale*-6, 0)
	translateVertices(tLeftVs, scale*-6, 0)
	translateVertices(tRightVs, scale*-6, 0)
	translateVertices(bottomVs, scale*-6, 0)
	translateVertices(bRightVs, scale*-6, 0)
	translateVertices(bLeftVs, scale*-6, 0)

	screen.DrawTriangles(middleVs, is, whiteImage, op)
	screen.DrawTriangles(tLeftVs, is, whiteImage, op)
	screen.DrawTriangles(tRightVs, is, whiteImage, op)
	screen.DrawTriangles(bLeftVs, is, whiteImage, op)
	screen.DrawTriangles(bRightVs, is, whiteImage, op)
	screen.DrawTriangles(bottomVs, is, whiteImage, op)
}
