package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TIME_PER_TICK = 1 * time.Second
)

type GameObj interface {
	Draw(img *ebiten.Image, prog float32)
	OnEvent(ethersim.Event) bool
}

type Game struct {
	prevTick time.Time
	objs     []GameObj
	sim      *ethersim.Simulation
}

func (g *Game) OnEvent(event ethersim.Event) {
	for _, obj := range g.objs {
		if obj.OnEvent(event) {
			break
		}
	}
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fmt.Printf("Click at %v, %v\n", x, y)
		g.OnEvent(ethersim.MouseClickEvent{X: x, Y: y, Button: ebiten.MouseButtonLeft})
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fmt.Printf("Release at %v, %v\n", x, y)
		g.OnEvent(ethersim.MouseReleaseEvent{X: x, Y: y, Button: ebiten.MouseButtonLeft})
	} else {
		x, y := ebiten.CursorPosition()
		g.OnEvent(ethersim.MouseMoveEvent{X: x, Y: y})
	}

	t := time.Now()
	if t.Sub(g.prevTick) <= 1*time.Second {
		return nil
	}
	g.prevTick = t
	g.sim.Tick()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
	n := time.Now()
	deltaT := n.Sub(g.prevTick)
	prog := float32(deltaT) / float32(TIME_PER_TICK)
	prog = min(1, prog)

	for _, obj := range g.objs {
		obj.Draw(screen, prog)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	sim := ethersim.MakeSimulation()
	n0 := ethersim.MakeNode(sim)
	n1 := n0.CreateNode(4)
	d0 := ethersim.MakeServer(&n0, 2)
	d1 := ethersim.MakeServer(&n0, 2)
	d2 := ethersim.MakeServer(&n1, 2)
	d3 := ethersim.MakeServer(&n1, 2)
	d4 := ethersim.MakeServer(&n1, 2)
	n0.MoveTo(300, 50)
	n1.MoveTo(300, 300)
	d1.MoveTo(600, 50)
	d2.MoveTo(50, 300)
	d3.MoveTo(600, 300)
	d4.MoveTo(300, 400)
	// ethersim.MakeDevice(2, n1)
	// ethersim.MakeDevice(3, n1)

	d0.QueueMessage(&ethersim.BaseMsg{V: true})
	d1.QueueMessage(&ethersim.BaseMsg{V: false})

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello World")
	if err := ebiten.RunGame(&Game{sim: sim, prevTick: time.Now(), objs: []GameObj{&n0, &n1, &d0, &d1, &d2, &d3, &d4}}); err != nil {
		log.Fatal(err)
	}

}
