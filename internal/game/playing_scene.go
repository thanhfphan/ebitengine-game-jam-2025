package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/thanhfphan/ebitengj2025/internal/ui"
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

	defaultFont := g.AssetManager.GetFont("default")

	// Setup table cards UI
	s.tableCards = ui.NewUITableCards(640, 360, 290)
	s.tableCards.SetTags(ui.TagInGame)
	g.UIManager.AddElement(s.tableCards)
	s.elements = append(s.elements, s.tableCards)

	// Setup player hand UI
	s.playerHand = ui.NewUIHand(370, 600, 640, 150)
	s.playerHand.SetTags(ui.TagInGame)
	s.playerHand.SetOnPlayCard(func(cardIndex int) {
		if g.Player != nil {
			g.PlayCard(g.Player.ID, cardIndex)
		}
	})
	g.UIManager.AddElement(s.playerHand)
	s.elements = append(s.elements, s.playerHand)

	// Setup buttons
	passBtn := ui.NewUIButton(860, 550, 100, 40, "Pass", defaultFont)
	passBtn.SetTags(ui.TagInGame)
	passBtn.OnClick = func() {
		if g.Player != nil {
			g.Pass(g.Player.ID)
		}
	}
	g.UIManager.AddElement(passBtn)
	s.elements = append(s.elements, passBtn)

	playBtn := ui.NewUIButton(860, 600, 100, 40, "Play", defaultFont)
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
	for _, hand := range s.botHands {
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

	// Tạo map các hình ảnh card
	cardImages := make(map[string]*ebiten.Image)

	// Lấy hình ảnh cho tất cả card trong TableStack
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

	// Cập nhật UI
	s.tableCards.UpdateFromTableStack(g.CardManager.TableStack, cardImages)
}

func (s *PlayingScene) Draw(screen *ebiten.Image, g *Game) {
	g.Renderer.DrawWorld(screen, g.World)

	if g.DebugMode {
		mx, my := ebiten.CursorPosition()
		g.Renderer.DrawDebug(screen, fmt.Sprintf("Cursor: %d, %d", mx, my))
	}
}
