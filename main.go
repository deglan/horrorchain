package main

import (
	"log"

	"github.com/deglan/horrorchain/constants"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(constants.WindowSizeWidth, constants.WindowSizeHeight)
	ebiten.SetWindowTitle("HORROR CHAIN!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
