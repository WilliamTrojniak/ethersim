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
	for i := range 10 {
		node.MoveTo(50+i*60, 50)
		d := node.CreateDevice(3)
		d.MoveTo(50+i*60, 100)

		if i%4 == 0 {
			d.QueueMessage(&ethersim.BaseMsg{V: true})
			d.QueueMessage(&ethersim.BaseMsg{V: true})
			d.QueueMessage(&ethersim.BaseMsg{V: true})
			d.QueueMessage(&ethersim.BaseMsg{V: true})
		}

		if i < 9 {
			node = node.CreateNode(5)
		}
	}

	ebiten.SetWindowSize(1200, 800)
	ebiten.SetWindowTitle("Hello World")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
