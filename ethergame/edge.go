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

func (e *Edge) Draw(img *ebiten.Image, prog float32) {

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
		col := color.RGBA{0, 0, 0, 0xFF}
		if !msg.Msg().Valid() {
			col = color.RGBA{0xFF, 0, 0, 0xFF}
		}

		c := Circle{
			pos: Vec2[int]{
				X: x1 + int(totalprog*dx),
				Y: y1 + int(totalprog*dy),
			},
			R: 8,
			c: col,
		}
		c.Draw(img, prog)
	}

}
func (e *Edge) Pos() Vec2[int] {
	p := Vec2[int]{}
	p.X = (e.n1.Pos().X + e.n2.Pos().X) / 2
	p.Y = (e.n1.Pos().Y + e.n2.Pos().Y) / 2
	return p
}
func (e *Edge) In(x, y int) bool         { return false } // not selectable
func (e *Edge) MoveTo(x, y int)          {}               // no-op
func (e *Edge) SetColor(col color.Color) { e.c = col }
