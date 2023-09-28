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
	rotations    = 10
)

var (
	whiteImage = ebiten.NewImage(3, 3)
)

var selectors = make([][]ebiten.Vertex, 12)

var gui bool

func init() {
	guiFlag := flag.Bool("gui", false, "if gui flag is set, the gui will be displayed. otherwise, the solver will be ran.")
	flag.Parse()
	gui = *guiFlag
	rand.Seed(time.Now().UnixNano()) // seed random number generator for randomizer
	whiteImage.Fill(color.White)
	state = NewState()
}

var state State

type Game struct {
	frame    int
	selected int
	stack    []Node // stack of nodes to unwind. if len(stack) == 0, no nodes to unwind
}

func (g *Game) Update() error {
	g.frame = (g.frame + 1) % 60

	if g.frame%60 == 0 && len(g.stack) != 0 {
		state = *(g.stack[len(g.stack)-1]).s
		g.stack = g.stack[:len(g.stack)-1] // pop off last element in stack
	}

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
		state.CW(g.selected)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		state.CCW(g.selected)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		state = NewState()
		state.randomize(rotations)
		_, node := Solve(state)
		var stack []Node
		for {
			stack = append(stack, node)
			if node.prev == nil {
				break
			}
			node = *(node.prev)
		}
		g.stack = stack
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		state = NewState()
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
		frame:    0,
		selected: 0,
	}); err != nil {
		log.Fatal(err)
	}

}
