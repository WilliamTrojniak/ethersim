package main

import (
	"log"

	"github.com/WilliamTrojniak/ethersim/ethergame"
	"github.com/WilliamTrojniak/ethersim/ethersim"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	sim := ethersim.MakeSimulation()
	game := ethergame.MakeGame(sim)

	node := game.MakeNode(sim)
	count := 5
	for i := range count {
		node.MoveTo(50+i*60, 50)
		d := node.CreateDevice(4)
		d.MoveTo(50+i*60, 100)

		if i%5 == 0 {
			d.QueueMessage(&ethersim.BaseMsg{V: true})
		}

		if i < count-1 {
			node = node.CreateNode(4)
		}
	}

	ebiten.SetWindowSize(1200, 800)
	ebiten.SetWindowTitle("Hello World")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
