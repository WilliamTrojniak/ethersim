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
	n0 := game.MakeNode(sim)
	n1 := n0.CreateNode(4)
	d0 := n0.CreateDevice(2)
	d1 := n0.CreateDevice(3)
	d2 := n1.CreateDevice(2)
	d3 := n1.CreateDevice(2)
	d4 := n1.CreateDevice(2)
	n0.MoveTo(300, 50)
	n1.MoveTo(300, 300)
	d1.MoveTo(600, 50)
	d2.MoveTo(50, 300)
	d3.MoveTo(600, 300)
	d4.MoveTo(300, 400)

	d0.QueueMessage(&ethersim.BaseMsg{V: true})
	d1.QueueMessage(&ethersim.BaseMsg{V: true})

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello World")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
