package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var _ Element = (*UICircle)(nil)

type UICircle struct {
	X, Y   int // Center point
	Radius int

	Visible bool
	ZIndex  int

	BorderColor     color.RGBA
	BackgroundColor color.RGBA

	IsHovered bool
	IsPressed bool
	OnClick   func()
}

func NewUICircle(x, y, r int) *UICircle {
	return &UICircle{
		X:               x,
		Y:               y,
		Radius:          r,
		Visible:         true,
		BackgroundColor: color.RGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff}, // Gray20
		BorderColor:     color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, // White
	}
}

func (u *UICircle) Update() {
	if !u.Visible {
		return
	}
	cx, cy := ebiten.CursorPosition()
	u.IsHovered = u.Contains(cx, cy)
	u.IsPressed = u.IsHovered && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

func (u *UICircle) Draw(screen *ebiten.Image) {
	if !u.Visible {
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

func (u *UICircle) Contains(x, y int) bool {
	if !u.Visible {
		return false
	}
	dx := float64(x - u.X)
	dy := float64(y - u.Y)
	return dx*dx+dy*dy <= float64(u.Radius*u.Radius)
}

func (u *UICircle) HandleMouseDown(x, y int) bool {
	return u.Contains(x, y)
}

func (u *UICircle) HandleMouseUp(x, y int) bool {
	if u.Contains(x, y) && u.IsPressed {
		if u.OnClick != nil {
			u.OnClick()
		}
		return true
	}
	return false
}

func (u *UICircle) GetZIndex() int {
	return u.ZIndex
}

func (u *UICircle) IsVisible() bool {
	return u.Visible
}

func (u *UICircle) SetVisible(v bool) {
	u.Visible = v
}

func (u *UICircle) SetZIndex(z int) {
	u.ZIndex = z
}
