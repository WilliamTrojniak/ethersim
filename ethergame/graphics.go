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

var (
	ColorDark   color.Color = color.RGBA{0x2D, 0x27, 0x27, 0xFF}
	ColorGrey   color.Color = color.RGBA{0x41, 0x35, 0x43, 0xFF}
	ColorPurple color.Color = color.RGBA{0x8f, 0x43, 0xee, 0xFF}

	// Vintage Color Palette
	ColorMaroon    color.Color = color.RGBA{0x8b, 0x1e, 0x3f, 0xFF}
	ColorSalmon    color.Color = color.RGBA{0xef, 0x76, 0x7a, 0xFF}
	ColorYellow    color.Color = color.RGBA{0xf5, 0xdd, 0x90, 0xFF}
	ColorTeal      color.Color = color.RGBA{0x49, 0xba, 0xaa, 0xFF}
	ColorCyan      color.Color = color.RGBA{0x77, 0x9f, 0xa1, 0xFF}
	ColorNavy      color.Color = color.RGBA{0x45, 0x69, 0x90, 0xFF}
	ColorFadedNavy color.Color = color.RGBA{0x76, 0x85, 0x94, 0xFF}
)

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
	x := float32(r.pos.X - r.W/2)
	y := float32(r.pos.Y - r.H/2)
	vector.DrawFilledRect(img, x, y, float32(r.W), float32(r.H), r.c, true)
	vector.StrokeRect(img, x, y, float32(r.W), float32(r.H), 2, color.Black, true)
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
	pos    Vec2[int]
	R      float32
	border bool
	c      color.Color
}

func (c *Circle) Pos() Vec2[int] { return c.pos }
func (c *Circle) MoveTo(x, y int) {
	c.pos.X = x
	c.pos.Y = y
}
func (c *Circle) Draw(img *ebiten.Image, prog float32) {
	x := float32(c.pos.X)
	y := float32(c.pos.Y)
	vector.DrawFilledCircle(img, x, y, c.R, c.c, true)
	if c.border {
		vector.StrokeCircle(img, x, y, c.R, 2, color.Black, true)
	}
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
