package ui

import "github.com/hajimehoshi/ebiten/v2"

var _ Element = (*UIImage)(nil)

type UIImage struct {
	X, Y          int
	Width, Height int
	Image         *ebiten.Image
	Visible       bool
	ZIndex        int
}

func NewUIImage(x, y int, image *ebiten.Image) *UIImage {
	width, height := image.Bounds().Dx(), image.Bounds().Dy()
	return &UIImage{
		X:       x,
		Y:       y,
		Width:   width,
		Height:  height,
		Image:   image,
		Visible: true,
		ZIndex:  0,
	}
}

func (u *UIImage) Contains(x int, y int) bool {
	return false
}

// Draw implements Element.
func (u *UIImage) Draw(*ebiten.Image) {
}

// GetZIndex implements Element.
func (u *UIImage) GetZIndex() int {
	return u.ZIndex
}

// HandleMouseDown implements Element.
func (u *UIImage) HandleMouseDown(x int, y int) bool {
	return false
}

// HandleMouseUp implements Element.
func (u *UIImage) HandleMouseUp(x int, y int) bool {
	return false
}

// IsStatic implements Element.
func (u *UIImage) IsStatic() bool {
	return true
}

// IsVisible implements Element.
func (u *UIImage) IsVisible() bool {
	return u.Visible
}

// SetVisible implements Element.
func (u *UIImage) SetVisible(v bool) {
	u.Visible = v
}

// SetZIndex implements Element.
func (u *UIImage) SetZIndex(idx int) {
	u.ZIndex = idx
}

// Update implements Element.
func (u *UIImage) Update() {
}
