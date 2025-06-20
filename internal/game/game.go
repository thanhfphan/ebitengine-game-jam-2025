package game

import (
	"fmt"

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
		// State:        GameStateMainMenu,
		State:        GameStatePlaying,
		CurrentTurn:  0,
		Players:      []*entity.Player{},
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

	table := ui.NewUICircle(640, 360, 290)
	table.SetVisible(true)
	g.UIManager.AddElement(table)

	passBtn := ui.NewUIButton(600, 400, 70, 90, "Pass", defaultFont)
	passBtn.SetVisible(true)
	g.UIManager.AddElement(passBtn)

	playBtn := ui.NewUIButton(500, 400, 70, 90, "Play", defaultFont)
	playBtn.SetVisible(true)
	g.UIManager.AddElement(playBtn)

	newGameBtn := ui.NewUIButton(500, 300, 150, 50, "New Game", defaultFont)
	playBtn.SetVisible(true)
	g.UIManager.AddElement(newGameBtn)

	// Quick test

	g.setupSoloMatch(3)

	return g, nil
}

func (g *Game) setupSoloMatch(botCount int) {
	g.CardManager.LoadDeck("default")
	g.TurnManager.Reset()
	g.Players = []*entity.Player{}
	g.Player = nil

	g.Player = entity.NewPlayer("P0", entity.TypePlayer, 0, 0)
	g.Players = append(g.Players, g.Player)
	g.TurnManager.AddPlayer(g.Player.ID, false)

	for i := 1; i <= botCount; i++ {
		bot := entity.NewPlayer(fmt.Sprintf("B%d", i), entity.TypeBot, 0, 0)
		g.Players = append(g.Players, bot)
		g.TurnManager.AddPlayer(bot.ID, true)
		g.AIManager.RegisterBot(bot.ID, ai.NewEasyBot())
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
		g.UpdatePlaying()
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
	}

}

// Layout implements ebiten.Game.
func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return 1280, 720
}

func (g *Game) UpdateMainMenu() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.State = GameStatePlaying
	}
}

func (g *Game) UpdatePlaying() {
	current := g.TurnManager.Current()
	if current == nil {
		return
	}
	if current.IsBot {
		g.AIManager.OnTurn(current.ID, g)
	}

	// TODO: improve this one
	if winnerOrder := g.TurnManager.FinishedOrder(); len(winnerOrder) > 0 {
		fmt.Println("Game finished. Winner order:", winnerOrder)
		g.State = GameStateMainMenu
	}
}

func (g *Game) DrawMainMenu(screen *ebiten.Image) {
	g.Renderer.DrawDebug(screen, "Main Menu - press <Enter> to start solo game")
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

func (g *Game) HandleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		g.DebugMode = !g.DebugMode
	}

	current := g.TurnManager.Current()
	if current == nil {
		return
	}
	if !current.IsBot && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if err := g.PlayCard(current.ID, 0); err != nil {
			fmt.Println("Error playing card:", err)
		}
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
	err := g.CardManager.PlayCard(player, cardIdx)
	if err != nil {
		return fmt.Errorf("Error playing card: %w", err)
	}

	g.TurnManager.MarkAllUnpassed()

	if player != nil && len(player.Hand) == 0 {
		g.TurnManager.MarkFinished(playerID)
	}

	for _, p := range g.Players {
		fmt.Println(p, "Hand:", p.Hand)
	}

	g.TurnManager.Next()

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
