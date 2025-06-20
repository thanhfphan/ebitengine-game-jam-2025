package ai

import (
	"github.com/thanhfphan/ebitengj2025/internal/entity"
)

type Bot interface {
	PlayTurn(g GameLike, botID string) error
}

type GameLike interface {
	GetPlayerState(id string) *PlayerState
	PlayCard(playerID string, index int) error
	Pass(playerID string)
}

type PlayerState struct {
	ID       string
	Hand     []*entity.Card
	IsBot    bool
	Passed   bool
	Finished bool
}
