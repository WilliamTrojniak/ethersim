package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/willtrojniak/ethersim/ethergame"
	"github.com/willtrojniak/ethersim/ethersim"
)

func main() {
	sim := ethersim.MakeSimulation()
	game := ethergame.MakeGame(sim)
	baseX := 550
	baseY := 200

	node := game.MakeNode(sim)
	count := 5
	for i := range count {
		node.MoveTo(baseX+50+i*60, baseY+50)
		d := node.CreateDevice(4)
		d.MoveTo(baseX+50+i*60, baseY+100)

		if i%5 == 0 {
			// d.QueueMessage(&ethersim.BaseMsg{V: true, Msg: "Hello", To: i + 3})
		}

		if i < count-1 {
			node = node.CreateNode(4)
		}
	}

	ebiten.SetWindowSize(1400, 800)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Hello World")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
