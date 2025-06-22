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
	f := g.AssetManager.GetFont("default")

	var (
		phi    = 1.618 // Golden ratio
		cx     = ScreenW / 2
		blockH = float64(ScreenH) / phi
		startY = int(ScreenH/2) - int(blockH/2)
		btnW   = ScreenW / 3
		btnH   = 56
		gapY   = 32
	)

	makeBtn := func(label string, y int, onClick func()) *ui.UIButton {
		b := ui.NewUIButton(cx-btnW/2, y, btnW, btnH, label, f)
		b.BackgroundColor = color.RGBA{0xF5, 0xDE, 0xB3, 0xFF}
		b.HoverColor = color.RGBA{0xFF, 0xD7, 0x00, 0xFF}
		b.PressedColor = color.RGBA{0xAA, 0x88, 0x00, 0xFF}
		b.TextColor = color.RGBA{0x4B, 0x2E, 0x2B, 0xFF}
		b.OnClick = onClick
		b.SetTags(ui.TagMenu)
		g.UIManager.AddElement(b)
		return b
	}

	// Title
	title := ui.NewUILabel(cx, int(startY), "FOOD CARDS", f)
	title.AlignCenter()
	title.TextColor = color.RGBA{255, 213, 0, 255}
	title.HoverColor = color.RGBA{255, 255, 255, 255}
	title.HoverScale = 1.08
	title.EnableHover = true
	title.SetTags(ui.TagMenu)
	g.UIManager.AddElement(title)
	s.elements = append(s.elements, title)

	// Buttons
	y := startY + 120
	s.elements = append(s.elements,
		makeBtn("New Game", y, func() { g.SetScene(NewPlayingScene()) }),
		makeBtn("Settings", y+btnH+gapY, func() { g.SetScene(NewSettingsScene()) }),
		makeBtn("Quit", y+2*(btnH+gapY), func() { g.State = GameStateQuit }),
	)

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
