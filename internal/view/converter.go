package view

import (
	"github.com/thanhfphan/ebitengj2025/internal/entity"
)

func FromEntityCard(card *entity.Card) Card {
	return Card{
		ID:            card.ID,
		Type:          card.Type,
		RequiredCards: card.RequiredIngredients,
	}
}

// FromEntityTableStack converts an entity.TableStack to a view.TableStack
func FromEntityTableStack(stack *entity.TableStack) TableStack {
	result := TableStack{
		Recipes:     make([]Card, 0, len(stack.Receipes)),
		Ingredients: make([]Card, 0, len(stack.Ingredients)),
	}

	for _, recipe := range stack.Receipes {
		result.Recipes = append(result.Recipes, FromEntityCard(recipe))
	}

	for _, ingredient := range stack.Ingredients {
		result.Ingredients = append(result.Ingredients, FromEntityCard(ingredient))
	}

	return result
}
