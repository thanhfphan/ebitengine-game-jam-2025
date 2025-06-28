package view

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Card represents a card view that can be used by the UI layer
type Card struct {
	ID    string
	Image *ebiten.Image
	Type  string // "ingredient" or "recipe"
	Name  string
	Icon  string

	IngredientID string // For ingredient cards

	RequiredIngredientIDs  []string        // For recipe cards
	CurrentIngredientCount map[string]bool // For recipe cards
}

// TableStack represents a collection of cards on the table
type TableStack struct {
	MapRecipes         map[string]Card
	MapIngredients     map[string]Card
	MapIngredientsByID map[string]bool
	StackRecipes       []string // Newest(last put on table) to oldest
}
