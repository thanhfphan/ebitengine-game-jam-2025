package game

type GameState int

const (
	GameStateMainMenu GameState = iota
	GameStatePlaying
	GameStatePaused
	GameStateSettings
	GameStateQuit
)
