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
	Update() error
	PostUpdate() error
	Draw(img *ebiten.Image, prog float32)
	OnEvent(ethersim.Event) bool
}

type Game struct {
	prevTick time.Time
	objs     []GameObj
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
	fmt.Printf("-------------tick-------------\n")

	for _, obj := range g.objs {
		obj.Update()
	}

	for _, obj := range g.objs {
		obj.PostUpdate()
	}

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
	n0 := ethersim.MakeNode()
	n1 := n0.CreateNode(2)
	d0 := ethersim.MakeServer(&n0, 2)
	d1 := ethersim.MakeServer(&n0, 3)
	d2 := ethersim.MakeServer(&n1, 1)
	d3 := ethersim.MakeServer(&n1, 1)
	n0.MoveTo(300, 50)
	n1.MoveTo(300, 300)
	d1.MoveTo(600, 50)
	d2.MoveTo(50, 300)
	d3.MoveTo(600, 300)
	// ethersim.MakeDevice(2, n1)
	// ethersim.MakeDevice(3, n1)

	d0.SendPacket(&ethersim.BaseMsg{V: true})
	// d1.SendPacket(&ethersim.BaseMsg{V: false})

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello World")
	if err := ebiten.RunGame(&Game{prevTick: time.Now(), objs: []GameObj{&n0, &n1, &d0, &d1, &d2, &d3}}); err != nil {
		log.Fatal(err)
	}

}
