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
	n1 := n0.CreateNode(20)
	n2 := n1.CreateNode(20)
	n3 := n2.CreateNode(20)
	n4 := n3.CreateNode(20)
	d0 := n0.CreateDevice(10)
	d1 := n1.CreateDevice(10)
	d2 := n2.CreateDevice(10)
	d3 := n3.CreateDevice(10)
	d4 := n4.CreateDevice(10)
	n0.MoveTo(100, 50)
	n1.MoveTo(200, 50)
	n2.MoveTo(300, 50)
	n3.MoveTo(400, 50)
	n4.MoveTo(500, 50)
	d0.MoveTo(100, 100)
	d1.MoveTo(200, 100)
	d2.MoveTo(300, 100)
	d3.MoveTo(400, 100)
	d4.MoveTo(500, 100)

	d0.QueueMessage(&ethersim.BaseMsg{V: true})
	d1.QueueMessage(&ethersim.BaseMsg{V: true})

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello World")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
