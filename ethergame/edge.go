package ethergame

import (
	"image/color"

	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Edge struct {
	n1   Graphic
	n2   Graphic
	edge *ethersim.NetworkEdge
	c    color.Color
}

func (g *Game) makeEdge(n1 Graphic, n2 Graphic, edge *ethersim.NetworkEdge) *Edge {
	e := &Edge{
		n1:   n1,
		n2:   n2,
		edge: edge,
		c:    ColorDark,
	}
	g.edges = append(g.edges, e)
	g.objs = append(g.objs, e)
	return e
}

func (e *Edge) Update() {}

func (e *Edge) Draw(img *ebiten.Image, prog float32) {
	if e.edge.IsResetting() {
		e.c = ColorSalmon
	} else {
		e.c = ColorDark
	}

	vector.StrokeLine(img, float32(e.n1.Pos().X), float32(e.n1.Pos().Y),
		float32(e.n2.Pos().X), float32(e.n2.Pos().Y), 4, e.c, true)

	x1 := e.n1.Pos().X
	y1 := e.n1.Pos().Y
	x2 := e.n2.Pos().X
	y2 := e.n2.Pos().Y
	dx := float32(x2 - x1)
	dy := float32(y2 - y1)

	w := 1.0 / float32(e.edge.Weight()) * prog

	if prog > 0.5 {
		dirs := make(map[int]int)
		for _, msg := range e.edge.Messages() {
			if v, ok := dirs[msg.Stage()]; !ok {
				dirs[msg.Stage()] = msg.Dir()
			} else if v != msg.Dir() {
				dirs[msg.Stage()] = 0
			}
		}

		for _, msg := range e.edge.Messages() {
			if !msg.Msg().Valid() {
				continue
			}
			// Case 1: Two messages will swap stages
			if v, ok := dirs[msg.Stage()+msg.Dir()]; ok && v != msg.Dir() {
				msg.Msg().Invalid()
			}
		}
	}

	for _, msg := range e.edge.Messages() {
		tickprog := float32(msg.Stage()) / float32(e.edge.Weight())
		totalprog := tickprog + w*float32(msg.Dir())
		var col color.Color
		col = ColorDark
		if !msg.Msg().Valid() {
			col = ColorSalmon
		}

		c := Circle{
			pos: Vec2[int]{
				X: x1 + int(totalprog*dx),
				Y: y1 + int(totalprog*dy),
			},
			c: col,
			R: 6,
		}
		c.Draw(img, prog)

		// w := Wave{
		// 	startPos: Vec2[int]{
		// 		X: x1 + int(totalprog*dx),
		// 		Y: y1 + int(totalprog*dy),
		// 	},
		// 	endPos: Vec2[int]{
		// 		X: x1 + int((totalprog+1/float32(e.edge.Weight()))*dx),
		// 		Y: y1 + int((totalprog+1/float32(e.edge.Weight()))*dy),
		// 	},
		// 	amplitude: 16,
		// }
		// w.Draw(img, prog)
	}

}
func (e *Edge) Pos() Vec2[int] {
	p := Vec2[int]{}
	p.X = (e.n1.Pos().X + e.n2.Pos().X) / 2
	p.Y = (e.n1.Pos().Y + e.n2.Pos().Y) / 2
	return p
}

func (e *Edge) SetColor(col color.Color) { e.c = col }
func (e *Edge) OnEvent(msg Event) bool   { return false }
