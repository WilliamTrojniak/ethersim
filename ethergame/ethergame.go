package ethergame

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

var face, err = loadFont(18)

const (
	TIME_PER_TICK = time.Millisecond * 50
)

type Game struct {
	prevTick        time.Time
	objs            []GameObject
	nodes           []*Node
	edges           []*Edge
	devices         []*Device
	sim             *ethersim.Simulation
	justPressedKeys []ebiten.Key
	speedFactor     float32
	paused          bool
	prog            float32
	activeWeight    int
	ui              *ebitenui.UI

	deviceDataContainer *widget.Container
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

func (g *Game) OnEvent(event Event) {
	switch e := event.(type) {
	case KeyJustPressedEvent:
		switch e.Key {
		case ebiten.KeySpace:
			g.paused = !g.paused
		case ebiten.Key1:
			g.activeWeight = 1
		case ebiten.Key2:
			g.activeWeight = 2
		case ebiten.Key3:
			g.activeWeight = 3
		case ebiten.Key4:
			g.activeWeight = 4
		case ebiten.Key5:
			g.activeWeight = 5
		case ebiten.Key6:
			g.activeWeight = 6
		case ebiten.Key7:
			g.activeWeight = 7
		case ebiten.Key8:
			g.activeWeight = 8
		case ebiten.Key9:
			g.activeWeight = 9
		}
	}

	for _, obj := range g.objs {
		if obj.OnEvent(event) {
			return
		}
	}

	switch e := event.(type) {
	case KeyJustPressedEvent:
		switch e.Key {
		case ebiten.KeyN:
			nn := g.MakeNode(g.sim)
			nn.clicked = true
			nn.selected = true
			return
		}
	}
}

func (g *Game) Update() error {
	g.ui.Update()
	for _, obj := range g.objs {
		obj.Update()
	}

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
		g.prevTick = t.Add(-(time.Duration(g.prog * float32(TIME_PER_TICK) * g.speedFactor)))
	}

	if t.Sub(g.prevTick) <= time.Duration(float32(TIME_PER_TICK)/g.speedFactor) || g.paused {
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
		g.prog = min(1, float32(deltaT)/float32(TIME_PER_TICK)*g.speedFactor)
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

	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) getEbitenUI() *ebitenui.UI {

	root := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(16)))))
	footer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionVertical))),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			VerticalPosition:  widget.AnchorLayoutPositionEnd,
			StretchHorizontal: true,
		})))

	controlsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
	)

	deviceDataRows := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{Bottom: 16}))),
	)

	sliderContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionVertical))),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			VerticalPosition:   widget.AnchorLayoutPositionEnd,
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
		})))
	sliderLabel := widget.NewText(widget.TextOpts.Text(fmt.Sprintf("Speed: %.2fx", 100.0/100.0), face, color.Black))
	slider := widget.NewSlider(
		// Set the slider orientation - n/s vs e/w
		widget.SliderOpts.Direction(widget.DirectionHorizontal),
		// Set the minimum and maximum value for the slider
		widget.SliderOpts.MinMax(50, 200),
		// Set the current value of the slider, without triggering a change event
		widget.SliderOpts.InitialCurrent(100),
		widget.SliderOpts.WidgetOpts(
			// Set the Widget to layout in the center on the screen
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.SliderOpts.Images(
			// Set the track images
			&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(ColorFadedNavy),
				Hover: image.NewNineSliceColor(ColorFadedNavy),
			},
			// Set the handle images
			&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.Black),
				Hover:   image.NewNineSliceColor(color.Black),
				Pressed: image.NewNineSliceColor(color.Black),
			},
		),
		// Set the size of the handle
		widget.SliderOpts.FixedHandleSize(20),
		// Set the offset to display the track
		widget.SliderOpts.TrackOffset(0),
		// Set the size to move the handle
		widget.SliderOpts.PageSizeFunc(func() int {
			return 1
		}),
		// Set the callback to call when the slider value is changed
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			g.speedFactor = float32(args.Slider.Current) / 100.0
			sliderLabel.Label = fmt.Sprintf("Speed: %.2fx", g.speedFactor)
		}),

		widget.SliderOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
			}),

			// Set the widget's dimensions
			widget.WidgetOpts.MinSize(400, 10),
		),
	)

	controlsLabel := widget.NewText(widget.TextOpts.Text(
		"[space]: Pause/Play | [n]: New Transceiver | [d]: New Device\n[m]: New Message | [0-9]: Set New Edge Weight",
		face,
		color.Black,
	))

	root.AddChild(footer)
	footer.AddChild(deviceDataRows)
	footer.AddChild(controlsContainer)
	controlsContainer.AddChild(controlsLabel)
	controlsContainer.AddChild(sliderContainer)
	sliderContainer.AddChild(slider)
	sliderContainer.AddChild(sliderLabel)

	g.deviceDataContainer = deviceDataRows
	return &ebitenui.UI{
		Container: root,
	}

}

func MakeGame(sim *ethersim.Simulation) *Game {
	g := &Game{
		prevTick:        time.Now(),
		objs:            make([]GameObject, 0),
		nodes:           make([]*Node, 0),
		edges:           make([]*Edge, 0),
		devices:         make([]*Device, 0),
		sim:             sim,
		justPressedKeys: make([]ebiten.Key, 0, 10),
		paused:          false,
		activeWeight:    3,
		ui:              nil,
		speedFactor:     1.0,
	}

	g.ui = g.getEbitenUI()

	return g
}
