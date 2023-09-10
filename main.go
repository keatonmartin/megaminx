package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
	"math/rand"
	"time"
)

const (
	screenWidth  = 450
	screenHeight = 300
)

var (
	whiteImage = ebiten.NewImage(3, 3)
)

var megaminx Megaminx

var selectors = make([][]ebiten.Vertex, 12)

var rotations int

func init() {

	rot := flag.Int("rotations", 20, "the number of clockwise rotations made when randomizing the puzzle")
	flag.Parse()
	rotations = *rot
	rand.Seed(time.Now().UnixNano()) // seed random number generator for randomizer
	whiteImage.Fill(color.White)
	megaminx = NewMegaminx()
}

type Game struct {
	m        Megaminx
	frame    int
	selected int
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) { // if left click just pressed
		// figure out where
		xi, yi := ebiten.CursorPosition()
		x := float32(xi)
		y := float32(yi)
		for i, s := range selectors {
			if x >= s[0].DstX && x <= s[1].DstX && y >= s[0].DstY && y <= s[2].DstY {
				g.selected = i
				break
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		megaminx.CW(g.selected)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		megaminx.CCW(g.selected)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		megaminx = NewMegaminx()
		randomize(rotations)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		megaminx = NewMegaminx()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawFaces(screen, float32(18))
	drawSelectors(screen)
	drawMarker(screen, g.selected)

	ebitenutil.DebugPrintAt(screen, "To restart, press R", 5, 260)
	ebitenutil.DebugPrintAt(screen, "To restart and randomize, press T", 5, 280)

	ebitenutil.DebugPrintAt(screen, "Clockwise: left arrow", 250, 280)
	ebitenutil.DebugPrintAt(screen, "Counter-clockwise: right arrow", 250, 260)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Megaminx Viewer")
	if err := ebiten.RunGame(&Game{
		m:        NewMegaminx(),
		frame:    0,
		selected: 0,
	}); err != nil {
		log.Fatal(err)
	}
}
