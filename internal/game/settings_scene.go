package game

import (
	"image/color"

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
	defaultFont := g.AssetManager.GetFont("nunito", 24)
	titleFont := g.AssetManager.GetFont("nunito", 48)
	smallFont := g.AssetManager.GetFont("nunito", 18)

	// Màu sắc giống main menu
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
	title.SetTags(ui.TagSettings)
	g.UIManager.AddElement(title)
	s.elements = append(s.elements, title)

	y := startY + 100

	// Music volume label + slider
	musicLabel := ui.NewUILabel(cx, y, "Music Volume", defaultFont)
	musicLabel.AlignCenter()
	musicLabel.SetTags(ui.TagSettings)
	g.UIManager.AddElement(musicLabel)
	s.elements = append(s.elements, musicLabel)

	y += spacing
	musicSlider := ui.NewUISlider(cx, y, 400, 30, g.AssetManager.GetMusicVolume())
	musicSlider.SetTags(ui.TagSettings)
	musicSlider.OnChange = func(value float64) {
		g.AssetManager.SetMusicVolume(value)
	}
	g.UIManager.AddElement(musicSlider)
	s.elements = append(s.elements, musicSlider)

	// Sound Effects volume
	y += spacing
	soundLabel := ui.NewUILabel(cx, y, "Sound Effects Volume", defaultFont)
	soundLabel.AlignCenter()
	soundLabel.SetTags(ui.TagSettings)
	g.UIManager.AddElement(soundLabel)
	s.elements = append(s.elements, soundLabel)

	y += spacing
	soundSlider := ui.NewUISlider(cx, y, 400, 30, g.AssetManager.GetMasterVolume())
	soundSlider.SetTags(ui.TagSettings)
	soundSlider.OnChange = func(value float64) {
		g.AssetManager.SetMasterVolume(value)
	}
	g.UIManager.AddElement(soundSlider)
	s.elements = append(s.elements, soundSlider)

	// Test sound button
	y += spacing + 10
	testBtn := ui.NewUIButton(cx-100, y, 200, 40, "Test Sound", smallFont)
	testBtn.BackgroundColor = colButtonBg
	testBtn.HoverColor = colButtonHover
	testBtn.PressedColor = colButtonPressed
	testBtn.TextColor = colButtonText
	testBtn.OnClick = func() {
		g.AssetManager.PlaySound("click")
	}
	testBtn.SetTags(ui.TagSettings)
	g.UIManager.AddElement(testBtn)
	s.elements = append(s.elements, testBtn)

	// Back button
	y += spacing + 20
	backBtn := ui.NewUIButton(cx-100, y, 200, 50, "Back", defaultFont)
	backBtn.BackgroundColor = colButtonBg
	backBtn.HoverColor = colButtonHover
	backBtn.PressedColor = colButtonPressed
	backBtn.TextColor = colButtonText
	backBtn.OnClick = func() {
		g.SetScene(NewMainMenuScene())
	}
	backBtn.SetTags(ui.TagSettings)
	g.UIManager.AddElement(backBtn)
	s.elements = append(s.elements, backBtn)

	// UI mask
	g.UIManager.SetMask(ui.TagSettings)
}

func (s *SettingsScene) Exit(g *Game) {
	for _, element := range s.elements {
		g.UIManager.RemoveElement(element)
	}
	s.elements = nil
}

func (s *SettingsScene) Update(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.SetScene(NewMainMenuScene())
	}
}

func (s *SettingsScene) Draw(screen *ebiten.Image, g *Game) {
	// UI Manager draws
}
