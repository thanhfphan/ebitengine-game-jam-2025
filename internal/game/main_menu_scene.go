package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thanhfphan/ebitengj2025/internal/ui"
)

type MainMenuScene struct {
	elements []ui.Element
}

func NewMainMenuScene() *MainMenuScene {
	return &MainMenuScene{
		elements: []ui.Element{},
	}
}

func (s *MainMenuScene) Enter(g *Game) {
	g.UIManager = ui.NewManager()
	defaultFont := g.AssetManager.GetFont("default")

	// Create menu title
	menuTitle := ui.NewUIButton(640, 200, 300, 60, "Food Cards", defaultFont)
	menuTitle.BackgroundColor = color.RGBA{0, 0, 0, 0} // Transparent background
	menuTitle.TextColor = color.RGBA{255, 255, 0, 255} // Yellow text
	menuTitle.SetTags(ui.TagMenu)
	g.UIManager.AddElement(menuTitle)
	s.elements = append(s.elements, menuTitle)

	// Create New Game button
	newGameBtn := ui.NewUIButton(640, 300, 200, 50, "New Game", defaultFont)
	newGameBtn.OnClick = func() {
		g.SetScene(NewPlayingScene())
	}
	newGameBtn.SetTags(ui.TagMenu)
	g.UIManager.AddElement(newGameBtn)
	s.elements = append(s.elements, newGameBtn)

	// Create Settings button
	settingsBtn := ui.NewUIButton(640, 370, 200, 50, "Settings", defaultFont)
	settingsBtn.OnClick = func() {
		g.SetScene(NewSettingsScene())
	}
	settingsBtn.SetTags(ui.TagMenu)
	g.UIManager.AddElement(settingsBtn)
	s.elements = append(s.elements, settingsBtn)

	// Create Quit button
	quitBtn := ui.NewUIButton(640, 440, 200, 50, "Quit", defaultFont)
	quitBtn.OnClick = func() {
		g.State = GameStateQuit
	}
	quitBtn.SetTags(ui.TagMenu)
	g.UIManager.AddElement(quitBtn)
	s.elements = append(s.elements, quitBtn)

	// Set UI visibility mask for menu
	g.UIManager.SetMask(ui.TagMenu)
}

func (s *MainMenuScene) Exit(g *Game) {
	// Remove all menu elements
	for _, element := range s.elements {
		g.UIManager.RemoveElement(element)
	}
	s.elements = nil
}

func (s *MainMenuScene) Update(g *Game) {
	// Handle input is already done in Game.Update()
}

func (s *MainMenuScene) Draw(screen *ebiten.Image, g *Game) {
	// UI Manager will draw all visible elements
}
