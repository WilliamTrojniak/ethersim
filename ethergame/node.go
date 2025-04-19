package ethergame

import (
	"image/color"

	"github.com/WilliamTrojniak/ethersim/ethersim"
)

type Node struct {
	*ethersim.NetworkNode
	Graphic
}

func (n *Node) OnEvent(e Event) bool { return false }

func MakeNode(s *ethersim.Simulation) *Node {
	return &Node{
		NetworkNode: ethersim.MakeNetworkNode(s),
		Graphic: &Circle{
			pos: Vec2[int]{50, 50},
			R:   4,
			c:   color.RGBA{0x00, 0x00, 0x00, 0xFF},
		},
	}
}

func (n *Node) CreateNode(w int) *Node {
	nn, edge := n.NetworkNode.CreateNode(w)

	out := &Node{}
	out.NetworkNode = nn
	out.Graphic = &Composite{
		Graphic: &Circle{
			pos: Vec2[int]{n.Pos().X + 64, 50},
			R:   4,
			c:   color.RGBA{0x00, 0x00, 0x00, 0xFF},
		},
		secondary: &Edge{
			n1:   n,
			n2:   out,
			edge: edge,
			c:    color.Black,
		},
	}
	return out
}
