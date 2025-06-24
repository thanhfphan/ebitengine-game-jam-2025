package ai

import (
	"github.com/thanhfphan/ebitengj2025/internal/entity"
)

type Bot interface {
	PlayTurn(g GameLike, botID string) error
}

type GameLike interface {
	GetPlayerState(id string) *PlayerState
	PlayCard(playerID string, cardID string) error
	Pass(playerID string)
}

type PlayerState struct {
	ID       string
	Hand     map[string]*entity.Card
	IsBot    bool
	Passed   bool
	Finished bool
}
