package ai

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	mrand "math/rand"

	"github.com/thanhfphan/ebitengj2025/internal/entity"
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
	var card *entity.Card
	i := 0
	for _, c := range player.Hand {
		if i == idx {
			card = c
			break
		}
		i++
	}

	return g.PlayCard(player.ID, card.ID)
}
