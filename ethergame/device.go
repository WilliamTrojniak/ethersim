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
		s.Graphic.SetColor(ColorTeal)
	} else if s.NetworkDevice.IncomingMsg() {
		s.Graphic.SetColor(ColorFadedNavy)
	} else {
		s.Graphic.SetColor(ColorNavy)
	}

	s.Graphic.Draw(screen, prog)
}

func (s *Device) OnEvent(e Event) bool {
	switch e := e.(type) {
	case MouseClickEvent:
		if s.In(e.X, e.Y) && e.Button == ebiten.MouseButtonLeft {
			s.clicked = !s.clicked
			s.selected = !s.selected
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
	case KeyJustPressedEvent:
		if !s.selected {
			return false
		}
		switch e.Key {
		case ebiten.KeyM:
			s.QueueMessage(&ethersim.BaseMsg{V: true})
			return true
		}
	}

	return false
}

func (n *Node) CreateDevice(w int) *Device {
	simDevice, simEdge := n.NetworkNode.CreateDevice(w)
	d := &Device{
		NetworkDevice: simDevice,
		Graphic: &Circle{
			pos:    Vec2[int]{50, 50},
			R:      24,
			c:      ColorDark,
			border: true,
		},
		clicked: false,
	}
	n.game.makeEdge(n, d, simEdge)
	n.game.devices = append(n.game.devices, d)
	n.game.objs = append(n.game.objs, d)

	return d
}
