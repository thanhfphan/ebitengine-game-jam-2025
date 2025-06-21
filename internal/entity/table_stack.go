package entity

type TableStack struct {
	Receipes      []*Card
	Ingredients   []*Card
	PlayerCardMap map[string][]string // playerID -> list of cardIDs
}

func NewTableStack() *TableStack {
	return &TableStack{
		Receipes:      []*Card{},
		Ingredients:   []*Card{},
		PlayerCardMap: make(map[string][]string),
	}
}

func (t *TableStack) Clear() {
	t.Receipes = []*Card{}
	t.Ingredients = []*Card{}
	t.PlayerCardMap = make(map[string][]string)
}

func (t *TableStack) AddCard(card *Card, playerID string) {
	if card.Type == CardRecipe {
		t.AddReceipe(card)
	} else {
		t.AddIngredient(card)
	}
	t.PlayerCardMap[playerID] = append(t.PlayerCardMap[playerID], card.ID)
}

func (t *TableStack) AddReceipe(card *Card) {
	t.Receipes = append(t.Receipes, card)
}

func (t *TableStack) AddIngredient(card *Card) {
	t.Ingredients = append(t.Ingredients, card)
}

func (t *TableStack) RemoveReceipe(card *Card) {
	for i, c := range t.Receipes {
		if c.ID == card.ID {
			t.Receipes = append(t.Receipes[:i], t.Receipes[i+1:]...)
			return
		}
	}
}

func (t *TableStack) RemoveReceipeAt(idx int) *Card {
	if idx < 0 || idx >= len(t.Receipes) {
		return nil
	}

	removeCard := t.Receipes[idx]
	t.Receipes = append(t.Receipes[:idx], t.Receipes[idx+1:]...)
	return removeCard
}

func (t *TableStack) RemoveIngredient(card *Card) {
	for i, c := range t.Ingredients {
		if c.ID == card.ID {
			t.Ingredients = append(t.Ingredients[:i], t.Ingredients[i+1:]...)
			return
		}
	}
}

func (t *TableStack) RemoveIngredientAt(idx int) *Card {
	if idx < 0 || idx >= len(t.Ingredients) {
		return nil
	}

	removeCard := t.Ingredients[idx]
	t.Ingredients = append(t.Ingredients[:idx], t.Ingredients[idx+1:]...)
	return removeCard
}

// HasPlayerCards checks if a player still has any cards on the table
func (t *TableStack) HasPlayerCards(playerID string) bool {
	// Get all card IDs played by this player
	cardIDs, exists := t.PlayerCardMap[playerID]
	if !exists || len(cardIDs) == 0 {
		return false
	}

	// Check if any of those cards are still on the table
	for _, cardID := range cardIDs {
		for _, card := range t.Ingredients {
			if card.ID == cardID {
				return true
			}
		}
	}

	return false
}

func (t *TableStack) RemoveCardFromPlayerTracking(cardID string) {
	for playerID, cardIDs := range t.PlayerCardMap {
		for i, id := range cardIDs {
			if id == cardID {
				t.PlayerCardMap[playerID] = append(cardIDs[:i], cardIDs[i+1:]...)
				break
			}
		}
	}
}
