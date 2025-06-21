package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/thanhfphan/ebitengj2025/assets/fonts"
	"github.com/thanhfphan/ebitengj2025/internal/ai"
	"github.com/thanhfphan/ebitengj2025/internal/am"
	"github.com/thanhfphan/ebitengj2025/internal/card"
	"github.com/thanhfphan/ebitengj2025/internal/entity"
	"github.com/thanhfphan/ebitengj2025/internal/renderer"
	"github.com/thanhfphan/ebitengj2025/internal/rules"
	"github.com/thanhfphan/ebitengj2025/internal/ui"
	"github.com/thanhfphan/ebitengj2025/internal/world"
)

var _ ebiten.Game = (*Game)(nil)

type Game struct {
	State       GameState
	CurrentTurn int
	Players     []*entity.Player
	Player      *entity.Player
	PlayerHand  *ui.UIHand
	BotHands    []*ui.UIBotHand
	DebugMode   bool

	AssetManager *am.AssetManager
	Renderer     *renderer.Renderer
	World        *world.World
	UIManager    *ui.Manager
	AIManager    *ai.Manager
	CardManager  *card.Manager
	TurnManager  *rules.TurnManager
}

func New() (*Game, error) {
	assetManager := am.NewAssetManager()
	renderer := renderer.New(assetManager)
	world := world.New()
	uiManager := ui.NewManager()
	aiManager := ai.NewManager()
	cardManager := card.NewManager()
	turnManager := rules.NewTurnManager()

	g := &Game{
		State:        GameStateMainMenu,
		CurrentTurn:  0,
		Players:      []*entity.Player{},
		BotHands:     []*ui.UIBotHand{},
		AssetManager: assetManager,
		Renderer:     renderer,
		World:        world,
		UIManager:    uiManager,
		AIManager:    aiManager,
		CardManager:  cardManager,
		TurnManager:  turnManager,
	}

	cardManager.OnDishMade = func(recipe *entity.Card) {
		fmt.Println("Recipe made:", recipe.Name)
	}
	cardManager.OnPlayCard = func(player *entity.Player, card *entity.Card) {
		fmt.Println("Card played:", card.Name, "by", player.Name)
	}

	g.AssetManager.LoadFont("default", fonts.MPlus1pRegular_ttf, 24)
	defaultFont := g.AssetManager.GetFont("default")

	// Setup table UI
	table := ui.NewUITable(640, 360, 290)
	// table.SetVisible(false)
	g.UIManager.AddElement(table)

	// Setup player hand UI
	g.PlayerHand = ui.NewUIHand(320, 550, 640, 150)
	g.PlayerHand.SetVisible(true)
	g.PlayerHand.SetOnPlayCard(func(cardIndex int) {
		if g.Player != nil {
			g.PlayCard(g.Player.ID, cardIndex)
		}
	})
	g.UIManager.AddElement(g.PlayerHand)

	// Setup buttons
	passBtn := ui.NewUIButton(860, 550, 100, 40, "Pass", defaultFont)
	passBtn.SetVisible(true)
	passBtn.OnClick = func() {
		if g.Player != nil {
			g.Pass(g.Player.ID)
		}
	}
	g.UIManager.AddElement(passBtn)

	playBtn := ui.NewUIButton(860, 600, 100, 40, "Play", defaultFont)
	playBtn.SetVisible(true)
	playBtn.OnClick = func() {
		g.PlayerHand.PlaySelected()
	}
	g.UIManager.AddElement(playBtn)

	menuTitle := ui.NewUIButton(640, 200, 300, 60, "Food Cards", defaultFont)
	menuTitle.SetVisible(true)
	menuTitle.BackgroundColor = color.RGBA{0, 0, 0, 0}
	menuTitle.TextColor = color.RGBA{255, 255, 0, 255} // Yellow text
	g.UIManager.AddElement(menuTitle)

	newGameBtn := ui.NewUIButton(640, 300, 200, 50, "New Game", defaultFont)
	newGameBtn.SetVisible(true)
	newGameBtn.OnClick = func() {
		g.setupSoloMatch(3)
		g.State = GameStatePlaying
		g.updateElementsVisibility()
	}
	g.UIManager.AddElement(newGameBtn)

	settingsBtn := ui.NewUIButton(640, 370, 200, 50, "Settings", defaultFont)
	settingsBtn.SetVisible(true)
	settingsBtn.OnClick = func() {
		g.State = GameStateSettings
		g.updateElementsVisibility()
	}
	g.UIManager.AddElement(settingsBtn)

	quitBtn := ui.NewUIButton(640, 440, 200, 50, "Quit", defaultFont)
	quitBtn.SetVisible(true)
	quitBtn.OnClick = func() {
		// Signal to close the game
		g.State = GameStateQuit
	}
	g.UIManager.AddElement(quitBtn)

	g.updateElementsVisibility()

	return g, nil
}

func (g *Game) setupSoloMatch(botCount int) {
	g.CardManager.LoadDeck("default")
	g.TurnManager.Reset()
	g.Players = []*entity.Player{}
	g.Player = nil
	for _, hand := range g.BotHands {
		g.UIManager.RemoveElement(hand)
	}
	g.BotHands = []*ui.UIBotHand{}

	g.Player = entity.NewPlayer("P0", entity.TypePlayer, 0, 0)
	g.Players = append(g.Players, g.Player)
	g.TurnManager.AddPlayer(g.Player.ID, false)

	defaultFont := g.AssetManager.GetFont("default")

	for i := 1; i <= botCount; i++ {
		bot := entity.NewPlayer(fmt.Sprintf("B%d", i), entity.TypeBot, 0, 0)
		g.Players = append(g.Players, bot)
		g.TurnManager.AddPlayer(bot.ID, true)
		g.AIManager.RegisterBot(bot.ID, ai.NewEasyBot())

		var botHand *ui.UIBotHand
		switch i {
		case 1: // Left side
			botHand = ui.NewUIBotHand(150, 360, 80, 120, defaultFont)
		case 2: // Top
			botHand = ui.NewUIBotHand(320, 170, 80, 120, defaultFont)
		case 3: // Right side
			botHand = ui.NewUIBotHand(930, 360, 80, 120, defaultFont)
		default:
			botHand = ui.NewUIBotHand(320, 170, 80, 120, defaultFont)
		}

		botHand.SetVisible(true)
		g.BotHands = append(g.BotHands, botHand)
		g.UIManager.AddElement(botHand)
	}

	g.CardManager.DealHands(g.Players)
}

// Update implements ebiten.Game.
func (g *Game) Update() error {
	g.HandleInput()

	switch g.State {
	case GameStateMainMenu:
		g.UpdateMainMenu()
	case GameStatePlaying:
		g.UIManager.Update()
		g.UpdatePlaying()
	case GameStateSettings:
		g.UIManager.Update()
		g.UpdateSettings()
	case GameStateQuit:
		return ebiten.Termination
	}

	g.Renderer.Update()

	return nil
}

// Draw implements ebiten.Game.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.Renderer.BackgroundColor())

	switch g.State {
	case GameStateMainMenu:
		g.DrawMainMenu(screen)
	case GameStatePlaying:
		g.DrawPlaying(screen)
	case GameStateSettings:
		g.DrawSettings(screen)
	}

}

// Layout implements ebiten.Game.
func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return 1280, 720
}

func (g *Game) UpdateMainMenu() {
	g.UIManager.Update()
}

func (g *Game) UpdateSettings() {
	// TODO: implement
}

func (g *Game) UpdatePlayerHand() {
	if g.Player == nil || g.PlayerHand == nil {
		return
	}

	// Update player's hand
	cardImages := make(map[string]*ebiten.Image)
	for _, card := range g.Player.Hand {
		cardImg := g.AssetManager.GetCardImage(card.ID)
		if cardImg == nil {
			cardImg = ebiten.NewImage(80, 120)
		}
		cardImages[card.ID] = cardImg
	}

	g.PlayerHand.UpdateCards(g.Player.Hand, cardImages)

	cardBackImage := g.AssetManager.GetCardBackImage()
	if cardBackImage == nil {
		cardBackImage = ebiten.NewImage(80, 120)
		cardBackImage.Fill(color.RGBA{100, 100, 100, 255}) // Gray color for card back
	}

	for i, botHand := range g.BotHands {
		if i+1 >= len(g.Players) {
			continue
		}

		// Index 0 is the player, so we start from index 1
		bot := g.Players[i+1]
		botHand.UpdateCards(bot.Hand, cardBackImage)
	}
}

func (g *Game) UpdateTurn() {
	current := g.TurnManager.Current()
	if current == nil {
		return
	}
	if current.IsBot {
		g.AIManager.OnTurn(current.ID, g)
	}
}

func (g *Game) UpdatePlaying() {
	// TODO: improve this one
	if winnerOrder := g.TurnManager.FinishedOrder(); len(winnerOrder) > 0 {
		fmt.Println("Game finished. Winner order:", winnerOrder)
		g.State = GameStateMainMenu
		g.updateElementsVisibility()
		return
	}

	g.UpdateTurn()
	// Update player's hand UI
	g.UpdatePlayerHand()
}

func (g *Game) DrawMainMenu(screen *ebiten.Image) {
	g.UIManager.Draw(screen)
}

func (g *Game) DrawPlaying(screen *ebiten.Image) {
	g.Renderer.DrawWorld(screen, g.World)
	//  g.RenderTableAndHands(screen, g.CardManager)
	g.UIManager.Draw(screen)

	if g.DebugMode {
		mx, my := ebiten.CursorPosition()
		g.Renderer.DrawDebug(screen, fmt.Sprintf("Cursor: %d, %d", mx, my))
	}
}

func (g *Game) DrawSettings(screen *ebiten.Image) {
	g.Renderer.DrawDebug(screen, "Settings Menu - Press ESC to return to main menu")
	g.UIManager.Draw(screen)
}

func (g *Game) HandleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		g.DebugMode = !g.DebugMode
	}

	// ESC key to return to main menu from settings
	if g.State == GameStateSettings && inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.State = GameStateMainMenu
		g.updateElementsVisibility()
	}

	if g.State == GameStatePlaying {
		current := g.TurnManager.Current()
		if current != nil && !current.IsBot && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			if err := g.PlayCard(current.ID, 0); err != nil {
				fmt.Println("Error playing card:", err)
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		g.UIManager.HandleMouseDown(mx, my)
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		g.UIManager.HandleMouseUp(mx, my)
	}
}

func (g *Game) NextTurn() {
	g.CurrentTurn = (g.CurrentTurn + 1) % len(g.Players)

	player := g.Players[g.CurrentTurn]
	if player.IsBot() {
		g.AIManager.OnTurn(player.ID, g) // Automatically play the bot's turn
	}
}

func (g *Game) GetPlayerState(id string) *ai.PlayerState {
	playerTurn := g.TurnManager.GetPlayerByID(id)
	if playerTurn == nil {
		return nil
	}

	player := g.GetPlayer(id)

	return &ai.PlayerState{
		ID:       playerTurn.ID,
		IsBot:    playerTurn.IsBot,
		Hand:     player.Hand,
		Passed:   playerTurn.Passed,
		Finished: playerTurn.Finished,
	}
}

func (g *Game) Pass(playerID string) {
	current := g.TurnManager.Current()
	if current == nil || current.ID != playerID || current.Finished {
		fmt.Println("Cannot pass because it's not your turn. Player:", playerID, "Current:", current)
		return
	}

	if err := g.TurnManager.Pass(playerID); err != nil {
		fmt.Println("Error passing turn:", err)
	}

	g.TurnManager.Next()
}

func (g *Game) PlayCard(playerID string, cardIdx int) error {
	current := g.TurnManager.Current()
	if current == nil || current.ID != playerID || current.Finished {
		fmt.Println("Cannot play card because it's not your turn. Player:", playerID, "Current:", current)
		return fmt.Errorf("Cannot play card because it's not your turn. Player: %s, Current: %v", playerID, current)
	}

	player := g.GetPlayer(playerID)
	hasDish := g.CardManager.PlayCard(player, cardIdx)

	g.TurnManager.MarkAllUnpassed()

	if player != nil && len(player.Hand) == 0 {
		g.TurnManager.MarkFinished(playerID)
	}

	if !hasDish {
		g.TurnManager.Next()
	}

	return nil
}

func (g *Game) GetPlayer(id string) *entity.Player {
	for _, p := range g.Players {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func (g *Game) updateElementsVisibility() {
	// Hide/show elements based on current game state
	for _, element := range g.UIManager.Elements {
		switch btn := element.(type) {
		case *ui.UIButton:
			if btn.Text == "New Game" || btn.Text == "Settings" || btn.Text == "Quit" || btn.Text == "Food Cards" {
				btn.SetVisible(g.State == GameStateMainMenu)
			}
			if btn.Text == "Pass" || btn.Text == "Play" {
				btn.SetVisible(g.State == GameStatePlaying)
			}
		case *ui.UIHand:
			btn.SetVisible(g.State == GameStatePlaying)
		case *ui.UIBotHand:
			btn.SetVisible(g.State == GameStatePlaying)
		case *ui.UITable:
			btn.SetVisible(g.State == GameStatePlaying)
		}
	}
}
