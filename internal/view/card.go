package view

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thanhfphan/ebitengj2025/internal/entity"
)

// Card represents a card view that can be used by the UI layer
type Card struct {
	ID            string
	Image         *ebiten.Image
	Type          entity.CartType
	RequiredCards []string // For recipe cards
}

// TableStack represents a collection of cards on the table
type TableStack struct {
	Recipes     []Card
	Ingredients []Card
}
