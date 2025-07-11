package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thanhfphan/ebitengj2025/internal/view"
	"golang.org/x/image/font"
)

var _ Element = (*UIHand)(nil)

type UIHand struct {
	X, Y           int
	Width, Height  int
	Cards          []*UICard
	selectedCard   *UICard
	onPlayCard     func(cardID string)
	onCardSelected func(cardID string) // New callback for card selection
	visible        bool
	zIndex         int
}

func NewUIHand(x, y, w, h int) *UIHand {
	return &UIHand{
		X:            x,
		Y:            y,
		Width:        w,
		Height:       h,
		Cards:        []*UICard{},
		selectedCard: nil,
		visible:      true,
		zIndex:       0,
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

	for i := 0; i < len(h.Cards); i++ {
		h.Cards[i].Draw(screen)
	}
}

func (h *UIHand) HandleMouseDown(x, y int) bool {
	if !h.visible {
		return false
	}

	// Check cards in reverse order (top to bottom visually)
	for i := len(h.Cards) - 1; i >= 0; i-- {
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
				h.selectedCard = nil

				// Notify that no card is selected
				if h.onCardSelected != nil {
					h.onCardSelected("")
				}
			} else {
				card.selected = true
				card.Y = h.Y - 20 // Move card up when selected

				// Notify that a card is selected
				if h.onCardSelected != nil {
					h.onCardSelected(card.ID)
				}
			}

			return true
		}
	}

	return false
}

func (h *UIHand) HandleMouseUp(x, y int) bool {
	return false // No action on mouse up for hand
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

func (h *UIHand) GetSelectedCardID() string {
	if h.selectedCard == nil {
		return ""
	}

	return h.selectedCard.ID
}

func (h *UIHand) SetOnPlayCard(callback func(cardID string)) {
	h.onPlayCard = callback
}

func (h *UIHand) PlaySelected() bool {
	id := h.GetSelectedCardID()
	if id != "" && h.onPlayCard != nil {
		h.onPlayCard(id)
		return true
	}

	return false
}

func (h *UIHand) UpdateCards(cards []view.Card, tableStack view.TableStack, fonts map[string]font.Face, ingredientNames map[string]string) {
	if len(cards) == 0 {
		h.Cards = []*UICard{}
		h.selectedCard = nil
		return
	}

	cardWidth := 80
	cardHeight := 120

	existingCards := make(map[string]*UICard)
	for _, card := range h.Cards {
		existingCards[card.ID] = card
	}

	newCards := make([]*UICard, 0, len(cards))
	for i, card := range cards {
		var uiCard *UICard

		if existing, ok := existingCards[card.ID]; ok {
			uiCard = existing
			uiCard.UpdateHightlightingHandRecipes(tableStack)
		} else {
			uiCard = NewUICard(card.ID, cardWidth, cardHeight)
			uiCard.SetCardData(card, fonts["title"], fonts["subtitle"], fonts["body"])
			uiCard.SetRequirementNames(ingredientNames)
			uiCard.UpdateHightlightingHandRecipes(tableStack)
		}

		overlap := 30
		uiCard.X = h.X + i*cardWidth + overlap
		if i > 0 {
			uiCard.X -= i * overlap
		}

		if uiCard != h.selectedCard {
			uiCard.Y = h.Y
		}

		newCards = append(newCards, uiCard)
	}

	if h.selectedCard != nil {
		stillExists := false
		for _, card := range newCards {
			if card.ID == h.selectedCard.ID {
				stillExists = true
				break
			}
		}
		if !stillExists {
			h.selectedCard = nil
		}
	}

	h.Cards = newCards
}

func (h *UIHand) IsVisible() bool             { return h.visible }
func (h *UIHand) SetVisible(v bool)           { h.visible = v }
func (h *UIHand) GetZIndex() int              { return h.zIndex }
func (h *UIHand) SetZIndex(z int)             { h.zIndex = z }
func (h *UIHand) IsStatic() bool              { return false }
func (h *UIHand) SetDraggable(draggable bool) {}
func (h *UIHand) SetPosition(x, y int) {
	h.X = x
	h.Y = y
}

func (h *UIHand) SetOnCardSelected(callback func(cardID string)) {
	h.onCardSelected = callback
}
