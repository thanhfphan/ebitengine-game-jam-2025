package ai

type Manager struct {
	bots map[string]Bot // map PlayerID -> Bot
}

func NewManager() *Manager {
	return &Manager{bots: make(map[string]Bot)}
}

func (m *Manager) RegisterBot(playerID string, bot Bot) {
	m.bots[playerID] = bot
}

func (m *Manager) OnTurn(playerID string, g GameLike) {
	if bot, ok := m.bots[playerID]; ok {
		_ = bot.PlayTurn(g, playerID)
	}
}
