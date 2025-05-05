package ethergame

import (
	"fmt"
	"math/rand"

	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/hajimehoshi/ebiten/v2"
)

type Device struct {
	game *Game
	*ethersim.NetworkDevice
	Circle
	clicked  bool
	selected bool
}

func (s *Device) Draw(screen *ebiten.Image, prog float32) {
	if s.selected {
		s.SetColor(ColorTeal)
	} else if s.NetworkDevice.IncomingMsg() {
		s.SetColor(ColorFadedNavy)
	} else {
		s.SetColor(ColorNavy)
	}

	p := Progress{
		C: &s.Circle,
		p: float32(s.Timeout()) / float32(s.TimeoutFrom()),
	}
	p.Draw(screen, prog)
	s.Circle.Draw(screen, prog)
}

func (s *Device) OnEvent(e Event) bool {
	switch e := e.(type) {
	case MouseClickEvent:
		if s.Circle.In(e.X, e.Y) && e.Button == ebiten.MouseButtonLeft {
			s.clicked = !s.clicked
			s.selected = !s.selected
			return false
		} else {
			s.selected = false
		}
		break
	case MouseMoveEvent:
		if s.clicked {
			s.Circle.MoveTo(e.X, e.Y)
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
			s.QueueMessage(&ethersim.BaseMsg{V: true, Msg: fmt.Sprintf("%v", rand.Intn(10)), To: rand.Intn(len(s.game.devices))})
			return true
		}
	}

	return false
}
func (d *Device) Update() {}
func (n *Node) CreateDevice(w int) *Device {
	simDevice, simEdge := n.NetworkNode.CreateDevice(w)
	d := &Device{
		game:          n.game,
		NetworkDevice: simDevice,
		Circle: Circle{
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
