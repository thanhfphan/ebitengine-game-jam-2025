package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thanhfphan/ebitengj2025/internal/view"
	"golang.org/x/image/font"
)

var _ Element = (*UICard)(nil)

type UICard struct {
	ID             string
	X, Y           int
	Width, Height  int
	BorderColor    color.RGBA
	SelectedColor  color.RGBA
	HighlightColor color.RGBA
	CanMakeDish    bool

	CardType         string // "ingredient" or "recipe"
	Name             string
	Icon             string
	Requirements     []string
	RequirementNames map[string]string // Map of requirement ID to name

	HighlightedRequirements map[string]bool // Requirements that are available on the table
	IsNeededForRecipe       bool            // If this ingredient is needed for a recipe on the table

	// Fonts
	TitleFont    font.Face
	SubtitleFont font.Face
	BodyFont     font.Face

	draggable   bool
	isDragging  bool
	dragOffsetX int
	dragOffsetY int
	visible     bool
	zIndex      int
	hovering    bool
	selected    bool
}

func NewUICard(id string, w, h int) *UICard {
	return &UICard{
		ID:                      id,
		Width:                   w,
		Height:                  h,
		draggable:               false,
		isDragging:              false,
		dragOffsetX:             0,
		dragOffsetY:             0,
		visible:                 true,
		BorderColor:             color.RGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF}, // Xám nhạt #AAAAAA
		SelectedColor:           color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}, // Green
		HighlightColor:          color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, // Red
		CardType:                "ingredient",
		HighlightedRequirements: make(map[string]bool),
		RequirementNames:        make(map[string]string),
	}
}

func (u *UICard) Draw(screen *ebiten.Image) {
	if !u.visible {
		return
	}

	bgColor := color.RGBA{0xFA, 0xF8, 0xF0, 0xFF} // Ingredient: #FAF8F0
	if u.CardType == "recipe" {
		bgColor = color.RGBA{0xFF, 0xF5, 0xCC, 0xFF} // Recipe: #FFF5CC
	}

	borderColor := u.BorderColor
	if u.selected {
		borderColor = u.SelectedColor
	} else if u.IsNeededForRecipe || u.CanMakeDish {
		borderColor = u.HighlightColor
	}
	// else if u.hovering {
	// }

	textColor := color.RGBA{0x44, 0x44, 0x44, 0xFF}     // #444444
	titleColor := color.RGBA{0x33, 0x33, 0x33, 0xFF}    // #333333
	highlightText := color.RGBA{0xD6, 0x86, 0x00, 0xFF} // #D68600

	x, y := float32(u.X), float32(u.Y)
	w, h := float32(u.Width), float32(u.Height)
	radius := float32(6)
	borderWidth := float32(2)

	vector.DrawFilledRect(screen, x+radius, y, w-radius*2, h, bgColor, false)
	vector.DrawFilledRect(screen, x, y+radius, w, h-radius*2, bgColor, false)
	vector.DrawFilledCircle(screen, x+radius, y+radius, radius, bgColor, false)
	vector.DrawFilledCircle(screen, x+w-radius, y+radius, radius, bgColor, false)
	vector.DrawFilledCircle(screen, x+radius, y+h-radius, radius, bgColor, false)
	vector.DrawFilledCircle(screen, x+w-radius, y+h-radius, radius, bgColor, false)

	vector.StrokeLine(screen, x+radius, y, x+w-radius, y, borderWidth, borderColor, false)     // top
	vector.StrokeLine(screen, x+w, y+radius, x+w, y+h-radius, borderWidth, borderColor, false) // right
	vector.StrokeLine(screen, x+radius, y+h, x+w-radius, y+h, borderWidth, borderColor, false) // bottom
	vector.StrokeLine(screen, x, y+radius, x, y+h-radius, borderWidth, borderColor, false)     // left

	drawArc(screen, x+radius, y+radius, radius, math.Pi, 1.5*math.Pi, borderWidth, borderColor)     // top-left
	drawArc(screen, x+w-radius, y+radius, radius, 1.5*math.Pi, 2*math.Pi, borderWidth, borderColor) // top-right
	drawArc(screen, x+w-radius, y+h-radius, radius, 0, 0.5*math.Pi, borderWidth, borderColor)       // bottom-right
	drawArc(screen, x+radius, y+h-radius, radius, 0.5*math.Pi, math.Pi, borderWidth, borderColor)   // bottom-left

	if u.TitleFont == nil {
		return
	}

	padding := 5
	if u.CardType == "ingredient" {
		text.Draw(screen, u.Icon, u.TitleFont, u.X+u.Width/2-10, u.Y+20, titleColor)
		vector.DrawFilledRect(screen, float32(u.X+padding), float32(u.Y+30), float32(u.Width-padding*2), 1, titleColor, false)
		if u.Name != "" {
			text.Draw(screen, u.Name, u.BodyFont, u.X+padding, u.Y+50, textColor)
		}
	} else if u.CardType == "recipe" {
		titleY := u.Y + 20
		title := u.Icon + " " + u.Name
		text.Draw(screen, title, u.TitleFont, u.X+padding, titleY, titleColor)

		vector.DrawFilledRect(screen, float32(u.X+padding), float32(titleY+5), float32(u.Width-padding*2), 1, titleColor, false)

		text.Draw(screen, "Require:", u.SubtitleFont, u.X+padding, titleY+20, textColor)

		reqY := titleY + 35
		for i, reqID := range u.Requirements {
			reqName := u.RequirementNames[reqID]
			if reqName == "" {
				reqName = reqID
			}
			col := textColor
			if u.HighlightedRequirements[reqID] {
				col = highlightText
			}
			text.Draw(screen, "- "+reqName, u.BodyFont, u.X+padding, reqY+i*15, col)
		}
	}
}

func drawArc(screen *ebiten.Image, cx, cy, r, start, end, width float32, col color.Color) {
	const segments = 10
	thetaStep := (end - start) / segments

	for i := 0; i < segments; i++ {
		theta1 := float64(start) + float64(thetaStep)*float64(i)
		theta2 := float64(start) + float64(thetaStep)*float64(i+1)

		x1 := float64(cx) + math.Cos(theta1)*float64(r)
		y1 := float64(cy) + math.Sin(theta1)*float64(r)
		x2 := float64(cx) + math.Cos(theta2)*float64(r)
		y2 := float64(cy) + math.Sin(theta2)*float64(r)

		vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), float32(width), col, false)
	}
}

// SetCardData sets the card data based on the view.Card
func (u *UICard) SetCardData(card view.Card, titleFont, subtitleFont, bodyFont font.Face) {
	u.CardType = card.Type
	u.Name = card.Name
	u.Icon = card.Icon
	u.Requirements = card.RequiredIngredientIDs
	u.TitleFont = titleFont
	u.SubtitleFont = subtitleFont
	u.BodyFont = bodyFont
}

// SetRequirementNames sets the names for requirement IDs
func (u *UICard) SetRequirementNames(names map[string]string) {
	u.RequirementNames = names
}

func (u *UICard) UpdateHightlightingHandRecipes(tableStack view.TableStack) {
	if u.CardType != "recipe" {
		return
	}

	// Clear previous highlighting
	for _, reqID := range u.Requirements {
		u.HighlightedRequirements[reqID] = false
	}

	for _, reqID := range u.Requirements {
		if _, has := tableStack.MapIngredientsByID[reqID]; has {
			u.HighlightedRequirements[reqID] = true
		}
	}
}

func (u *UICard) Update() {
	if !u.visible {
		return
	}

	mx, my := ebiten.CursorPosition()
	u.hovering = u.Contains(mx, my)

	if u.isDragging {
		x, y := ebiten.CursorPosition()
		newX := x - u.dragOffsetX
		newY := y - u.dragOffsetY

		// The parent component will handle the constraint logic
		u.X = newX
		u.Y = newY
	}
}

func (u *UICard) Contains(x, y int) bool {
	return u.visible &&
		x >= u.X && x < u.X+u.Width &&
		y >= u.Y && y < u.Y+u.Height
}

func (u *UICard) HandleMouseDown(x, y int) bool {
	if !u.visible || !u.Contains(x, y) {
		return false
	}

	u.selected = true

	if u.draggable {
		u.isDragging = true
		u.dragOffsetX = x - u.X
		u.dragOffsetY = y - u.Y
	}

	return true
}

func (u *UICard) HandleMouseUp(x, y int) bool {
	if !u.visible {
		return false
	}

	afterMouseUpFnc := func() {
		u.isDragging = false
		u.selected = false
	}

	if u.isDragging || u.selected {
		afterMouseUpFnc()
		return true
	}

	return false
}

func (u *UICard) IsVisible() bool             { return u.visible }
func (u *UICard) SetVisible(v bool)           { u.visible = v }
func (u *UICard) GetZIndex() int              { return u.zIndex }
func (u *UICard) SetZIndex(z int)             { u.zIndex = z }
func (u *UICard) IsStatic() bool              { return false }
func (u *UICard) SetDraggable(draggable bool) { u.draggable = draggable }
func (u *UICard) IsDraggable() bool           { return u.draggable }
func (u *UICard) SetPosition(x, y int) {
	u.X = x
	u.Y = y
}

func (u *UICard) IsDragging() bool { return u.isDragging }
