package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/thanhfphan/ebitengj2025/internal/entity"
	"github.com/thanhfphan/ebitengj2025/internal/ui"
	"github.com/thanhfphan/ebitengj2025/internal/view"
	"golang.org/x/image/font"
)

var _ Scene = (*PlayingScene)(nil)

type PlayingScene struct {
	elements   []ui.Element
	playerHand *ui.UIHand
	botHands   []*ui.UIBotHand
	tableCards *ui.UITableCards
	isPaused   bool
	pauseMenu  *PauseMenu
	uiManager  *ui.Manager
	bgImage    *ebiten.Image
}

type PauseMenu struct {
	elements []ui.Element
	overlay  *ebiten.Image
}

func NewPlayingScene() *PlayingScene {
	return &PlayingScene{
		elements:  []ui.Element{},
		botHands:  []*ui.UIBotHand{},
		isPaused:  false,
		pauseMenu: nil,
		uiManager: ui.NewManager(),
	}

}

func (s *PlayingScene) Enter(g *Game) {
	g.CurrentUIManager = s.uiManager
	if s.isPaused { // FIXME: might want to add resume method instead??
		return
	}

	s.bgImage = g.AssetManager.GetImage(ImagePlayBG)
	tableBgImage := g.AssetManager.GetImage(ImageTableBG)
	settingsIconImg := g.AssetManager.GetImage(ImageSettingIcon)

	defaultFont := g.AssetManager.GetFont("nunito", 24)
	centerX, centerY := ScreenW/2, ScreenH/2

	if settingsIconImg != nil {
		settingsBtn := ui.NewUIImageButton(ScreenW-60, 20, 40, 40, settingsIconImg)
		settingsBtn.BackgroundColor = color.RGBA{0x00, 0x00, 0x00, 0x00} // Black
		settingsBtn.HoverColor = color.RGBA{0xFF, 0xE0, 0x7A, 0xFF}
		settingsBtn.PressedColor = color.RGBA{0xD9, 0xC3, 0x90, 0xFF}
		settingsBtn.OnClick = func() {
			s.togglePause(g)
		}
		s.uiManager.AddElement(settingsBtn)
		s.elements = append(s.elements, settingsBtn)
	}

	// Setup table cards UI
	s.tableCards = ui.NewUITableCards(centerX, centerY, TableRadius, tableBgImage)
	s.uiManager.AddElement(s.tableCards)
	s.elements = append(s.elements, s.tableCards)

	// Setup player hand UI
	handWidth := 500
	handHeight := 160
	s.playerHand = ui.NewUIHand(centerX-handWidth/2, 600, handWidth, handHeight)
	s.playerHand.SetOnPlayCard(func(cardID string) {
		if g.Player != nil {
			g.PlayCard(g.Player.ID, cardID)
		}
	})

	// Add card selection handler to highlight matching recipes
	s.playerHand.SetOnCardSelected(func(cardID string) {
		s.highlightMatchingRecipes(g, cardID)
		if err := g.AssetManager.PlaySound(SoundSelect); err != nil {
			fmt.Println("Error playing sound:", err)
		}
	})

	s.uiManager.AddElement(s.playerHand)
	s.elements = append(s.elements, s.playerHand)

	// Setup buttons
	btnX := centerX + 300
	passBtn := ui.NewUIButton(btnX, 600, 100, 40, "Pass", defaultFont)
	passBtn.OnClick = func() {
		if g.Player != nil {
			g.Pass(g.Player.ID)
		}
	}
	s.uiManager.AddElement(passBtn)
	s.elements = append(s.elements, passBtn)

	playBtn := ui.NewUIButton(btnX, 650, 100, 40, "Play", defaultFont)
	playBtn.OnClick = func() {
		s.playerHand.PlaySelected()
	}
	s.uiManager.AddElement(playBtn)
	s.elements = append(s.elements, playBtn)

	s.initPauseMenu(g)

	s.setupGame(g)
}

func (s *PlayingScene) setupGame(g *Game) {
	numBots := 3
	s.botHands = g.setupGameData(numBots)

	// Position bot hands
	reserved := math.Pi / 3
	if numBots >= 4 {
		reserved = 2 * math.Pi / 3
	}
	arc := 2*math.Pi - reserved // For bots
	angleGap := arc / float64(numBots+1)
	startAngle := 3*math.Pi/2 + reserved/2
	centerX, centerY := ScreenW/2, ScreenH/2
	for i, hand := range s.botHands {
		angle := math.Mod(startAngle+angleGap*float64(i+1), 2*math.Pi)
		x := int(float64(centerX) + float64(TableRadius)*math.Cos(angle))
		y := int(float64(centerY) - float64(TableRadius)*math.Sin(angle))
		hand.SetPosition(x-hand.Width/2, y-hand.Height/2)
	}
}

func (s *PlayingScene) Exit(g *Game) {
}

func (s *PlayingScene) Update(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.togglePause(g)
		return
	}

	// Don't update game state if paused
	if s.isPaused {
		return
	}

	// Check for game over
	if winnerOrder := g.TurnManager.FinishedOrder(); len(winnerOrder) == len(g.Players) {
		fmt.Println("Game finished. Winner order:", winnerOrder)
		g.PopScene()
		g.PushScene(NewMainMenuScene())
		return
	}

	g.UpdateTurn()

	s.UpdateHands(g)
}

func (s *PlayingScene) UpdateHands(g *Game) {
	if s.tableCards == nil {
		return
	}

	fonts := map[string]font.Face{
		"title":    g.AssetManager.GetFont("nunito", 14),
		"subtitle": g.AssetManager.GetFont("nunito", 12),
		"body":     g.AssetManager.GetFont("nunito", 10),
	}
	ingredientNames := g.CardManager.GetMapIngredientNames()
	viewTableStack := ToViewTableStack(g.CardManager.TableStack)

	// Update table cards
	s.tableCards.UpdateFromTableStack(viewTableStack, fonts, ingredientNames)

	// Update player's hand
	viewPlayerCards := make([]view.Card, 0, len(g.Player.Hand))
	for _, id := range g.Player.OrderHand {
		card := g.Player.GetCard(id)
		viewPlayerCards = append(viewPlayerCards, ToViewCard(card))
	}
	s.playerHand.UpdateCards(viewPlayerCards, viewTableStack, fonts, ingredientNames)

	// Update bot hands
	cardBackImage := g.AssetManager.GetImage(ImageCardBack)
	for i, botHand := range s.botHands {
		if i+1 >= len(g.Players) {
			continue
		}

		// Index 0 is the player, so we start from index 1
		bot := g.Players[i+1]
		botViewCards := make([]view.Card, 0, len(bot.Hand))
		for _, card := range bot.Hand {
			botViewCards = append(botViewCards, ToViewCard(card))
		}
		botHand.UpdateCards(botViewCards, cardBackImage)
	}
}

func (s *PlayingScene) Draw(screen *ebiten.Image, g *Game) {
	if s.bgImage != nil {
		op := &ebiten.DrawImageOptions{}

		sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
		bw, bh := s.bgImage.Bounds().Dx(), s.bgImage.Bounds().Dy()

		sx := float64(sw) / float64(bw)
		sy := float64(sh) / float64(bh)

		op.GeoM.Scale(sx, sy)
		screen.DrawImage(s.bgImage, op)
	}

	// Draw pause overlay if paused
	if s.isPaused && s.pauseMenu != nil && s.pauseMenu.overlay != nil {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(s.pauseMenu.overlay, op)
	}

	if g.DebugMode {
		mx, my := ebiten.CursorPosition()
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Cursor: %d, %d", mx, my))
	}
}

func (s *PlayingScene) GetUIManager() *ui.Manager {
	return s.uiManager
}

// Add a new function to highlight matching recipes
func (s *PlayingScene) highlightMatchingRecipes(g *Game, cardID string) {
	if cardID == "" {
		s.tableCards.ResetCanMakeDish()
		return
	}

	selectedCard := g.Player.GetCard(cardID)
	if selectedCard == nil || selectedCard.Type != entity.CardTypeIngredient {
		return
	}

	s.tableCards.ResetCanMakeDish()
	viewTableStack := ToViewTableStack(g.CardManager.TableStack)
	s.tableCards.UpdateCanMakeDish(selectedCard.IngredientID, viewTableStack)
}

// Initialize the pause menu
func (s *PlayingScene) initPauseMenu(g *Game) {
	s.pauseMenu = &PauseMenu{
		elements: []ui.Element{},
	}

	s.pauseMenu.overlay = ebiten.NewImage(ScreenW, ScreenH)
	s.pauseMenu.overlay.Fill(color.RGBA{0, 0, 0, 180})

	// Setup pause menu buttons
	titleFont := g.AssetManager.GetFont("nunito", 48)
	defaultFont := g.AssetManager.GetFont("nunito", 24)

	centerX := ScreenW / 2
	startY := ScreenH / 3
	btnWidth := 300
	btnHeight := 50
	btnSpacing := 70

	// Pause title
	pauseTitle := ui.NewUILabel(centerX, startY-80, "PAUSED", titleFont)
	pauseTitle.AlignCenter()
	pauseTitle.TextColor = color.RGBA{0xFF, 0xE7, 0x4D, 0xFF}
	s.uiManager.AddElement(pauseTitle)
	s.pauseMenu.elements = append(s.pauseMenu.elements, pauseTitle)

	// Button colors
	colButtonBg := color.RGBA{0xF3, 0xE2, 0xC3, 0xFF}
	colButtonHover := color.RGBA{0xFF, 0xE0, 0x7A, 0xFF}
	colButtonPressed := color.RGBA{0xD9, 0xC3, 0x90, 0xFF}
	colButtonText := color.RGBA{0x36, 0x55, 0x34, 0xFF}

	// Resume button
	resumeBtn := ui.NewUIButton(centerX-btnWidth/2, startY, btnWidth, btnHeight, "Resume", defaultFont)
	resumeBtn.BackgroundColor = colButtonBg
	resumeBtn.HoverColor = colButtonHover
	resumeBtn.PressedColor = colButtonPressed
	resumeBtn.TextColor = colButtonText
	resumeBtn.OnClick = func() {
		s.togglePause(g)
	}
	s.uiManager.AddElement(resumeBtn)
	s.pauseMenu.elements = append(s.pauseMenu.elements, resumeBtn)

	// Settings button
	settingsBtn := ui.NewUIButton(centerX-btnWidth/2, startY+btnSpacing, btnWidth, btnHeight, "Settings", defaultFont)
	settingsBtn.BackgroundColor = colButtonBg
	settingsBtn.HoverColor = colButtonHover
	settingsBtn.PressedColor = colButtonPressed
	settingsBtn.TextColor = colButtonText
	settingsBtn.OnClick = func() {
		g.PushScene(NewSettingsScene())
	}
	s.uiManager.AddElement(settingsBtn)
	s.pauseMenu.elements = append(s.pauseMenu.elements, settingsBtn)

	// Main menu button
	menuBtn := ui.NewUIButton(centerX-btnWidth/2, startY+2*btnSpacing, btnWidth, btnHeight, "Return to Main Menu", defaultFont)
	menuBtn.BackgroundColor = colButtonBg
	menuBtn.HoverColor = colButtonHover
	menuBtn.PressedColor = colButtonPressed
	menuBtn.TextColor = colButtonText
	menuBtn.OnClick = func() {
		g.PopScene()
		g.PushScene(NewMainMenuScene())
	}
	s.uiManager.AddElement(menuBtn)
	s.pauseMenu.elements = append(s.pauseMenu.elements, menuBtn)

	// Hide pause menu initially
	for _, element := range s.pauseMenu.elements {
		element.SetVisible(false)
	}
}

// Toggle pause state
func (s *PlayingScene) togglePause(g *Game) {
	s.isPaused = !s.isPaused

	// Show/hide pause menu elements
	for _, element := range s.pauseMenu.elements {
		element.SetVisible(s.isPaused)
	}
}
