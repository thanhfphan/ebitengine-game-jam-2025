package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/thanhfphan/ebitengj2025/internal/ui"
)

var _ Scene = (*SettingsScene)(nil)

type SettingsScene struct {
	elements  []ui.Element
	bgImage   *ebiten.Image
	uiManager *ui.Manager
}

func NewSettingsScene() *SettingsScene {
	return &SettingsScene{
		elements: []ui.Element{},
	}
}

func (s *SettingsScene) Enter(g *Game) {
	s.uiManager = ui.NewManager()
	g.CurrentUIManager = s.uiManager

	s.bgImage = g.AssetManager.GetImage(ImageMainBG)

	defaultFont := g.AssetManager.GetFont("nunito", 24)
	titleFont := g.AssetManager.GetFont("nunito", 48)
	smallFont := g.AssetManager.GetFont("nunito", 18)

	var (
		colButtonBg      = color.RGBA{0xF3, 0xE2, 0xC3, 0xFF}
		colButtonHover   = color.RGBA{0xFF, 0xE0, 0x7A, 0xFF}
		colButtonPressed = color.RGBA{0xD9, 0xC3, 0x90, 0xFF}
		colButtonText    = color.RGBA{0x36, 0x55, 0x34, 0xFF}
		colTitle         = color.RGBA{0xFF, 0xE7, 0x4D, 0xFF}
		colTitleHover    = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	)

	cx := ScreenW / 2
	startY := 100
	spacing := 60

	// Title
	title := ui.NewUILabel(cx, startY, "SETTINGS", titleFont)
	title.AlignCenter()
	title.TextColor = colTitle
	title.HoverColor = colTitleHover
	title.HoverScale = 1.08
	title.EnableHover = true
	s.uiManager.AddElement(title)
	s.elements = append(s.elements, title)

	y := startY + 100

	// Music volume label + slider
	musicLabel := ui.NewUILabel(cx, y, "Music Volume", defaultFont)
	musicLabel.AlignCenter()
	s.uiManager.AddElement(musicLabel)
	s.elements = append(s.elements, musicLabel)

	y += spacing
	musicSlider := ui.NewUISlider(cx, y, 400, 30, g.AssetManager.GetMusicVolume())
	musicSlider.OnChange = func(value float64) {
		g.AssetManager.SetMusicVolume(value)
	}
	s.uiManager.AddElement(musicSlider)
	s.elements = append(s.elements, musicSlider)

	// Sound Effects volume
	y += spacing
	soundLabel := ui.NewUILabel(cx, y, "Sound Effects Volume", defaultFont)
	soundLabel.AlignCenter()
	s.uiManager.AddElement(soundLabel)
	s.elements = append(s.elements, soundLabel)

	y += spacing
	soundSlider := ui.NewUISlider(cx, y, 400, 30, g.AssetManager.GetMasterVolume())
	soundSlider.OnChange = func(value float64) {
		g.AssetManager.SetMasterVolume(value)
	}
	s.uiManager.AddElement(soundSlider)
	s.elements = append(s.elements, soundSlider)

	// Test sound button
	y += spacing + 10
	testBtn := ui.NewUIButton(cx-100, y, 200, 40, "Test Sound", smallFont)
	testBtn.BackgroundColor = colButtonBg
	testBtn.HoverColor = colButtonHover
	testBtn.PressedColor = colButtonPressed
	testBtn.TextColor = colButtonText
	testBtn.OnClick = func() {
		if err := g.AssetManager.PlaySound(SoundPlay); err != nil {
			fmt.Println("Error playing sound:", err)
		}
	}
	s.uiManager.AddElement(testBtn)
	s.elements = append(s.elements, testBtn)

	// Back button
	y += spacing + 20
	backBtn := ui.NewUIButton(cx-100, y, 200, 50, "Back", defaultFont)
	backBtn.BackgroundColor = colButtonBg
	backBtn.HoverColor = colButtonHover
	backBtn.PressedColor = colButtonPressed
	backBtn.TextColor = colButtonText
	backBtn.OnClick = func() {
		g.PopScene()
	}
	s.uiManager.AddElement(backBtn)
	s.elements = append(s.elements, backBtn)

}

func (s *SettingsScene) Exit(g *Game) {
}

func (s *SettingsScene) Update(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.PopScene() // Pop scene on Escape key
	}
}

func (s *SettingsScene) Draw(screen *ebiten.Image, g *Game) {
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

func (s *SettingsScene) GetUIManager() *ui.Manager {
	return s.uiManager
}
