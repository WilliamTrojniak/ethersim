package ethersim

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Graphic interface {
	Draw(screen *ebiten.Image, prog float32)
	Pos() Vec2[int]
	In(x, y int) bool
	MoveTo(x, y int)
	SetColor(c color.Color)
}

type Rect struct {
	pos Vec2[int]
	W   int
	H   int
	c   color.Color
}

func (r *Rect) Pos() Vec2[int] {
	return r.pos
}

func (r *Rect) MoveTo(x, y int) {
	r.pos.X = x
	r.pos.Y = y
}

func (r *Rect) Draw(img *ebiten.Image, prog float32) {
	vector.DrawFilledRect(img, float32(r.pos.X-r.W/2), float32(r.pos.Y-r.H/2), float32(r.W), float32(r.H), r.c, false)
}

func (r *Rect) In(x, y int) bool {
	hW := r.W / 2
	hH := r.H / 2
	return x >= r.pos.X-hW && x <= r.pos.X+hW && y >= r.pos.Y-hH && y <= r.pos.Y+hH
}

func (r *Rect) SetColor(c color.Color) {
	r.c = c
}

type Circle struct {
	pos Vec2[int]
	R   float32
	c   color.Color
}

func (c *Circle) Pos() Vec2[int] { return c.pos }
func (c *Circle) MoveTo(x, y int) {
	c.pos.X = x
	c.pos.Y = y
}
func (c *Circle) Draw(img *ebiten.Image, prog float32) {
	vector.DrawFilledCircle(img, float32(c.pos.X), float32(c.pos.Y), c.R, c.c, true)
}
func (c *Circle) In(x, y int) bool {
	dx := x - c.pos.X
	dy := y - c.pos.Y
	d2 := dx*dx + dy*dy
	return d2 <= int(c.R*c.R)
}
func (c *Circle) SetColor(col color.Color) {
	c.c = col
}

type Edge struct {
	n1   Graphic
	n2   Graphic
	edge *NetworkEdge
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

	w := 1.0 / float32(e.edge.weight) * prog

	if prog > 0.5 {
		dirs := make(map[int]int)
		for _, msg := range e.edge.messages {
			if v, ok := dirs[msg.stage]; !ok {
				dirs[msg.stage] = msg.dir
			} else if v != msg.dir {
				dirs[msg.stage] = 0
			}
		}

		for _, msg := range e.edge.messages {
			if !msg.msg.Valid() {
				continue
			}
			// Case 1: Two messages will swap stages
			if v, ok := dirs[msg.stage+msg.dir]; ok && v != msg.dir {
				msg.msg.Invalid()
			}
		}
	}

	for _, msg := range e.edge.messages {
		tickprog := float32(msg.stage) / float32(e.edge.weight)
		totalprog := tickprog + w*float32(msg.dir)
		col := color.RGBA{0, 0, 0, 0xFF}
		if !msg.msg.Valid() {
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

type Composite struct {
	Graphic
	secondary Graphic
}

func (c *Composite) Draw(img *ebiten.Image, prog float32) {
	c.secondary.Draw(img, prog)
	c.Graphic.Draw(img, prog)
}
