package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/thanhfphan/ebitengj2025/internal/ui"
)

type SettingsScene struct {
	elements []ui.Element
}

func NewSettingsScene() *SettingsScene {
	return &SettingsScene{
		elements: []ui.Element{},
	}
}

func (s *SettingsScene) Enter(g *Game) {
	g.UIManager = ui.NewManager()
	defaultFont := g.AssetManager.GetFont("default")

	// Create settings title
	title := ui.NewUIButton(640, 200, 300, 60, "Settings", defaultFont)
	title.SetTags(ui.TagSettings)
	g.UIManager.AddElement(title)
	s.elements = append(s.elements, title)

	// Create back button
	backBtn := ui.NewUIButton(640, 500, 200, 50, "Back", defaultFont)
	backBtn.OnClick = func() {
		g.SetScene(NewMainMenuScene())
	}
	backBtn.SetTags(ui.TagSettings)
	g.UIManager.AddElement(backBtn)
	s.elements = append(s.elements, backBtn)

	// Set UI visibility mask for settings
	g.UIManager.SetMask(ui.TagSettings)
}

func (s *SettingsScene) Exit(g *Game) {
	// Remove all settings elements
	for _, element := range s.elements {
		g.UIManager.RemoveElement(element)
	}
	s.elements = nil
}

func (s *SettingsScene) Update(g *Game) {
	// Check for ESC key to return to main menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.SetScene(NewMainMenuScene())
	}
}

func (s *SettingsScene) Draw(screen *ebiten.Image, g *Game) {
	// UI Manager will draw all visible elements
}
