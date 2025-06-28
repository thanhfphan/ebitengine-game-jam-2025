package entity

type CartType int

const (
	CardTypeIngredient CartType = iota
	CardTypeRecipe
)

type Card struct {
	Entity
	Type                CartType
	IngredientID        string   // If Type is CartIngredient, this is the ID of the ingredient
	RequiredIngredients []string // If Type is CartRecipe, this is the list of required ingredients
}

func NewCard(name string, cartType CartType) *Card {
	entity := NewEntity(TypeCard, name)
	return &Card{
		Entity: *entity,
		Type:   cartType,
	}
}
