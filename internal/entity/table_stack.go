package entity

type CardOnTable struct {
	Card     *Card
	PlayerID string
}

type TableStack struct {
	cards     map[string]*CardOnTable // cardID -> CardOnTable
	playOrder []string                // Oldest to newest(last put on table)
}

func NewTableStack() *TableStack {
	return &TableStack{
		cards:     make(map[string]*CardOnTable),
		playOrder: []string{},
	}
}

func (t *TableStack) Clear() {
	t.cards = make(map[string]*CardOnTable)
}

func (t *TableStack) AddCard(card *Card, playerID string) {
	if _, ok := t.cards[card.ID]; ok {
		return
	}

	t.cards[card.ID] = &CardOnTable{
		Card:     card,
		PlayerID: playerID,
	}
	t.playOrder = append(t.playOrder, card.ID)
}

func (t *TableStack) RemoveCard(cardID string) {
	delete(t.cards, cardID)
	for i, id := range t.playOrder {
		if id == cardID {
			t.playOrder = append(t.playOrder[:i], t.playOrder[i+1:]...)
			break
		}
	}
}

// GetAllCardsInOrder returns all cards in the order they were played
func (t *TableStack) GetAllCardsInOrder() []*Card {
	var result []*Card
	for _, id := range t.playOrder {
		if cardEntry, ok := t.cards[id]; ok {
			result = append(result, cardEntry.Card)
		}
	}
	return result
}

// GetAllCardsInReverseOrder returns all cards in the reverse order they were played
func (t *TableStack) GetAllCardsInReverseOrder() []*Card {
	var result []*Card
	for i := len(t.playOrder) - 1; i >= 0; i-- {
		if cardEntry, ok := t.cards[t.playOrder[i]]; ok {
			result = append(result, cardEntry.Card)
		}
	}
	return result
}

func (t *TableStack) GetCardsByType(cardType CartType) []*Card {
	var result []*Card
	for _, c := range t.cards {
		if c.Card.Type == cardType {
			result = append(result, c.Card)
		}
	}
	return result
}

func (t *TableStack) GetCardsByPlayer(playerID string) []*Card {
	var result []*Card
	for _, c := range t.cards {
		if c.PlayerID == playerID {
			result = append(result, c.Card)
		}
	}
	return result
}

func (t *TableStack) HasPlayerCards(playerID string) bool {
	for _, c := range t.cards {
		if c.PlayerID == playerID {
			return true
		}
	}
	return false
}
