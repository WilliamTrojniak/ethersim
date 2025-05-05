package ethergame

import (
	"fmt"
	"image/color"

	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Node struct {
	game *Game
	*ethersim.NetworkNode
	Circle
	clicked  bool
	selected bool
	ui       *widget.Text
}

func (n *Node) Draw(screen *ebiten.Image, prog float32) {
	if n.selected {
		n.SetColor(ColorTeal)
	} else if n.IsResetting() {
		n.SetColor(ColorOrange)
	} else if n.IsTransmitting() {
		n.SetColor(ColorGreen)
	} else {
		n.SetColor(color.Black)
	}

	p := Progress{C: &n.Circle, p: float32(n.Timeout()) / float32(n.TimeoutFrom())}
	p.Draw(screen, prog)
	n.Circle.Draw(screen, prog)
}

func (n *Node) OnEvent(e Event) bool {
	switch e := e.(type) {
	case MouseClickEvent:
		if n.In(e.X, e.Y) && e.Button == ebiten.MouseButtonLeft {
			n.clicked = !n.clicked
			n.selected = !n.selected
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
			nn := n.CreateNode(n.game.activeWeight)
			n.selected = false
			nn.clicked = true
			nn.selected = true
			return true
		case ebiten.KeyD:
			d := n.CreateDevice(n.game.activeWeight)
			n.selected = false
			d.clicked = true
			d.selected = true
			return true
		}
		return false
	}

	return false
}

func (g *Game) makeNode(n *ethersim.NetworkNode) *Node {
	nn := &Node{
		game:        g,
		NetworkNode: n,
		Circle: Circle{
			pos:    Vec2[int]{50, 50},
			R:      8,
			c:      color.Black,
			border: true,
		},
		clicked:  false,
		selected: false,
	}
	nn.ui = nn.createUI()
	g.deviceDataContainer.AddChild(nn.ui)
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
func (n *Node) getLabel() string {
	return fmt.Sprintf("(T%v) | Max Timeout: %v | Queued: %v | Sending: %v | To: %v", n.Id(), n.TimeoutRange(), n.NQueued(), n.SendingValue(), n.SendingTo())
}

func (n *Node) createUI() *widget.Text {
	row := widget.NewText(widget.TextOpts.Text(
		n.getLabel(),
		face,
		color.Black,
	))
	return row
}

func (n *Node) Update() {
	n.ui.Label = n.getLabel()
}
