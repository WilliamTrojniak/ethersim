package ethergame

import (
	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/hajimehoshi/ebiten/v2"
)

type Device struct {
	*ethersim.NetworkDevice
	Graphic
	clicked  bool
	selected bool
}

func (s *Device) Draw(screen *ebiten.Image, prog float32) {
	if s.selected {
		s.Graphic.SetColor(ColorYellow)
	} else if s.NetworkDevice.IncomingMsg() {
		s.Graphic.SetColor(ColorPurple)
	} else {
		s.Graphic.SetColor(ColorPurple)
	}

	s.Graphic.Draw(screen, prog)
}

func (s *Device) OnEvent(e Event) bool {
	switch e := e.(type) {
	case MouseClickEvent:
		if s.In(e.X, e.Y) && e.Button == ebiten.MouseButtonLeft {
			s.clicked = true
			s.selected = true
			return false
		} else {
			s.selected = false
		}
		break
	case MouseMoveEvent:
		if s.clicked {
			s.MoveTo(e.X, e.Y)
			return false
		}
		break
	case MouseReleaseEvent:
		prev := s.clicked
		s.clicked = false
		return prev
	}

	return false
}

func MakeDevice(n *Node, w int) (*Device, *Edge) {
	simDevice, simEdge := n.CreateDevice(w)
	d := &Device{
		NetworkDevice: simDevice,
		Graphic: &Rect{
			pos: Vec2[int]{50, 50},
			W:   64,
			H:   64,
			c:   ColorDark,
		},
		clicked: false,
	}
	e := makeEdge(n, d, simEdge)

	return d, e
}
