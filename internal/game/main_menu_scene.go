package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thanhfphan/ebitengj2025/internal/ui"
)

var _ Scene = (*MainMenuScene)(nil)

type MainMenuScene struct {
	elements  []ui.Element
	bgImage   *ebiten.Image
	uiManager *ui.Manager
}

func NewMainMenuScene() *MainMenuScene {
	return &MainMenuScene{
		elements: []ui.Element{},
	}
}

func (s *MainMenuScene) Enter(g *Game) {
	s.uiManager = ui.NewManager()
	g.CurrentUIManager = s.uiManager

	s.bgImage = g.AssetManager.GetImage(ImageMainBG)

	defaultFont := g.AssetManager.GetFont("nunito", 24)
	titleFont := g.AssetManager.GetFont("nunito", 48)

	var (
		colButtonBg      = color.RGBA{0xF3, 0xE2, 0xC3, 0xFF} // #F3E2C3 – bẽ sáng (tre, rơm khô)
		colButtonHover   = color.RGBA{0xFF, 0xE0, 0x7A, 0xFF} // #FFE07A – vàng nắng lúa
		colButtonPressed = color.RGBA{0xD9, 0xC3, 0x90, 0xFF} // #D9C390 – bẽ đậm
		colButtonText    = color.RGBA{0x36, 0x55, 0x34, 0xFF} // #365534 – xanh lá sẫm

		colTitle      = color.RGBA{0xFF, 0xE7, 0x4D, 0xFF} // #FFE74D – vàng tươi nổi bật
		colTitleHover = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF} // White when hover
	)

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
		b := ui.NewUIButton(cx-btnW/2, y, btnW, btnH, label, defaultFont)
		b.BackgroundColor = colButtonBg
		b.HoverColor = colButtonHover
		b.PressedColor = colButtonPressed
		b.TextColor = colButtonText
		b.OnClick = onClick
		s.uiManager.AddElement(b)
		return b
	}

	title := ui.NewUILabel(cx, startY, "FOOD CARDS", titleFont)
	title.AlignCenter()
	title.TextColor = colTitle
	title.HoverColor = colTitleHover
	title.HoverScale = 1.08
	title.EnableHover = true
	s.uiManager.AddElement(title)
	s.elements = append(s.elements, title)

	y := startY + 120
	s.elements = append(s.elements,
		makeBtn("New Game", y, func() {
			g.PopScene()
			g.PushScene(NewPlayingScene())
		}),
		makeBtn("Settings", y+btnH+gapY, func() {
			g.PushScene(NewSettingsScene())
		}),
		makeBtn("Quit", y+2*(btnH+gapY), func() { g.State = GameStateQuit }),
	)

}

func (s *MainMenuScene) Exit(g *Game) {
}

func (s *MainMenuScene) Update(g *Game) {
	// Handle input is already done in Game.Update()
}

func (s *MainMenuScene) Draw(screen *ebiten.Image, g *Game) {
	if s.bgImage != nil {
		op := &ebiten.DrawImageOptions{}

		sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
		bw, bh := s.bgImage.Bounds().Dx(), s.bgImage.Bounds().Dy()

		sx := float64(sw) / float64(bw)
		sy := float64(sh) / float64(bh)

		op.GeoM.Scale(sx, sy)
		screen.DrawImage(s.bgImage, op)
	}

}
func (s *MainMenuScene) GetUIManager() *ui.Manager {
	return s.uiManager
}
