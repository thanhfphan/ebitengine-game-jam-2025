package game

import "github.com/hajimehoshi/ebiten/v2"

// Scene represents a game scene (menu, gameplay, settings, etc.)
type Scene interface {
    Enter(g *Game)
    Exit(g *Game)
    Update(g *Game)
    Draw(screen *ebiten.Image, g *Game)
}