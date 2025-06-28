package ui

import (
	"image/color"
	"math"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
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
	BackgroundImage *ebiten.Image

	needsResolveOverlaps bool

	visible bool
	zIndex  int
}

func NewUITableCards(tableX, tableY, tableRadius int, backgroundImage *ebiten.Image) *UITableCards {
	return &UITableCards{
		X:               tableX,
		Y:               tableY,
		Radius:          tableRadius,
		Cards:           []*UICard{},
		visible:         true,
		zIndex:          0,
		BackgroundColor: color.RGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff}, // Gray20
		BorderColor:     color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, // White
		BackgroundImage: backgroundImage,
	}
}

func (u *UITableCards) Update() {
	if !u.visible {
		return
	}

	for _, card := range u.Cards {
		card.Update()

		// Constrain card position to stay within table bounds
		if card.IsDraggable() && card.isDragging {
			dx := float64(card.X + card.Width/2 - u.X)
			dy := float64(card.Y + card.Height/2 - u.Y)
			distance := math.Sqrt(dx*dx + dy*dy)

			maxDistance := float64(u.Radius - card.Width/2)
			if distance > maxDistance {
				dx /= distance
				dy /= distance

				card.X = u.X + int(dx*maxDistance) - card.Width/2
				card.Y = u.Y + int(dy*maxDistance) - card.Height/2
			}
		}
	}

	if u.needsResolveOverlaps {
		u.resolveCardOverlaps()
		u.needsResolveOverlaps = false
	}
}

func (u *UITableCards) Draw(screen *ebiten.Image) {
	if !u.visible {
		return
	}

	cx := float64(u.X)
	cy := float64(u.Y)
	r := float64(u.Radius)

	imgW := u.BackgroundImage.Bounds().Dx()
	imgH := u.BackgroundImage.Bounds().Dy()
	scale := (r * 2) / float64(imgW)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		cx-(float64(imgW)*scale)/2,
		cy-(float64(imgH)*scale)/2,
	)
	screen.DrawImage(u.BackgroundImage, op)

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

	handle := false
	for _, card := range u.Cards {
		if card.HandleMouseUp(x, y) {
			handle = true
		}
	}

	if handle {
		u.needsResolveOverlaps = true
	}

	return handle
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
func (u *UITableCards) UpdateFromTableStack(tableStack view.TableStack, fonts map[string]font.Face, ingredientNames map[string]string) {
	u.cleanupCards(tableStack)

	cardAdded := false

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
			card := NewUICard(recipe.ID, 80, 120)
			card.SetDraggable(true)
			card.SetCardData(recipe, fonts["title"], fonts["subtitle"], fonts["body"])
			card.SetRequirementNames(ingredientNames)
			card.UpdateHightlightingHandRecipes(tableStack)

			// Place new cards in a more distributed way
			angle := rand.Float64() * 2 * math.Pi
			distance := float64(u.Radius) * 0.6 * rand.Float64()

			targetX := u.X + int(math.Cos(angle)*distance) - card.Width/2
			targetY := u.Y + int(math.Sin(angle)*distance) - card.Height/2

			card.X = targetX
			card.Y = targetY

			u.Cards = append(u.Cards, card)
			cardAdded = true
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
			card := NewUICard(ingredient.ID, 80, 120)
			card.SetDraggable(true)
			card.SetCardData(ingredient, fonts["title"], fonts["subtitle"], fonts["body"])

			// Place new cards in a more distributed way
			angle := rand.Float64() * 2 * math.Pi
			distance := float64(u.Radius) * 0.6 * rand.Float64()

			targetX := u.X + int(math.Cos(angle)*distance) - card.Width/2
			targetY := u.Y + int(math.Sin(angle)*distance) - card.Height/2

			card.X = targetX
			card.Y = targetY

			u.Cards = append(u.Cards, card)
			cardAdded = true
		}
	}

	if cardAdded {
		u.needsResolveOverlaps = true
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
func (u *UITableCards) SetDraggable(draggable bool) {}
func (u *UITableCards) SetPosition(x, y int) {
	u.X = x
	u.Y = y
}

// Add a method to get all cards
func (u *UITableCards) GetCards() []*UICard {
	return u.Cards
}

// resolveCardOverlaps prevents cards from completely overlapping
func (u *UITableCards) resolveCardOverlaps() {
	const minVisibleWidth = 30 // Minimum visible width of a card

	for i := 0; i < len(u.Cards); i++ {
		for j := i + 1; j < len(u.Cards); j++ {
			card1 := u.Cards[i]
			card2 := u.Cards[j]

			// Skip if either card is being dragged
			if card1.IsDragging() || card2.IsDragging() {
				continue
			}

			// Calculate overlap
			overlapX := min(card1.X+card1.Width, card2.X+card2.Width) - max(card1.X, card2.X)
			overlapY := min(card1.Y+card1.Height, card2.Y+card2.Height) - max(card1.Y, card2.Y)

			// If there's significant overlap in both dimensions
			if overlapX > card1.Width-minVisibleWidth && overlapY > 0 {
				// Calculate push direction (from center of card1 to center of card2)
				centerX1 := float64(card1.X) + float64(card1.Width)/2
				centerY1 := float64(card1.Y) + float64(card1.Height)/2
				centerX2 := float64(card2.X) + float64(card2.Width)/2
				centerY2 := float64(card2.Y) + float64(card2.Height)/2

				dirX := centerX2 - centerX1
				dirY := centerY2 - centerY1

				// Normalize direction
				length := math.Sqrt(dirX*dirX + dirY*dirY)
				if length > 0 {
					dirX /= length
					dirY /= length
				} else {
					// If centers are exactly the same, push in a random direction
					dirX = 1
					dirY = 0
				}

				// Push cards apart
				pushAmount := float64(minVisibleWidth) / 2

				// Ensure we don't push cards outside the table
				newX1 := int(float64(card1.X) - dirX*pushAmount)
				newY1 := int(float64(card1.Y) - dirY*pushAmount)
				newX2 := int(float64(card2.X) + dirX*pushAmount)
				newY2 := int(float64(card2.Y) + dirY*pushAmount)

				// Check if new positions are within table bounds
				u.constrainCardToTable(card1, newX1, newY1)
				u.constrainCardToTable(card2, newX2, newY2)
			}
		}
	}
}

// constrainCardToTable ensures a card stays within the table bounds
func (u *UITableCards) constrainCardToTable(card *UICard, newX, newY int) {
	// Calculate distance from new position to table center
	centerX := newX + card.Width/2
	centerY := newY + card.Height/2
	dx := float64(centerX - u.X)
	dy := float64(centerY - u.Y)
	distance := math.Sqrt(dx*dx + dy*dy)

	maxDistance := float64(u.Radius - card.Width/2)

	// If new position is within bounds, apply it
	if distance <= maxDistance {
		card.X = newX
		card.Y = newY
	} else {
		// Otherwise, place the card at the boundary
		dx /= distance
		dy /= distance
		card.X = u.X + int(dx*maxDistance) - card.Width/2
		card.Y = u.Y + int(dy*maxDistance) - card.Height/2
	}
}
