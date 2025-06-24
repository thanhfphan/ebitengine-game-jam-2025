package view

import (
	"github.com/thanhfphan/ebitengj2025/internal/entity"
)

func FromEntityCard(card *entity.Card) Card {
	cardType := "ingredient"
	if card.Type == entity.CardRecipe {
		cardType = "recipe"
	}

	icon := "+" // TODO: Extract icon from card data

	return Card{
		ID:                     card.ID,
		Type:                   cardType,
		Name:                   card.Name,
		Icon:                   icon,
		IngredientID:           card.IngredientID,
		RequiredIngredientIDs:  card.RequiredIngredients,
		CurrentIngredientCount: make(map[string]bool),
	}
}

// FromEntityTableStack converts an entity.TableStack to a view.TableStack
func FromEntityTableStack(stack *entity.TableStack) TableStack {
	result := TableStack{
		MapRecipes:     make(map[string]Card),
		MapIngredients: make(map[string]Card),
		OrderRecipes:   []string{},
	}

	for _, ing := range stack.Ingredients {
		result.MapIngredients[ing.IngredientID] = FromEntityCard(ing)
	}

	for _, recipe := range stack.Receipes {
		recipeCard := FromEntityCard(recipe)
		for _, reqID := range recipeCard.RequiredIngredientIDs {
			if _, has := result.MapIngredients[reqID]; !has {
				continue
			}
			recipeCard.CurrentIngredientCount[reqID] = true
		}
		result.MapRecipes[recipe.ID] = recipeCard
	}

	for i := len(stack.Receipes) - 1; i >= 0; i-- {
		recipe := stack.Receipes[i]
		result.OrderRecipes = append(result.OrderRecipes, recipe.ID)
	}

	return result
}
