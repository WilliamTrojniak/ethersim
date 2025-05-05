package ethergame

import (
	"fmt"
	"image/color"

	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Device struct {
	*ethersim.NetworkDevice
	Circle
	clicked  bool
	selected bool
	ui       *widget.Text
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
			s.QueueMessage(&ethersim.BaseMsg{V: true})
			return true
		}
	}

	return false
}

func (d *Device) getLabel() string {
	return fmt.Sprintf("(D%v) Queued: %v, Timeout Range %v", d.Id(), len(d.QueuedMessages()), d.TimeoutRange())
}

func (d *Device) createUI() *widget.Text {
	row := widget.NewText(widget.TextOpts.Text(
		d.getLabel(),
		face,
		color.Black,
	))
	return row
}

func (d *Device) Update() {
	d.ui.Label = d.getLabel()
}

func (n *Node) CreateDevice(w int) *Device {
	simDevice, simEdge := n.NetworkNode.CreateDevice(w)
	d := &Device{
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

	d.ui = d.createUI()
	n.game.deviceDataContainer.AddChild(d.ui)

	return d
}
