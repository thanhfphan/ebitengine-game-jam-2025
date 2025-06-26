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
	"github.com/thanhfphan/ebitengj2025/internal/view"
	"github.com/thanhfphan/ebitengj2025/internal/world"
	"golang.org/x/image/font"
)

var (
	_ ebiten.Game = (*Game)(nil)
	_ ai.GameLike = (*Game)(nil)
)

var (
	ScreenW, ScreenH      = 1280, 720
	TableRadius           = 300
	CardWidth, CardHeight = 80, 120
)

const (
	MusicBackground = "music_background"
	SoundSelect     = "sound_select"
	SoundRecipeMade = "sound_recipe_made"
	SoundPlay       = "sound_play"
)

type Game struct {
	State     GameState
	Players   []*entity.Player
	Player    *entity.Player
	DebugMode bool

	AssetManager *am.AssetManager
	Renderer     *renderer.Renderer
	World        *world.World
	UIManager    *ui.Manager
	AIManager    *ai.Manager
	CardManager  *card.Manager
	TurnManager  *rules.TurnManager

	currentScene Scene
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
		err := g.AssetManager.PlaySound(SoundRecipeMade)
		if err != nil {
			fmt.Println("Error playing sound:", err)
		}
	}
	cardManager.OnPlayCard = func(player *entity.Player, card *entity.Card) {
		fmt.Println("Card played:", card.Name, "by", player.Name, "(", player.ID, ")")
	}

	// Fonts
	g.AssetManager.LoadFont("nunito", fonts.NunitoRegular_ttf, 24)
	g.AssetManager.LoadFont("nunito", fonts.NunitoRegular_ttf, 32)
	g.AssetManager.LoadFont("nunito", fonts.NunitoRegular_ttf, 48)
	g.AssetManager.LoadFont("nunito", fonts.NunitoRegular_ttf, 18)

	// Music and sounds
	if err := g.AssetManager.LoadMusic(MusicBackground, "assets/sounds/vietnam-bamboo-flute.ogg"); err != nil {
		return nil, err
	}
	assetManager.SetMusicVolume(0.2)
	if err := g.AssetManager.LoadSound(SoundSelect, "assets/sounds/card-swipe.wav"); err != nil {
		return nil, err
	}
	if err := g.AssetManager.LoadSound(SoundPlay, "assets/sounds/poppop.wav"); err != nil {
		return nil, err
	}
	if err := g.AssetManager.LoadSound(SoundRecipeMade, "assets/sounds/ding-effect.wav"); err != nil {
		return nil, err
	}

	g.AssetManager.PlayMusic(MusicBackground)

	// Images
	if err := g.AssetManager.LoadImage("main_bg", "assets/images/backgrounds/main_bg.png"); err != nil {
		fmt.Println("Error loading background image:", err)
	}

	g.SetScene(NewMainMenuScene())

	return g, nil
}

func (g *Game) setupSoloMatch(botCount int) []*ui.UIBotHand {
	g.CardManager.LoadDeck("default")
	g.TurnManager.Reset()
	g.Players = []*entity.Player{}
	g.Player = nil
	botHands := []*ui.UIBotHand{}

	g.Player = entity.NewPlayer("P0", entity.TypePlayer, 0, 0)
	g.Players = append(g.Players, g.Player)
	g.TurnManager.AddPlayer(g.Player.ID, false)

	defaultFont := g.AssetManager.GetFont("nunito", 24)

	for i := 1; i <= botCount; i++ {
		bot := entity.NewPlayer(fmt.Sprintf("B%d", i), entity.TypeBot, 0, 0)
		g.Players = append(g.Players, bot)
		g.TurnManager.AddPlayer(bot.ID, true)
		g.AIManager.RegisterBot(bot.ID, ai.NewEasyBot())

		botHand := ui.NewUIBotHand(0, 0, CardWidth, CardHeight, defaultFont) // Position will be set later
		botHands = append(botHands, botHand)
		g.UIManager.AddElement(botHand)
	}

	g.CardManager.DealHands(g.Players)
	return botHands
}

func (g *Game) SetScene(scene Scene) {
	if g.currentScene != nil {
		g.currentScene.Exit(g)
	}
	g.currentScene = scene
	if g.currentScene != nil {
		g.currentScene.Enter(g)
	}
}

// Update implements ebiten.Game.
func (g *Game) Update() error {
	g.HandleInput()

	if g.currentScene != nil {
		g.currentScene.Update(g)
	}

	g.UIManager.Update()
	g.Renderer.Update()
	g.AssetManager.Update()

	if g.State == GameStateQuit {
		return ebiten.Termination
	}

	return nil
}

// Draw implements ebiten.Game.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.Renderer.BackgroundColor())

	if g.currentScene != nil {
		g.currentScene.Draw(screen, g)
	}

	g.UIManager.Draw(screen)
}

// Layout implements ebiten.Game.
func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return 1280, 720
}

func (g *Game) HandleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		g.DebugMode = !g.DebugMode
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

func (g *Game) UpdateHands(playerHand *ui.UIHand, botHands []*ui.UIBotHand) {
	if g.Player == nil {
		return
	}

	// Get fonts for card text
	fonts := map[string]font.Face{
		"title":    g.AssetManager.GetFont("nunito", 16),
		"subtitle": g.AssetManager.GetFont("nunito", 12),
		"body":     g.AssetManager.GetFont("nunito", 10),
	}

	// Convert entity.TableStack to view.TableStack for highlighting
	viewTableStack := ToViewTableStack(g.CardManager.TableStack)

	// Update player's hand
	cardImages := make(map[string]*ebiten.Image)
	viewCards := make([]view.Card, 0, len(g.Player.Hand))

	for _, id := range g.Player.OrderHand {
		card := g.Player.GetCard(id)
		viewCards = append(viewCards, ToViewCard(card))
		cardImg := g.AssetManager.GetCardImage(card.ID)
		if cardImg == nil {
			cardImg = ebiten.NewImage(80, 120)
		}
		cardImages[card.ID] = cardImg
	}
	playerHand.UpdateCards(viewCards, cardImages, viewTableStack, fonts)

	// Update bot hands
	cardBackImage := g.AssetManager.GetCardBackImage()
	if cardBackImage == nil {
		cardBackImage = ebiten.NewImage(80, 120)
		cardBackImage.Fill(color.RGBA{100, 100, 100, 255}) // Gray color for card back
	}
	for i, botHand := range botHands {
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

func (g *Game) UpdateTurn() {
	current := g.TurnManager.Current()
	if current == nil {
		return
	}
	if current.IsBot {
		g.AIManager.OnTurn(current.ID, g)
	}
}

// GetPlayerState implements ai.GameLike.
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

// Pass implements ai.GameLike.
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

// PlayCard implements ai.GameLike.
func (g *Game) PlayCard(playerID string, cardID string) error {
	current := g.TurnManager.Current()
	if current == nil || current.ID != playerID || current.Finished {
		fmt.Println("Cannot play card because it's not your turn. Player:", playerID, "Current:", current)
		return fmt.Errorf("Cannot play card because it's not your turn. Player: %s, Current: %v", playerID, current)
	}

	player := g.GetPlayer(playerID)
	err := g.CardManager.PlayCard(player, cardID)
	if err != nil {
		fmt.Println("Error playing card:", err)
		return err
	}

	if err := g.AssetManager.PlaySound(SoundPlay); err != nil {
		fmt.Println("Error playing sound:", err)
	}

	g.TurnManager.MarkAllUnpassed()

	hasDish := false
	for g.CardManager.TryMakeDish() {
		hasDish = true
		for _, p := range g.Players {
			if len(p.Hand) == 0 {
				g.TurnManager.MarkHandEmpty(p.ID)
				if !g.CardManager.TableStack.HasPlayerCards(p.ID) {
					g.TurnManager.MarkFinished(p.ID)
				}
			}
		}
	}

	for _, p := range g.Players {
		if len(p.Hand) == 0 {
			g.TurnManager.MarkHandEmpty(p.ID)
		}
	}

	if !hasDish || len(player.Hand) == 0 {
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
