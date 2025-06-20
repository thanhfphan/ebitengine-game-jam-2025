package ui

import "github.com/hajimehoshi/ebiten/v2"

type Element interface {
	Draw(*ebiten.Image)
	Update()
	IsVisible() bool
	SetVisible(bool)
	Contains(x, y int) bool
	HandleMouseDown(x, y int) bool
	HandleMouseUp(x, y int) bool
	GetZIndex() int
	SetZIndex(int)
	IsStatic() bool
}
