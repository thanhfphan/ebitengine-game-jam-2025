package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

var _ Element = (*UILabel)(nil)

type UILabel struct {
	X, Y       int
	Text       string
	Font       font.Face
	TextColor  color.Color
	HoverColor color.Color

	EnableHover bool
	HoverScale  float64
	isHovered   bool

	visible  bool
	zIndex   int
	tags     Tag
	centered bool
}

func NewUILabel(x, y int, text string, font font.Face) *UILabel {
	return &UILabel{
		X:          x,
		Y:          y,
		Text:       text,
		Font:       font,
		TextColor:  color.White,
		HoverColor: color.White,
		visible:    true,
		HoverScale: 1.0,
	}
}

func (l *UILabel) Update() {
	if !l.visible || !l.EnableHover {
		return
	}

	mx, my := ebiten.CursorPosition()
	bounds := text.BoundString(l.Font, l.Text)
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	x := l.X
	if l.centered {
		x = l.X - width/2
	}

	l.isHovered = mx >= x && mx <= x+width &&
		my >= l.Y-height && my <= l.Y
}

func (l *UILabel) Draw(screen *ebiten.Image) {
	if !l.visible {
		return
	}

	textColor := l.TextColor
	if l.isHovered && l.EnableHover {
		textColor = l.HoverColor
	}

	x := l.X
	if l.centered {
		bounds := text.BoundString(l.Font, l.Text)
		width := bounds.Max.X - bounds.Min.X
		x = l.X - width/2
	}

	text.Draw(screen, l.Text, l.Font, x, l.Y, textColor)
}

func (l *UILabel) AlignCenter() {
	l.centered = true
}

func (l *UILabel) HandleMouseDown(x, y int) bool {
	return false
}

func (l *UILabel) HandleMouseUp(x, y int) bool {
	return false
}

func (l *UILabel) IsVisible() bool {
	return l.visible
}

func (l *UILabel) SetVisible(visible bool) {
	l.visible = visible
}

func (l *UILabel) ZIndex() int {
	return l.zIndex
}

func (l *UILabel) SetZIndex(z int) {
	l.zIndex = z
}

func (l *UILabel) GetTags() Tag {
	return l.tags
}

func (l *UILabel) SetTags(tags Tag) {
	l.tags = tags
}

func (l *UILabel) Contains(x, y int) bool {
	if !l.visible {
		return false
	}

	bounds := text.BoundString(l.Font, l.Text)
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	labelX := l.X
	if l.centered {
		labelX = l.X - width/2
	}

	return x >= labelX && x <= labelX+width &&
		y >= l.Y-height && y <= l.Y
}

func (l *UILabel) SetPosition(x, y int) {
	l.X = x
	l.Y = y
}

func (l *UILabel) GetZIndex() int {
	return l.zIndex
}

func (l *UILabel) IsStatic() bool {
	return true
}

func (l *UILabel) SetDraggable(bool) {
	// Labels aren't draggable
}
