package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/WilliamTrojniak/ethersim/ethergame"
	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TIME_PER_TICK = 1 * time.Second
)

type Game struct {
	prevTick time.Time
	objs     []ethergame.GameObject
	sim      *ethersim.Simulation
}

func (g *Game) OnEvent(event ethergame.Event) {
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
		g.OnEvent(ethergame.MouseClickEvent{X: x, Y: y, Button: ebiten.MouseButtonLeft})
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fmt.Printf("Release at %v, %v\n", x, y)
		g.OnEvent(ethergame.MouseReleaseEvent{X: x, Y: y, Button: ebiten.MouseButtonLeft})
	} else {
		x, y := ebiten.CursorPosition()
		g.OnEvent(ethergame.MouseMoveEvent{X: x, Y: y})
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
	n0 := ethergame.MakeNode(sim)
	n1 := n0.CreateNode(4)
	d0 := ethergame.MakeDevice(n0, 2)
	d1 := ethergame.MakeDevice(n0, 3)
	d2 := ethergame.MakeDevice(n1, 2)
	d3 := ethergame.MakeDevice(n1, 2)
	d4 := ethergame.MakeDevice(n1, 2)
	n0.MoveTo(300, 50)
	n1.MoveTo(300, 300)
	d1.MoveTo(600, 50)
	d2.MoveTo(50, 300)
	d3.MoveTo(600, 300)
	d4.MoveTo(300, 400)

	d0.QueueMessage(&ethersim.BaseMsg{V: true})
	d1.QueueMessage(&ethersim.BaseMsg{V: true})

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello World")
	if err := ebiten.RunGame(&Game{sim: sim, prevTick: time.Now(), objs: []ethergame.GameObject{n0, n1, d0, d1, d2, d3, d4}}); err != nil {
		log.Fatal(err)
	}

}
