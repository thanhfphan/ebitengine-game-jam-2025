package ui

import "github.com/hajimehoshi/ebiten/v2"

// Tag represents UI element categories for visibility control
type Tag uint8

const (
	TagNone     Tag = 0
	TagMenu     Tag = 1 << iota // Main menu elements
	TagInGame                   // In-game elements
	TagSettings                 // Settings elements
)

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
	GetTags() Tag
	SetTags(Tag)
	SetDraggable(bool)
}
