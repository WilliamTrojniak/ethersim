package ethergame

import (
	"image/color"
	"time"

	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TIME_PER_TICK = time.Millisecond * 100
)

type Game struct {
	prevTick        time.Time
	objs            []GameObject
	nodes           []*Node
	edges           []*Edge
	devices         []*Device
	sim             *ethersim.Simulation
	justPressedKeys []ebiten.Key
	paused          bool
	prog            float32
	activeWeight    int
}

func (g *Game) OnEvent(event Event) {
	switch e := event.(type) {
	case KeyJustPressedEvent:
		switch e.Key {
		case ebiten.KeySpace:
			g.paused = !g.paused
		case ebiten.Key1:
			g.activeWeight = 10
		case ebiten.Key2:
			g.activeWeight = 20
		case ebiten.Key3:
			g.activeWeight = 30
		case ebiten.Key4:
			g.activeWeight = 40
		case ebiten.Key5:
			g.activeWeight = 50
		case ebiten.Key6:
			g.activeWeight = 60
		case ebiten.Key7:
			g.activeWeight = 70
		case ebiten.Key8:
			g.activeWeight = 80
		case ebiten.Key9:
			g.activeWeight = 90

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
		g.OnEvent(MouseClickEvent{X: x, Y: y, Button: ebiten.MouseButtonLeft})
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.OnEvent(MouseReleaseEvent{X: x, Y: y, Button: ebiten.MouseButtonLeft})
	} else {
		x, y := ebiten.CursorPosition()
		g.OnEvent(MouseMoveEvent{X: x, Y: y})
	}

	g.justPressedKeys = inpututil.AppendJustPressedKeys(g.justPressedKeys[:0])
	for _, key := range g.justPressedKeys {
		g.OnEvent(KeyJustPressedEvent{Key: key})
	}

	t := time.Now()
	if g.paused {
		g.prevTick = t.Add(-(time.Duration(g.prog * float32(TIME_PER_TICK))))
	}

	if t.Sub(g.prevTick) <= 1*TIME_PER_TICK || g.paused {
		return nil
	}

	g.prevTick = t
	g.sim.Tick()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
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
	return outsideWidth, outsideHeight
}

func MakeGame(sim *ethersim.Simulation) *Game {
	return &Game{
		prevTick:        time.Now(),
		objs:            make([]GameObject, 0),
		nodes:           make([]*Node, 0),
		edges:           make([]*Edge, 0),
		devices:         make([]*Device, 0),
		sim:             sim,
		justPressedKeys: make([]ebiten.Key, 0, 10),
		paused:          false,
		activeWeight:    10,
	}
}
