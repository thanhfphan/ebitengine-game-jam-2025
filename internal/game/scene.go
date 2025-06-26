package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thanhfphan/ebitengj2025/internal/ui"
)

// Scene represents a game scene (menu, gameplay, settings, etc.)
type Scene interface {
	Enter(g *Game)
	Exit(g *Game)
	Update(g *Game)
	Draw(screen *ebiten.Image, g *Game)
	GetUIManager() *ui.Manager
}

type SceneManager interface {
	PushScene(scene Scene)
	PopScene()
	ReplaceScene(scene Scene)
	CurrentScene() Scene
}
