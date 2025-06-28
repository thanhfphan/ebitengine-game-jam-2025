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
	"github.com/thanhfphan/ebitengj2025/internal/rules"
	"github.com/thanhfphan/ebitengj2025/internal/ui"
)

var (
	_ ebiten.Game  = (*Game)(nil)
	_ ai.GameLike  = (*Game)(nil)
	_ SceneManager = (*Game)(nil)
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
	SoundDrag       = "sound_drag"

	ImageMainBG      = "main_bg"
	ImagePlayBG      = "play_bg"
	ImageTableBG     = "table_bg"
	ImageCardBack    = "card_back"
	ImageSettingIcon = "setting_icon"
)

type Game struct {
	State     GameState
	Players   []*entity.Player
	Player    *entity.Player
	DebugMode bool

	AssetManager     *am.AssetManager
	CurrentUIManager *ui.Manager
	AIManager        *ai.Manager
	CardManager      *card.Manager
	TurnManager      *rules.TurnManager

	sceneStack []Scene // Scene stack for managing scenes
}

func New() (*Game, error) {
	assetManager := am.NewAssetManager()
	aiManager := ai.NewManager()
	cardManager := card.NewManager()
	turnManager := rules.NewTurnManager()

	g := &Game{
		State:        GameStateNormal,
		Players:      []*entity.Player{},
		AssetManager: assetManager,
		AIManager:    aiManager,
		CardManager:  cardManager,
		TurnManager:  turnManager,
		sceneStack:   []Scene{},
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
	if err := g.AssetManager.LoadImage(ImageMainBG, "assets/images/backgrounds/main_bg.jpg"); err != nil {
		fmt.Println("Error loading background image:", err)
	}
	if err := g.AssetManager.LoadImage(ImagePlayBG, "assets/images/backgrounds/tet.jpg"); err != nil {
		fmt.Println("Error loading background image:", err)
	}
	if err := g.AssetManager.LoadImage(ImageTableBG, "assets/images/backgrounds/tablecard.png"); err != nil {
		fmt.Println("Error loading background image:", err)
	}
	if err := g.AssetManager.LoadImage(ImageCardBack, "assets/images/backgrounds/card_back.png"); err != nil {
		fmt.Println("Error loading background image:", err)
	}
	if err := g.AssetManager.LoadImage(ImageSettingIcon, "assets/images/ui/setting_icon.png"); err != nil {
		fmt.Println("Error loading setting icon:", err)
	}

	g.PushScene(NewMainMenuScene())

	return g, nil
}

func (g *Game) setupGameData(botCount int) []*ui.UIBotHand {
	g.Players = []*entity.Player{}
	botHands := []*ui.UIBotHand{}
	g.CardManager.LoadDeck("default")
	g.TurnManager.Reset()

	g.Player = entity.NewPlayer("P0", entity.TypePlayer)
	g.Players = append(g.Players, g.Player)
	g.TurnManager.AddPlayer(g.Player.ID, false)

	defaultFont := g.AssetManager.GetFont("nunito", 32)

	for i := 1; i <= botCount; i++ {
		bot := entity.NewPlayer(fmt.Sprintf("B%d", i), entity.TypeBot)
		g.Players = append(g.Players, bot)
		g.TurnManager.AddPlayer(bot.ID, true)
		g.AIManager.RegisterBot(bot.ID, ai.NewEasyBot())

		botHand := ui.NewUIBotHand(0, 0, CardWidth, CardHeight, defaultFont) // Position will be set later
		botHands = append(botHands, botHand)
		g.CurrentUIManager.AddElement(botHand)
	}

	g.CardManager.DealHands(g.Players)

	return botHands
}

// Update implements ebiten.Game.
func (g *Game) Update() error {
	g.HandleInput()

	currentScene := g.CurrentScene()
	if currentScene != nil {
		currentScene.Update(g)
	}

	if g.CurrentUIManager != nil {
		g.CurrentUIManager.Update()
	}

	g.AssetManager.Update()

	if g.State == GameStateQuit {
		return ebiten.Termination
	}

	return nil
}

// Draw implements ebiten.Game.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x00, 0x00, 0xFF})

	currentScene := g.CurrentScene()
	if currentScene != nil {
		currentScene.Draw(screen, g)
	}

	if g.CurrentUIManager != nil {
		g.CurrentUIManager.Draw(screen)
	}
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
		if g.CurrentUIManager != nil {
			g.CurrentUIManager.HandleMouseDown(mx, my)
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if g.CurrentUIManager != nil {
			g.CurrentUIManager.HandleMouseUp(mx, my)
		}
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

// PushScene adds a new scene to the top of the stack
func (g *Game) PushScene(scene Scene) {
	g.sceneStack = append(g.sceneStack, scene)
	scene.Enter(g)
}

// PopScene removes the current scene from the top of the stack
func (g *Game) PopScene() {
	if len(g.sceneStack) == 0 {
		return
	}

	currentScene := g.sceneStack[len(g.sceneStack)-1]
	currentScene.Exit(g)

	g.sceneStack = g.sceneStack[:len(g.sceneStack)-1]

	if len(g.sceneStack) > 0 {
		g.sceneStack[len(g.sceneStack)-1].Enter(g)
	}
}

// ReplaceScene replaces the current scene with a new one
func (g *Game) ReplaceScene(scene Scene) {
	if len(g.sceneStack) > 0 {
		currentScene := g.sceneStack[len(g.sceneStack)-1]
		currentScene.Exit(g)

		g.sceneStack[len(g.sceneStack)-1] = scene
	} else {
		g.sceneStack = append(g.sceneStack, scene)
	}

	scene.Enter(g)
}

func (g *Game) CurrentScene() Scene {
	if len(g.sceneStack) == 0 {
		return nil
	}
	return g.sceneStack[len(g.sceneStack)-1]
}
