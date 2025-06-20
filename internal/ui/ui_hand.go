package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thanhfphan/ebitengj2025/internal/entity"
)

var _ Element = (*UIHand)(nil)

type UIHand struct {
	X, Y          int
	Width, Height int
	Cards         []*UICard
	Spacing       int
	MaxCards      int

	visible      bool
	zIndex       int
	selectedCard *UICard
	onPlayCard   func(cardIndex int)
}

func NewUIHand(x, y, width, height int) *UIHand {
	return &UIHand{
		X:        x,
		Y:        y,
		Width:    width,
		Height:   height,
		Cards:    []*UICard{},
		Spacing:  10,
		MaxCards: 10,
		visible:  true,
		zIndex:   0,
	}
}

func (h *UIHand) SetCards(cards []*entity.Card, cardImages map[string]*ebiten.Image) {
	h.Cards = []*UICard{}

	if len(cards) == 0 {
		return
	}

	cardWidth := 80
	cardHeight := 120

	// Adjust spacing based on number of cards
	availableWidth := h.Width - cardWidth
	if len(cards) > 1 {
		h.Spacing = min(h.Spacing, availableWidth/(len(cards)-1))
	}

	for i, card := range cards {
		img := cardImages[card.ID]
		if img == nil {
			fmt.Println("Card image not found for ID:", card.ID)
			// TODO: Use a default/placeholder image if specific card image not found
			continue
		}

		uiCard := NewUICard(card.ID, img, cardWidth, cardHeight)
		uiCard.X = h.X + i*(cardWidth+h.Spacing)
		uiCard.Y = h.Y

		h.Cards = append(h.Cards, uiCard)
	}
}

func (h *UIHand) Update() {
	if !h.visible {
		return
	}

	for _, card := range h.Cards {
		// Reset vertical position for non-selected cards
		if card != h.selectedCard && card.Y < h.Y {
			card.Y = h.Y
		}
		card.Update()
	}
}

func (h *UIHand) Draw(screen *ebiten.Image) {
	if !h.visible {
		return
	}

	// Draw cards from right to left to ensure proper overlap
	for i := len(h.Cards) - 1; i >= 0; i-- {
		h.Cards[i].Draw(screen)
	}
}

func (h *UIHand) HandleMouseDown(x, y int) bool {
	if !h.visible {
		return false
	}

	// Check cards in reverse order (top to bottom visually)
	for i := 0; i < len(h.Cards); i++ {
		card := h.Cards[i]
		if card.Contains(x, y) {
			if h.selectedCard != nil && h.selectedCard != card {
				// Deselect previous card
				h.selectedCard.Y = h.Y
				h.selectedCard.selected = false
			}

			h.selectedCard = card
			if card.selected {
				card.selected = false // Deselect if already selected
				card.Y = h.Y
			} else {
				card.selected = true
				card.Y = h.Y - 20 // Move card up when selected
			}

			return true
		}
	}

	return false
}

func (h *UIHand) HandleMouseUp(x, y int) bool {
	return h.visible
}

func (h *UIHand) Contains(x, y int) bool {
	if !h.visible {
		return false
	}

	// Check if point is within the hand area
	handArea := y >= h.Y-20 && y <= h.Y+h.Height &&
		x >= h.X && x <= h.X+h.Width

	// Also check individual cards
	for _, card := range h.Cards {
		if card.Contains(x, y) {
			return true
		}
	}

	return handArea
}

func (h *UIHand) GetSelectedCardIndex() int {
	if h.selectedCard == nil {
		return -1
	}

	for i, card := range h.Cards {
		if card == h.selectedCard {
			return i
		}
	}

	return -1
}

func (h *UIHand) SetOnPlayCard(callback func(cardIndex int)) {
	h.onPlayCard = callback
}

func (h *UIHand) PlaySelected() bool {
	idx := h.GetSelectedCardIndex()
	if idx >= 0 && h.onPlayCard != nil {
		h.onPlayCard(idx)
		return true
	}

	return false
}

func (h *UIHand) IsVisible() bool {
	return h.visible
}

func (h *UIHand) SetVisible(v bool) {
	h.visible = v
}

func (h *UIHand) GetZIndex() int {
	return h.zIndex
}

func (h *UIHand) SetZIndex(z int) {
	h.zIndex = z
}
