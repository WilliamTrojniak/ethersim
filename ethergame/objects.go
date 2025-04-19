package ethergame

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type GameObject interface {
	Draw(img *ebiten.Image, prog float32)
	OnEvent(Event) bool
}
