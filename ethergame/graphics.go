package ethergame

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
	vector.DrawFilledRect(img, float32(r.pos.X-r.W/2), float32(r.pos.Y-r.H/2), float32(r.W), float32(r.H), r.c, true)
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

type Composite struct {
	Graphic
	secondary Graphic
}

func (c *Composite) Draw(img *ebiten.Image, prog float32) {
	c.secondary.Draw(img, prog)
	c.Graphic.Draw(img, prog)
}
