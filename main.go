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
	prevTick        time.Time
	objs            []ethergame.GameObject
	nodes           []*ethergame.Node
	edges           []*ethergame.Edge
	devices         []*ethergame.Device
	sim             *ethersim.Simulation
	justPressedKeys []ebiten.Key
	paused          bool
	prog            float32
}

func (g *Game) OnEvent(event ethergame.Event) {
	switch e := event.(type) {
	case ethergame.KeyJustPressedEvent:
		if e.Key == ebiten.KeySpace {
			g.paused = !g.paused
		}
	}

	for _, obj := range g.objs {
		if obj.OnEvent(event) {
			break
		}
	}
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.OnEvent(ethergame.MouseClickEvent{X: x, Y: y, Button: ebiten.MouseButtonLeft})
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.OnEvent(ethergame.MouseReleaseEvent{X: x, Y: y, Button: ebiten.MouseButtonLeft})
	} else {
		x, y := ebiten.CursorPosition()
		g.OnEvent(ethergame.MouseMoveEvent{X: x, Y: y})
	}

	g.justPressedKeys = inpututil.AppendJustPressedKeys(g.justPressedKeys[:0])
	for _, key := range g.justPressedKeys {
		fmt.Printf("Key pressed %v\b", key.String())
		g.OnEvent(ethergame.KeyJustPressedEvent{Key: key})
	}

	t := time.Now()
	if g.paused {
		g.prevTick = t.Add(-(time.Duration(g.prog * float32(TIME_PER_TICK))))
	}

	if t.Sub(g.prevTick) <= 1*time.Second || g.paused {
		return nil
	}

	g.prevTick = t
	g.sim.Tick()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xFF, 0xFF, 0x00, 0xFF})
	n := time.Now()
	deltaT := n.Sub(g.prevTick)
	if !g.paused {
		g.prog = min(1, float32(deltaT)/float32(TIME_PER_TICK))
	}

	for _, edge := range g.edges {
		edge.Draw(screen, g.prog)
	}

	for _, node := range g.nodes {
		node.Draw(screen, g.prog)
	}

	for _, dev := range g.devices {
		dev.Draw(screen, g.prog)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func MakeGame(sim *ethersim.Simulation) *Game {
	return &Game{
		prevTick:        time.Now(),
		objs:            make([]ethergame.GameObject, 0),
		nodes:           make([]*ethergame.Node, 0),
		edges:           make([]*ethergame.Edge, 0),
		devices:         make([]*ethergame.Device, 0),
		sim:             sim,
		justPressedKeys: make([]ebiten.Key, 0, 10),
		paused:          false,
	}
}

func main() {
	sim := ethersim.MakeSimulation()
	game := MakeGame(sim)
	n0 := ethergame.MakeNode(sim)
	n1, e0 := n0.CreateNode(4)
	d0, e1 := ethergame.MakeDevice(n0, 2)
	d1, e2 := ethergame.MakeDevice(n0, 3)
	d2, e3 := ethergame.MakeDevice(n1, 2)
	d3, e4 := ethergame.MakeDevice(n1, 2)
	d4, e5 := ethergame.MakeDevice(n1, 2)
	n0.MoveTo(300, 50)
	n1.MoveTo(300, 300)
	d1.MoveTo(600, 50)
	d2.MoveTo(50, 300)
	d3.MoveTo(600, 300)
	d4.MoveTo(300, 400)

	game.nodes = append(game.nodes, n0, n1)
	game.edges = append(game.edges, e0, e1, e2, e3, e4, e5)
	game.devices = append(game.devices, d0, d1, d2, d3, d4)
	game.objs = append(game.objs, n0, n1, e0, e1, e2, e3, e4, e5, d0, d1, d2, d3, d4)

	d0.QueueMessage(&ethersim.BaseMsg{V: true})
	d1.QueueMessage(&ethersim.BaseMsg{V: true})

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello World")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
