package ui

import "github.com/hajimehoshi/ebiten/v2"

var _ Element = (*UIImage)(nil)

type UIImage struct {
	X, Y          int
	Width, Height int
	Image         *ebiten.Image

	visible bool
	zindex  int
	tags    Tag
}

func NewUIImage(x, y int, image *ebiten.Image) *UIImage {
	width, height := image.Bounds().Dx(), image.Bounds().Dy()
	return &UIImage{
		X:       x,
		Y:       y,
		Width:   width,
		Height:  height,
		Image:   image,
		visible: true,
		zindex:  0,
	}
}

func (u *UIImage) Contains(x int, y int) bool {
	return false
}

// Draw implements Element.
func (u *UIImage) Draw(*ebiten.Image) {
}

// HandleMouseDown implements Element.
func (u *UIImage) HandleMouseDown(x int, y int) bool {
	return false
}

// HandleMouseUp implements Element.
func (u *UIImage) HandleMouseUp(x int, y int) bool {
	return false
}

// Update implements Element.
func (u *UIImage) Update() {
}

func (u *UIImage) IsStatic() bool    { return true }
func (u *UIImage) GetZIndex() int    { return u.zindex }
func (u *UIImage) IsVisible() bool   { return u.visible }
func (u *UIImage) SetVisible(v bool) { u.visible = v }
func (u *UIImage) SetZIndex(idx int) { u.zindex = idx }
func (u *UIImage) GetTags() Tag      { return u.tags }
func (u *UIImage) SetTags(t Tag)     { u.tags = t }
