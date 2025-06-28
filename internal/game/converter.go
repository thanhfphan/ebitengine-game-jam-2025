package game

import (
	"github.com/thanhfphan/ebitengj2025/internal/entity"
	"github.com/thanhfphan/ebitengj2025/internal/view"
)

func ToViewCard(card *entity.Card) view.Card {
	cardType := "ingredient"
	if card.Type == entity.CardTypeRecipe {
		cardType = "recipe"
	}

	return view.Card{
		ID:                     card.ID,
		Type:                   cardType,
		Name:                   card.Name,
		Icon:                   "", // TODO: Extract icon from card data
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
		StackRecipes:   []string{},
	}

	for _, card := range stack.GetCardsByType(entity.CardTypeIngredient) {
		if card.IngredientID == "" {
			// should not happen
			panic("Ingredient card has no ingredient ID")
		}
		result.MapIngredients[card.IngredientID] = ToViewCard(card)
	}

	for _, card := range stack.GetAllCardsInReverseOrder() {
		if card.Type != entity.CardTypeRecipe {
			continue
		}

		recipeCard := ToViewCard(card)
		for _, reqID := range recipeCard.RequiredIngredientIDs {
			if _, has := result.MapIngredients[reqID]; !has {
				continue
			}
			recipeCard.CurrentIngredientCount[reqID] = true
		}
		result.MapRecipes[card.ID] = recipeCard
		result.StackRecipes = append(result.StackRecipes, card.ID)
	}

	return result
}
