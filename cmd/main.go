package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thanhfphan/ebitengj2025/internal/game"
)

func main() {
	game, err := game.New()
	if err != nil {
		log.Fatalf("Game init error: %v", err)
	}

	ebiten.SetWindowTitle("Food Cards")
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}
