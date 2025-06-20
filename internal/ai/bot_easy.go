package ai

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	mrand "math/rand"
)

var _ Bot = (*EasyBot)(nil)

type EasyBot struct {
	rand *mrand.Rand
}

func NewEasyBot() *EasyBot {
	var seed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &seed)

	return &EasyBot{
		rand: mrand.New(mrand.NewSource(seed)),
	}
}

func (b *EasyBot) PlayTurn(g GameLike, botID string) error {
	player := g.GetPlayerState(botID)
	if player == nil || player.Finished {
		return nil
	}

	if len(player.Hand) == 0 {
		return nil
	}

	idx := rand.Intn(len(player.Hand)) // pick random card

	return g.PlayCard(player.ID, idx)
}
