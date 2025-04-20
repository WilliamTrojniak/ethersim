package ethergame

import (
	"image/color"

	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/hajimehoshi/ebiten/v2"
)

type Node struct {
	game *Game
	*ethersim.NetworkNode
	Graphic
	clicked  bool
	selected bool
}

func (n *Node) OnEvent(e Event) bool {
	switch e := e.(type) {
	case MouseClickEvent:
		if n.In(e.X, e.Y) && e.Button == ebiten.MouseButtonLeft {
			n.clicked = true
			n.selected = true
			return false
		} else {
			n.selected = false
		}
		break
	case MouseMoveEvent:
		if n.clicked {
			n.MoveTo(e.X, e.Y)
			return false
		}
		break
	case MouseReleaseEvent:
		n.clicked = false
		return false
	case KeyJustPressedEvent:
		if !n.selected {
			return false
		}
		switch e.Key {
		case ebiten.KeyN:
			n.CreateNode(n.game.activeWeight)
			break
		case ebiten.KeyD:
			n.CreateDevice(n.game.activeWeight)
			break
		}
		return false
	}

	return false
}

func (g *Game) makeNode(n *ethersim.NetworkNode) *Node {
	nn := &Node{
		game:        g,
		NetworkNode: n,
		Graphic: &Circle{
			pos: Vec2[int]{50, 50},
			R:   8,
			c:   color.Black,
		},
		clicked:  false,
		selected: false,
	}
	g.nodes = append(g.nodes, nn)
	g.objs = append(g.objs, nn)
	return nn
}

func (g *Game) MakeNode(s *ethersim.Simulation) *Node {
	return g.makeNode(ethersim.MakeNetworkNode(s))
}

func (n *Node) CreateNode(w int) *Node {
	simNode, simEdge := n.NetworkNode.CreateNode(w)
	nn := n.game.makeNode(simNode)
	n.game.makeEdge(n, nn, simEdge)
	return nn
}
