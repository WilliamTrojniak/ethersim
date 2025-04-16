package ethersim

import "github.com/hajimehoshi/ebiten/v2"

type Event interface{}

type MouseClickEvent struct {
	X, Y   int
	Button ebiten.MouseButton
}
type MouseReleaseEvent struct {
	X, Y   int
	Button ebiten.MouseButton
}

type MouseMoveEvent struct {
	X, Y int
}
