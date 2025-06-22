package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thanhfphan/ebitengj2025/internal/entity"
)

var _ Element = (*UITableCards)(nil)

type UITableCards struct {
	X, Y   int // Center point
	Radius int
	Cards  []*UICard

	BorderColor     color.RGBA
	BackgroundColor color.RGBA

	visible bool
	zIndex  int
	tags    Tag
}

func NewUITableCards(tableX, tableY, tableRadius int) *UITableCards {
	return &UITableCards{
		X:               tableX,
		Y:               tableY,
		Radius:          tableRadius,
		Cards:           []*UICard{},
		visible:         true,
		zIndex:          0,
		BackgroundColor: color.RGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff}, // Gray20
		BorderColor:     color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, // White
	}
}

func (u *UITableCards) Update() {
	if !u.visible {
		return
	}

	for _, card := range u.Cards {
		card.Update()
	}
}

func (u *UITableCards) Draw(screen *ebiten.Image) {
	if !u.visible {
		return
	}

	cx := float32(u.X)
	cy := float32(u.Y)
	r := float32(u.Radius)
	vector.DrawFilledCircle(screen, cx, cy, r, u.BackgroundColor, false)
	vector.StrokeCircle(screen, cx, cy, r, 1, u.BorderColor, false)

	for _, card := range u.Cards {
		card.Draw(screen)
	}
}

func (u *UITableCards) Contains(x, y int) bool {
	if !u.visible {
		return false
	}

	for _, card := range u.Cards {
		if card.Contains(x, y) {
			return true
		}
	}

	return false
}

func (u *UITableCards) HandleMouseDown(x, y int) bool {
	if !u.visible {
		return false
	}

	for _, card := range u.Cards {
		if card.HandleMouseDown(x, y) {
			return true
		}
	}

	return false
}

func (u *UITableCards) HandleMouseUp(x, y int) bool {
	if !u.visible {
		return false
	}

	for _, card := range u.Cards {
		if card.HandleMouseUp(x, y) {
			return true
		}
	}

	return false
}

func (u *UITableCards) UpdateFromTableStack(tableStack *entity.TableStack, cardImages map[string]*ebiten.Image) {
	u.cleanupCards(tableStack)

	for i, recipe := range tableStack.Receipes {
		found := false
		for _, c := range u.Cards {
			if c.ID == recipe.ID {
				found = true
				break
			}
		}

		if !found {
			img := cardImages[recipe.ID]
			if img == nil {
				continue
			}

			card := NewUICard(recipe.ID, img, 80, 120)
			card.SetTags(u.tags)
			card.SetDraggable(true)

			targetX := u.X - 100 + i*40
			targetY := u.Y - u.Radius/2 - 20

			card.X = targetX
			card.Y = targetY

			u.Cards = append(u.Cards, card)
			return
		}
	}

	for i, ingredient := range tableStack.Ingredients {
		found := false
		for _, c := range u.Cards {
			if c.ID == ingredient.ID {
				found = true
				break
			}
		}

		if !found {
			img := cardImages[ingredient.ID]
			if img == nil {
				continue
			}

			card := NewUICard(ingredient.ID, img, 80, 120)
			card.SetTags(u.tags)
			card.SetDraggable(true)

			targetX := u.X - 80 + i*30
			targetY := u.Y + u.Radius/2 - 60
			card.X = targetX
			card.Y = targetY

			u.Cards = append(u.Cards, card)
			return
		}
	}
}

// cleanupCards removes any cards that are no longer in the table stack.
// If new cards are added, they will be added in the next call to
func (u *UITableCards) cleanupCards(tableStack *entity.TableStack) {
	cardIDs := make(map[string]bool)
	for _, recipe := range tableStack.Receipes {
		cardIDs[recipe.ID] = true
	}

	for _, ingredient := range tableStack.Ingredients {
		cardIDs[ingredient.ID] = true
	}

	newCards := make([]*UICard, 0)
	for _, card := range u.Cards {
		if cardIDs[card.ID] {
			newCards = append(newCards, card)
		}
	}

	u.Cards = newCards
}

func (u *UITableCards) IsVisible() bool             { return u.visible }
func (u *UITableCards) SetVisible(v bool)           { u.visible = v }
func (u *UITableCards) GetZIndex() int              { return u.zIndex }
func (u *UITableCards) SetZIndex(z int)             { u.zIndex = z }
func (u *UITableCards) IsStatic() bool              { return true }
func (u *UITableCards) GetTags() Tag                { return u.tags }
func (u *UITableCards) SetTags(tag Tag)             { u.tags = tag }
func (u *UITableCards) SetDraggable(draggable bool) {}
