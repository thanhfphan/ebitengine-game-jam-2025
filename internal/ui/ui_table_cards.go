package ui

import (
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thanhfphan/ebitengj2025/internal/view"
	"golang.org/x/image/font"
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

func (u *UITableCards) ResetCanMakeDish() {
	for _, card := range u.Cards {
		card.CanMakeDish = false
	}
}

func (u *UITableCards) UpdateCanMakeDish(ingredientID string, tableStack view.TableStack) {
	if ingredientID == "" {
		return
	}

	for i := 0; i < len(tableStack.StackRecipes); i++ {
		recipeID := tableStack.StackRecipes[i]
		recipe := tableStack.MapRecipes[recipeID]

		if !slices.Contains(recipe.RequiredIngredientIDs, ingredientID) {
			continue
		}
		// already have this ingredient
		if recipe.CurrentIngredientCount[ingredientID] {
			continue
		}
		count := 0
		for _, has := range recipe.CurrentIngredientCount {
			if has {
				count++
			}
		}
		// not enough cards to make this recipe(need 1 more)
		if len(recipe.RequiredIngredientIDs) != count+1 {
			continue
		}

		for _, card := range u.Cards {
			if card.ID == recipeID {
				card.CanMakeDish = true
			}
		}
		break // Only highlight the topmost recipe
	}

}

// UpdateFromTableStack updates the UI cards based on the view.TableStack
func (u *UITableCards) UpdateFromTableStack(tableStack view.TableStack, cardImages map[string]*ebiten.Image, fonts map[string]font.Face) {
	u.cleanupCards(tableStack)

	// Create a map of ingredient names for recipe requirements
	ingredientNames := make(map[string]string)
	for _, ing := range tableStack.MapIngredients {
		ingredientNames[ing.ID] = ing.Name
	}

	for _, recipe := range tableStack.MapRecipes {
		found := false
		for _, c := range u.Cards {
			if c.ID == recipe.ID {
				found = true
				c.UpdateHightlightingHandRecipes(tableStack)
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
			card.SetCardData(recipe, fonts["title"], fonts["subtitle"], fonts["body"])
			card.SetRequirementNames(ingredientNames)
			card.UpdateHightlightingHandRecipes(tableStack)

			i := len(u.Cards) + 1
			targetX := u.X - 100 + i*40
			targetY := u.Y - u.Radius/2 - 20

			card.X = targetX
			card.Y = targetY

			u.Cards = append(u.Cards, card)
			return
		}
	}

	for _, ingredient := range tableStack.MapIngredients {
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
			card.SetCardData(ingredient, fonts["title"], fonts["subtitle"], fonts["body"])

			i := len(u.Cards) + 1
			targetX := u.X - 80 + i*30
			targetY := u.Y + u.Radius/2 - 60
			card.X = targetX
			card.Y = targetY

			u.Cards = append(u.Cards, card)
			return
		}
	}
}

// cleanupCards removes any cards that are no longer in the table stack
func (u *UITableCards) cleanupCards(tableStack view.TableStack) {
	cardIDs := make(map[string]bool)
	for _, recipe := range tableStack.MapRecipes {
		cardIDs[recipe.ID] = true
	}

	for _, ingredient := range tableStack.MapIngredients {
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
func (u *UITableCards) SetPosition(x, y int) {
	u.X = x
	u.Y = y
}

// Add a method to get all cards
func (u *UITableCards) GetCards() []*UICard {
	return u.Cards
}
