package ai

import (
	crand "crypto/rand"
	"encoding/binary"
	mrand "math/rand"
	"time"
)

type Manager struct {
	bots     map[string]Bot // map PlayerID -> Bot
	rand     *mrand.Rand
	thinking map[string]time.Time // map PlayerID -> time when bot started thinking
}

func NewManager() *Manager {
	var seed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &seed)

	return &Manager{
		bots:     make(map[string]Bot),
		rand:     mrand.New(mrand.NewSource(seed)),
		thinking: make(map[string]time.Time),
	}
}

func (m *Manager) RegisterBot(playerID string, bot Bot) {
	m.bots[playerID] = bot
}

func (m *Manager) OnTurn(playerID string, g GameLike) {
	if bot, ok := m.bots[playerID]; ok {
		if startTime, isThinking := m.thinking[playerID]; isThinking {
			thinkDuration := time.Since(startTime)

			// Random thinking time between 1200ms and 2000ms
			minThinkTime := 1200 * time.Millisecond
			maxThinkTime := 2000 * time.Millisecond
			thinkTime := minThinkTime + time.Duration(m.rand.Int63n(int64(maxThinkTime-minThinkTime)))

			if thinkDuration < thinkTime {
				return
			}

			delete(m.thinking, playerID)

			_ = bot.PlayTurn(g, playerID)
		} else {
			m.thinking[playerID] = time.Now()
		}
	}
}

func (m *Manager) IsThinking(playerID string) bool {
	_, isThinking := m.thinking[playerID]
	return isThinking
}
