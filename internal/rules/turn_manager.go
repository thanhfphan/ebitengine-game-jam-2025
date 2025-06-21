package rules

import (
	"errors"
)

type PlayerTurn struct {
	ID       string
	IsBot    bool
	Passed   bool
	Finished bool
}

type TurnManager struct {
	players []*PlayerTurn
	index   int
	order   []string // finished order
}

func NewTurnManager() *TurnManager {
	return &TurnManager{
		players: []*PlayerTurn{},
		index:   0,
		order:   []string{},
	}
}

func (tm *TurnManager) AddPlayer(id string, isBot bool) {
	tm.players = append(tm.players, &PlayerTurn{
		ID:    id,
		IsBot: isBot,
	})
}

func (tm *TurnManager) Reset() {
	tm.index = 0
	tm.order = []string{}
	tm.players = []*PlayerTurn{}
}

func (tm *TurnManager) Current() *PlayerTurn {
	if len(tm.players) == 0 {
		return nil
	}
	return tm.players[tm.index%len(tm.players)]
}

func (tm *TurnManager) MarkFinished(playerID string) {
	for _, p := range tm.players {
		if p.ID == playerID && !p.Finished {
			p.Finished = true
			tm.order = append(tm.order, p.ID)
			break
		}
	}
}

func (tm *TurnManager) FinishedOrder() []string {
	return tm.order
}

// Next moves to the next player in turn (skipping finished players).
func (tm *TurnManager) Next() {
	for i := 1; i <= len(tm.players); i++ {
		idx := (tm.index + i) % len(tm.players)
		p := tm.players[idx]
		if p.Finished || p.Passed {
			continue
		}
		tm.index = idx
		return
	}

	// If all players have passed, reset the index to the first player
	//tm.index = tm.index
}

func (tm *TurnManager) Pass(playerID string) error {
	for _, p := range tm.players {
		if p.ID == playerID && !p.Finished {
			p.Passed = true
			return nil
		}
	}
	return errors.New("invalid pass")
}

func (tm *TurnManager) MarkAllUnpassed() {
	for _, p := range tm.players {
		p.Passed = false
	}
}

func (tm *TurnManager) GetPlayerByID(id string) *PlayerTurn {
	for _, p := range tm.players {
		if p.ID == id {
			return p
		}
	}
	return nil
}
