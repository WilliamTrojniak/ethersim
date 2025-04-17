package ethersim

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Vec2[T any] struct {
	X T
	Y T
}

type Server struct {
	*Device
	Graphic
	clicked bool
}

func (s *Server) Update() error     { return nil }
func (s *Server) PostUpdate() error { return nil }
func (s *Server) Draw(img *ebiten.Image, prog float32) {
	if s.network.IncomingMsg(s.Device) {
		s.SetColor(color.Black)
	} else {
		s.SetColor(color.RGBA{0xA9, 0xAF, 0xD1, 0xFF})
	}

	s.Graphic.Draw(img, prog)
}

func (s *Server) OnEvent(e Event) bool {
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

func MakeServer(n *Node, w int) Server {
	s := Server{}
	s.Device = n.createDevice(w)
	edge := n.edges[len(n.edges)-1]

	s.Graphic = &Composite{
		Graphic: &Rect{
			pos: Vec2[int]{50, 50},
			W:   64,
			H:   64,
			c:   color.RGBA{0xFF, 0x00, 0x00, 0xFF},
		},
		secondary: &Edge{
			n1:   n,
			n2:   &s,
			edge: edge,
			c:    color.Black,
		},
	}
	return s
}

type Node struct {
	*NetworkNode
	Graphic
}

func (n *Node) Update() error {
	return n.NetworkNode.Update()

}

func (n *Node) Draw(img *ebiten.Image, prog float32) {
	n.Graphic.Draw(img, prog)
}

func (n *Node) OnEvent(e Event) bool {
	return false
}

func MakeNode() Node {
	return Node{
		MakeNetworkNode(),
		&Circle{
			pos: Vec2[int]{50, 50},
			R:   4,
			c:   color.RGBA{0x00, 0x00, 0x00, 0xFF},
		},
	}
}

func (n *Node) CreateNode(w int) Node {
	nn := n.NetworkNode.CreateNode(w)
	edge := n.edges[len(n.edges)-1]

	out := Node{}
	out.NetworkNode = nn
	out.Graphic = &Composite{
		Graphic: &Circle{
			pos: Vec2[int]{n.Pos().X + 64, 50},
			R:   4,
			c:   color.RGBA{0x00, 0x00, 0x00, 0xFF},
		},
		secondary: &Edge{
			n1:   n,
			n2:   &out,
			edge: edge,
			c:    color.Black,
		},
	}
	return out
}
