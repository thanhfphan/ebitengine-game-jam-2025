package game

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/thanhfphan/ebitengj2025/internal/ui"
	"github.com/thanhfphan/ebitengj2025/internal/view"
)

type PlayingScene struct {
	elements   []ui.Element
	playerHand *ui.UIHand
	botHands   []*ui.UIBotHand
	tableCards *ui.UITableCards
}

func NewPlayingScene() *PlayingScene {
	return &PlayingScene{
		elements: []ui.Element{},
		botHands: []*ui.UIBotHand{},
	}
}

func (s *PlayingScene) Enter(g *Game) {
	g.UIManager = ui.NewManager()

	defaultFont := g.AssetManager.GetFont("nunito", 24)

	centerX, centerY := ScreenW/2, ScreenH/2

	// Setup table cards UI
	s.tableCards = ui.NewUITableCards(centerX, centerY, TableRadius)
	s.tableCards.SetTags(ui.TagInGame)
	g.UIManager.AddElement(s.tableCards)
	s.elements = append(s.elements, s.tableCards)

	// Setup player hand UI
	handWidth := 800
	handHeight := 160
	s.playerHand = ui.NewUIHand(centerX-handWidth/2, 600, handWidth, handHeight)
	s.playerHand.SetTags(ui.TagInGame)
	s.playerHand.SetOnPlayCard(func(cardIndex int) {
		if g.Player != nil {
			g.PlayCard(g.Player.ID, cardIndex)
		}
	})
	g.UIManager.AddElement(s.playerHand)
	s.elements = append(s.elements, s.playerHand)

	// Setup buttons
	btnX := centerX + 300
	passBtn := ui.NewUIButton(btnX, 600, 100, 40, "Pass", defaultFont)
	passBtn.SetTags(ui.TagInGame)
	passBtn.OnClick = func() {
		if g.Player != nil {
			g.Pass(g.Player.ID)
		}
	}
	g.UIManager.AddElement(passBtn)
	s.elements = append(s.elements, passBtn)

	playBtn := ui.NewUIButton(btnX, 650, 100, 40, "Play", defaultFont)
	playBtn.SetTags(ui.TagInGame)
	playBtn.OnClick = func() {
		s.playerHand.PlaySelected()
	}
	g.UIManager.AddElement(playBtn)
	s.elements = append(s.elements, playBtn)

	// Setup back to menu button
	menuBtn := ui.NewUIButton(100, 50, 120, 40, "Menu", defaultFont)
	menuBtn.SetTags(ui.TagInGame)
	menuBtn.OnClick = func() {
		g.SetScene(NewMainMenuScene())
	}
	g.UIManager.AddElement(menuBtn)
	s.elements = append(s.elements, menuBtn)

	// Setup game
	s.botHands = g.setupSoloMatch(3)
	numBots := len(s.botHands)
	reserved := math.Pi / 3
	if numBots >= 4 {
		reserved = 2 * math.Pi / 3
	}
	arc := 2*math.Pi - reserved // For bots
	angleGap := arc / float64(numBots+1)
	startAngle := 3*math.Pi/2 + reserved/2
	for i, hand := range s.botHands {
		angle := math.Mod(startAngle+angleGap*float64(i+1), 2*math.Pi)
		x := int(float64(centerX) + float64(TableRadius)*math.Cos(angle))
		y := int(float64(centerY) - float64(TableRadius)*math.Sin(angle))
		hand.SetPosition(x-hand.Width/2, y-hand.Height/2)
		hand.SetTags(ui.TagInGame)
	}

	g.UIManager.SetMask(ui.TagInGame)
}

func (s *PlayingScene) Exit(g *Game) {
	for _, element := range s.elements {
		g.UIManager.RemoveElement(element)
	}
	s.elements = nil
	s.playerHand = nil
	s.botHands = nil
	s.tableCards = nil
}

func (s *PlayingScene) Update(g *Game) {
	// Check for game over
	if winnerOrder := g.TurnManager.FinishedOrder(); len(winnerOrder) == len(g.Players) {
		fmt.Println("Game finished. Winner order:", winnerOrder)
		g.SetScene(NewMainMenuScene())
		return
	}

	g.UpdateTurn()
	g.UpdateHands(s.playerHand, s.botHands)

	s.UpdateTableCards(g)

	// Check for ESC key to return to main menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.SetScene(NewMainMenuScene())
	}
}

func (s *PlayingScene) UpdateTableCards(g *Game) {
	if s.tableCards == nil {
		return
	}

	cardImages := make(map[string]*ebiten.Image)

	for _, recipe := range g.CardManager.TableStack.Receipes {
		cardImg := g.AssetManager.GetCardImage(recipe.ID)
		if cardImg != nil {
			cardImages[recipe.ID] = cardImg
		}
	}

	for _, ingredient := range g.CardManager.TableStack.Ingredients {
		cardImg := g.AssetManager.GetCardImage(ingredient.ID)
		if cardImg != nil {
			cardImages[ingredient.ID] = cardImg
		}
	}

	viewTableStack := view.FromEntityTableStack(g.CardManager.TableStack)

	s.tableCards.UpdateFromTableStack(viewTableStack, cardImages)
}

func (s *PlayingScene) Draw(screen *ebiten.Image, g *Game) {
	g.Renderer.DrawWorld(screen, g.World)

	if g.DebugMode {
		mx, my := ebiten.CursorPosition()
		g.Renderer.DrawDebug(screen, fmt.Sprintf("Cursor: %d, %d", mx, my))
	}
}
