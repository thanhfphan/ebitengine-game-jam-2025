package ui

import (
	"image/color"

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
	Image          *ebiten.Image
	BorderColor    color.RGBA
	HoverColor     color.RGBA
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
	tags        Tag
}

func NewUICard(id string, img *ebiten.Image, w, h int) *UICard {
	return &UICard{
		ID:                      id,
		Image:                   img,
		Width:                   w,
		Height:                  h,
		draggable:               false,
		isDragging:              false,
		dragOffsetX:             0,
		dragOffsetY:             0,
		visible:                 true,
		BorderColor:             color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		HoverColor:              color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff},
		SelectedColor:           color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff},
		HighlightColor:          color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff},
		CardType:                "ingredient",
		HighlightedRequirements: make(map[string]bool),
		RequirementNames:        make(map[string]string),
	}
}

func (u *UICard) Draw(screen *ebiten.Image) {
	if !u.visible {
		return
	}

	// Draw card background
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(
		float64(u.Width)/float64(u.Image.Bounds().Dx()),  // scale x
		float64(u.Height)/float64(u.Image.Bounds().Dy()), // scale y
	)
	op.GeoM.Translate(float64(u.X), float64(u.Y))
	screen.DrawImage(u.Image, op)

	// Determine border color
	borderColor := u.BorderColor
	if u.selected {
		borderColor = u.SelectedColor
	} else if u.IsNeededForRecipe {
		borderColor = u.HighlightColor
	} else if u.hovering {
		borderColor = u.HoverColor
	} else if u.CanMakeDish {
		borderColor = u.HighlightColor
	}

	// Draw border
	vector.DrawFilledRect(screen, float32(u.X), float32(u.Y), float32(u.Width), 1, borderColor, false)
	vector.DrawFilledRect(screen, float32(u.X), float32(u.Y), 1, float32(u.Height), borderColor, false)
	vector.DrawFilledRect(screen, float32(u.X+u.Width-1), float32(u.Y), 1, float32(u.Height), borderColor, false)
	vector.DrawFilledRect(screen, float32(u.X), float32(u.Y+u.Height-1), float32(u.Width), 1, borderColor, false)

	// Draw card content
	if u.TitleFont != nil {
		padding := 5

		if u.CardType == "ingredient" {
			iconText := u.Icon
			text.Draw(screen, iconText, u.TitleFont, u.X+u.Width/2-10, u.Y+20, color.White)
			vector.DrawFilledRect(screen, float32(u.X+padding), float32(u.Y+30), float32(u.Width-padding*2), 1, color.White, false)
			if u.Name != "" {
				text.Draw(screen, u.Name, u.BodyFont, u.X+padding, u.Y+50, color.White)
			}
		} else if u.CardType == "recipe" {
			titleY := u.Y + 20
			iconText := u.Icon
			title := iconText + " " + u.Name
			text.Draw(screen, title, u.TitleFont, u.X+padding, titleY, color.White)

			vector.DrawFilledRect(screen, float32(u.X+padding), float32(titleY+5), float32(u.Width-padding*2), 1, color.White, false)

			text.Draw(screen, "Require:", u.SubtitleFont, u.X+padding, titleY+20, color.White)

			reqY := titleY + 35
			for i, reqID := range u.Requirements {
				reqName := u.RequirementNames[reqID]
				if reqName == "" {
					reqName = reqID
				}

				textColor := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
				if u.HighlightedRequirements[reqID] {
					textColor = u.HighlightColor
				}

				text.Draw(screen, "• "+reqName, u.BodyFont, u.X+padding, reqY+i*15, textColor)
			}
		}
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
		if _, has := tableStack.MapIngredients[reqID]; has {
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
		u.X = x - u.dragOffsetX
		u.Y = y - u.dragOffsetY
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
func (u *UICard) GetTags() Tag                { return u.tags }
func (u *UICard) SetTags(t Tag)               { u.tags = t }
func (u *UICard) SetDraggable(draggable bool) { u.draggable = draggable }
func (u *UICard) IsDraggable() bool           { return u.draggable }
func (u *UICard) SetPosition(x, y int) {
	u.X = x
	u.Y = y
}
