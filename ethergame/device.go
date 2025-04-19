package ethergame

import (
	"image/color"

	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/hajimehoshi/ebiten/v2"
)

type Device struct {
	*ethersim.NetworkDevice
	Graphic
	clicked bool
}

func (s *Device) Update() error     { return nil }
func (s *Device) PostUpdate() error { return nil }
func (s *Device) OnEvent(e Event) bool {
	switch e := e.(type) {
	case MouseClickEvent:
		if s.In(e.X, e.Y) && e.Button == ebiten.MouseButtonLeft {
			s.clicked = true
			return true
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

func MakeDevice(n *Node, w int) *Device {
	d := &Device{}
	simDevice, simEdge := n.CreateDevice(w)
	d.NetworkDevice = simDevice

	d.Graphic = &Composite{
		Graphic: &Rect{
			pos: Vec2[int]{50, 50},
			W:   64,
			H:   64,
			c:   color.RGBA{0xFF, 0x00, 0x00, 0xFF},
		},
		secondary: &Edge{
			n1:   n,
			n2:   d,
			edge: simEdge,
			c:    color.Black,
		},
	}
	return d
}
