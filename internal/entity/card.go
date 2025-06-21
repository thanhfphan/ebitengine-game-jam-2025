package entity

type CartType int

const (
	CardIngredient CartType = iota
	CardRecipe
)

type Card struct {
	Entity
	Type                CartType
	IngredientID        string   // If Type is CartIngredient, this is the ID of the ingredient
	RequiredIngredients []string // If Type is CartRecipe, this is the list of required ingredients
}

func NewCard(x, y float64, name string, cartType CartType) *Card {
	entity := NewEntity(TypeCard, name, x, y)
	return &Card{
		Entity: *entity,
		Type:   cartType,
	}
}
