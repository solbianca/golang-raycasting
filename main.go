package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"shadow2d/shadow2d"
)

func main() {
	game := shadow2d.NewGame()
	ebiten.SetWindowSize(shadow2d.ScreenWidth, shadow2d.ScreenHeight)
	ebiten.SetWindowTitle("Light and Shadow in 2D (Ebiten Demo)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
