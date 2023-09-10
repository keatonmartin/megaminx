package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
)

const (
	screenWidth  = 450
	screenHeight = 300
)

var (
	whiteImage = ebiten.NewImage(3, 3)
	// whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

var megaminx Megaminx

var selectors = make([][]ebiten.Vertex, 12)

func init() {
	whiteImage.Fill(color.White)
	megaminx = NewMegaminx()
	//megaminx.CCW(3)
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

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		randomize(20)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawFaces(screen, float32(18))
	drawSelectors(screen)
	drawMarker(screen, g.selected)
	ebitenutil.DebugPrint(screen, "Megaminx Viewer")
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
