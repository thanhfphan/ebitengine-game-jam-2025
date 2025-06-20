package entity

type TableStack struct {
	Receipes    []*Card
	Ingredients []*Card
}

func NewTableStack() *TableStack {
	return &TableStack{
		Receipes:    []*Card{},
		Ingredients: []*Card{},
	}
}

func (t *TableStack) Clear() {
	t.Receipes = []*Card{}
	t.Ingredients = []*Card{}
}

func (t *TableStack) AddCard(card *Card) {
	if card.Type == CardRecipe {
		t.AddReceipe(card)
	} else {
		t.AddIngredient(card)
	}
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
