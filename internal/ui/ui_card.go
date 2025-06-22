package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var _ Element = (*UICard)(nil)

type UICard struct {
	ID            string // TODO: might want remove this one?
	X, Y          int
	Width, Height int
	Image         *ebiten.Image
	BorderColor   color.RGBA
	HoverColor    color.RGBA
	SelectedColor color.RGBA

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
		ID:            id,
		Image:         img,
		Width:         w,
		Height:        h,
		draggable:     false,
		isDragging:    false,
		dragOffsetX:   0,
		dragOffsetY:   0,
		visible:       true,
		BorderColor:   color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		HoverColor:    color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff},
		SelectedColor: color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff},
	}
}

func (c *UICard) Draw(screen *ebiten.Image) {
	if !c.visible {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(
		float64(c.Width)/float64(c.Image.Bounds().Dx()),  // scale x
		float64(c.Height)/float64(c.Image.Bounds().Dy()), // scale y
	)
	op.GeoM.Translate(float64(c.X), float64(c.Y))
	screen.DrawImage(c.Image, op)

	borderColor := c.BorderColor
	if c.selected {
		borderColor = c.SelectedColor
	} else if c.hovering {
		borderColor = c.HoverColor
	}

	vector.DrawFilledRect(screen, float32(c.X), float32(c.Y), float32(c.Width), 1, borderColor, false)
	vector.DrawFilledRect(screen, float32(c.X), float32(c.Y), 1, float32(c.Height), borderColor, false)
	vector.DrawFilledRect(screen, float32(c.X+c.Width-1), float32(c.Y), 1, float32(c.Height), borderColor, false)
	vector.DrawFilledRect(screen, float32(c.X), float32(c.Y+c.Height-1), float32(c.Width), 1, borderColor, false)
}

func (c *UICard) Update() {
	if !c.visible {
		return
	}

	mx, my := ebiten.CursorPosition()
	c.hovering = c.Contains(mx, my)

	if c.isDragging {
		x, y := ebiten.CursorPosition()
		c.X = x - c.dragOffsetX
		c.Y = y - c.dragOffsetY
	}
}

func (c *UICard) Contains(x, y int) bool {
	return c.visible &&
		x >= c.X && x < c.X+c.Width &&
		y >= c.Y && y < c.Y+c.Height
}

func (c *UICard) HandleMouseDown(x, y int) bool {
	if !c.visible || !c.Contains(x, y) {
		return false
	}

	c.selected = true

	if c.draggable {
		c.isDragging = true
		c.dragOffsetX = x - c.X
		c.dragOffsetY = y - c.Y
	}

	return true
}

func (c *UICard) HandleMouseUp(x, y int) bool {
	if !c.visible {
		return false
	}

	afterMouseUpFnc := func() {
		c.isDragging = false
		c.selected = false
	}

	if c.isDragging || c.selected {
		afterMouseUpFnc()
		return true
	}

	return false
}

func (c *UICard) IsVisible() bool             { return c.visible }
func (c *UICard) SetVisible(v bool)           { c.visible = v }
func (c *UICard) GetZIndex() int              { return c.zIndex }
func (c *UICard) SetZIndex(z int)             { c.zIndex = z }
func (c *UICard) IsStatic() bool              { return false }
func (c *UICard) GetTags() Tag                { return c.tags }
func (c *UICard) SetTags(t Tag)               { c.tags = t }
func (c *UICard) SetDraggable(draggable bool) { c.draggable = draggable }
func (c *UICard) IsDraggable() bool           { return c.draggable }
func (c *UICard) SetPosition(x, y int) {
	c.X = x
	c.Y = y
}
