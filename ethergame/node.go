package ethergame

import (
	"image/color"

	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/hajimehoshi/ebiten/v2"
)

type Node struct {
	*ethersim.NetworkNode
	Graphic
	clicked bool
}

func (n *Node) OnEvent(e Event) bool {
	switch e := e.(type) {
	case MouseClickEvent:
		if n.In(e.X, e.Y) && e.Button == ebiten.MouseButtonLeft {
			n.clicked = true
			return true
		}
		break
	case MouseMoveEvent:
		if n.clicked {
			n.MoveTo(e.X, e.Y)
			return false
		}
		break
	case MouseReleaseEvent:
		prev := n.clicked
		n.clicked = false
		return prev
	}

	return false
}

func makeNode(n *ethersim.NetworkNode) *Node {
	return &Node{
		NetworkNode: n,
		Graphic: &Circle{
			pos: Vec2[int]{50, 50},
			R:   8,
			c:   color.Black,
		},
		clicked: false,
	}
}

func MakeNode(s *ethersim.Simulation) *Node {
	return makeNode(ethersim.MakeNetworkNode(s))
}

func (n *Node) CreateNode(w int) (*Node, *Edge) {
	simNode, simEdge := n.NetworkNode.CreateNode(w)
	nn := makeNode(simNode)
	e := makeEdge(n, nn, simEdge)
	return nn, e
}
