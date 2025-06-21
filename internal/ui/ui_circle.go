package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var _ Element = (*UITable)(nil)

type UITable struct {
	X, Y   int // Center point
	Radius int

	BorderColor     color.RGBA
	BackgroundColor color.RGBA

	IsHovered bool
	IsPressed bool
	OnClick   func()

	visible bool
	zindex  int
	tags    Tag
}

func NewUITable(x, y, r int) *UITable {
	return &UITable{
		X:               x,
		Y:               y,
		Radius:          r,
		visible:         true,
		BackgroundColor: color.RGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff}, // Gray20
		BorderColor:     color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, // White
	}
}

func (u *UITable) Update() {
	if !u.visible {
		return
	}
	cx, cy := ebiten.CursorPosition()
	u.IsHovered = u.Contains(cx, cy)
	u.IsPressed = u.IsHovered && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

func (u *UITable) Draw(screen *ebiten.Image) {
	if !u.visible {
		return
	}
	cx := float32(u.X)
	cy := float32(u.Y)
	r := float32(u.Radius)

	// Fill circle
	vector.DrawFilledCircle(screen, cx, cy, r, u.BackgroundColor, false)
	// Stroke border
	vector.StrokeCircle(screen, cx, cy, r, 1, u.BorderColor, false)
}

func (u *UITable) Contains(x, y int) bool {
	if !u.visible {
		return false
	}
	dx := float64(x - u.X)
	dy := float64(y - u.Y)
	return dx*dx+dy*dy <= float64(u.Radius*u.Radius)
}

func (u *UITable) HandleMouseDown(x, y int) bool {
	return u.Contains(x, y)
}

func (u *UITable) HandleMouseUp(x, y int) bool {
	if u.Contains(x, y) && u.IsPressed {
		if u.OnClick != nil {
			u.OnClick()
		}
		return true
	}
	return false
}

func (u *UITable) GetZIndex() int              { return u.zindex }
func (u *UITable) IsVisible() bool             { return u.visible }
func (u *UITable) SetVisible(v bool)           { u.visible = v }
func (u *UITable) SetZIndex(z int)             { u.zindex = z }
func (b *UITable) IsStatic() bool              { return true }
func (b *UITable) GetTags() Tag                { return b.tags }
func (b *UITable) SetTags(t Tag)               { b.tags = t }
func (b *UITable) SetDraggable(draggable bool) {}
