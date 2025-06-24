package game

import (
	"github.com/thanhfphan/ebitengj2025/internal/entity"
	"github.com/thanhfphan/ebitengj2025/internal/view"
)

func ToViewCard(card *entity.Card) view.Card {
	cardType := "ingredient"
	if card.Type == entity.CardRecipe {
		cardType = "recipe"
	}

	icon := "+" // TODO: Extract icon from card data

	return view.Card{
		ID:                     card.ID,
		Type:                   cardType,
		Name:                   card.Name,
		Icon:                   icon,
		IngredientID:           card.IngredientID,
		RequiredIngredientIDs:  card.RequiredIngredients,
		CurrentIngredientCount: make(map[string]bool),
	}
}

// ToViewTableStack converts an entity.TableStack to a view.TableStack
func ToViewTableStack(stack *entity.TableStack) view.TableStack {
	result := view.TableStack{
		MapRecipes:     make(map[string]view.Card),
		MapIngredients: make(map[string]view.Card),
		OrderRecipes:   []string{},
	}

	for _, ing := range stack.Ingredients {
		result.MapIngredients[ing.IngredientID] = ToViewCard(ing)
	}

	for _, recipe := range stack.Receipes {
		recipeCard := ToViewCard(recipe)
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
