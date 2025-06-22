package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

var _ Element = (*UILabel)(nil)

type UILabel struct {
	X, Y      int
	Text      string
	Visible   bool
	ZIndex    int
	Tags      Tag
	TextColor color.RGBA
	Font      font.Face
	align     Align

	// Hover
	HoverColor  color.RGBA
	HoverScale  float64
	IsHovered   bool
	EnableHover bool
}

// NewUILabel create a label with default alignment is AlignLeft.
func NewUILabel(x, y int, text string, f font.Face) *UILabel {
	return &UILabel{
		X:           x,
		Y:           y,
		Text:        text,
		Visible:     true,
		ZIndex:      0,
		Tags:        TagNone,
		TextColor:   color.RGBA{0xff, 0xff, 0xff, 0xff},
		Font:        f,
		align:       AlignLeft,
		HoverColor:  color.RGBA{0xff, 0xcc, 0x00, 0xff}, // Yellow
		HoverScale:  1.0,
		IsHovered:   false,
		EnableHover: false,
	}
}

func (l *UILabel) AlignCenter() { l.align = AlignCenter }
func (l *UILabel) AlignLeft()   { l.align = AlignLeft }
func (l *UILabel) AlignRight()  { l.align = AlignRight }

func (l *UILabel) Update() {
	if !l.Visible || !l.EnableHover || l.Font == nil {
		l.IsHovered = false
		return
	}
	x, y := ebiten.CursorPosition()
	l.IsHovered = l.Contains(x, y)
}

func (l *UILabel) Draw(screen *ebiten.Image) {
	if !l.Visible || l.Font == nil {
		return
	}

	bounds := text.BoundString(l.Font, l.Text)
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y
	x := l.X
	switch l.align {
	case AlignCenter:
		x = l.X - w/2
	case AlignRight:
		x = l.X - w
	}
	y := l.Y + h

	clr := l.TextColor
	if l.IsHovered && l.EnableHover {
		clr = l.HoverColor
	}

	op := &ebiten.DrawImageOptions{}
	scale := 1.0
	if l.IsHovered && l.EnableHover && l.HoverScale > 1.0 {
		scale = l.HoverScale
	}
	op.GeoM.Scale(scale, scale)

	tx := float64(x)
	ty := float64(y)
	if scale > 1.0 {
		tx -= float64(w) * (scale - 1.0) / 2
		ty -= float64(h) * (scale - 1.0) / 2
	}
	op.GeoM.Translate(tx, ty)
	op.ColorScale.Scale(float32(clr.R)/255, float32(clr.G)/255, float32(clr.B)/255, float32(clr.A)/255)

	text.DrawWithOptions(screen, l.Text, l.Font, op)
}

func (l *UILabel) Contains(x, y int) bool {
	if !l.Visible || l.Font == nil {
		return false
	}
	bounds := text.BoundString(l.Font, l.Text)
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y

	x0 := l.X
	switch l.align {
	case AlignCenter:
		x0 = l.X - w/2
	case AlignRight:
		x0 = l.X - w
	}
	y0 := l.Y
	return x >= x0 && x < x0+w && y >= y0 && y < y0+h
}

func (l *UILabel) HandleMouseDown(_, _ int) bool { return false }
func (l *UILabel) HandleMouseUp(_, _ int) bool   { return false }

func (l *UILabel) IsVisible() bool      { return l.Visible }
func (l *UILabel) SetVisible(v bool)    { l.Visible = v }
func (l *UILabel) GetZIndex() int       { return l.ZIndex }
func (l *UILabel) SetZIndex(z int)      { l.ZIndex = z }
func (l *UILabel) IsStatic() bool       { return true }
func (l *UILabel) GetTags() Tag         { return l.Tags }
func (l *UILabel) SetTags(t Tag)        { l.Tags = t }
func (l *UILabel) SetDraggable(_ bool)  {}
func (l *UILabel) SetPosition(x, y int) { l.X, l.Y = x, y }
