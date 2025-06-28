package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

var _ Element = (*UIButton)(nil)

type UIButton struct {
	X, Y            int
	Width, Height   int
	Text            string
	Visible         bool
	ZIndex          int
	OnClick         func()
	IsHovered       bool
	IsPressed       bool
	BackgroundColor color.RGBA
	HoverColor      color.RGBA
	PressedColor    color.RGBA
	TextColor       color.RGBA
	Font            font.Face
	Image           *ebiten.Image
}

func NewUIButton(x, y, width, height int, text string, font font.Face) *UIButton {
	return &UIButton{
		X:               x,
		Y:               y,
		Width:           width,
		Height:          height,
		Text:            text,
		Visible:         true,
		ZIndex:          0,
		BackgroundColor: color.RGBA{0x60, 0x60, 0x60, 0xff},
		HoverColor:      color.RGBA{0x80, 0x80, 0x80, 0xff},
		PressedColor:    color.RGBA{0x40, 0x40, 0x40, 0xff},
		TextColor:       color.RGBA{0xff, 0xff, 0xff, 0xff},
		Font:            font,
		Image:           nil,
	}
}

func NewUIImageButton(x, y, width, height int, image *ebiten.Image) *UIButton {
	return &UIButton{
		X:               x,
		Y:               y,
		Width:           width,
		Height:          height,
		Text:            "",
		Visible:         true,
		ZIndex:          0,
		BackgroundColor: color.RGBA{0x60, 0x60, 0x60, 0xff},
		HoverColor:      color.RGBA{0x80, 0x80, 0x80, 0xff},
		PressedColor:    color.RGBA{0x40, 0x40, 0x40, 0xff},
		TextColor:       color.RGBA{0xff, 0xff, 0xff, 0xff},
		Font:            nil,
		Image:           image,
	}
}

func (b *UIButton) Update() {
	if !b.Visible {
		return
	}

	x, y := ebiten.CursorPosition()
	b.IsHovered = b.Contains(x, y)

	b.IsPressed = b.IsHovered && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

func (b *UIButton) Draw(screen *ebiten.Image) {
	if !b.Visible {
		return
	}

	bgColor := b.BackgroundColor
	if b.IsPressed {
		bgColor = b.PressedColor
	} else if b.IsHovered {
		bgColor = b.HoverColor
	}

	if b.Image != nil {
		op := &ebiten.DrawImageOptions{}

		imgW, imgH := b.Image.Bounds().Dx(), b.Image.Bounds().Dy()
		scaleX := float64(b.Width) / float64(imgW)
		scaleY := float64(b.Height) / float64(imgH)

		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}

		op.GeoM.Scale(scale, scale)

		imgWidth, imgHeight := float64(imgW)*scale, float64(imgH)*scale
		op.GeoM.Translate(
			float64(b.X)+(float64(b.Width)-imgWidth)/2,
			float64(b.Y)+(float64(b.Height)-imgHeight)/2,
		)

		screen.DrawImage(b.Image, op)
	} else if b.Font != nil && b.Text != "" {
		vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), bgColor, false)
		// Border
		vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), 1, color.RGBA{0x20, 0x20, 0x20, 0xff}, false)
		vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), 1, float32(b.Height), color.RGBA{0x20, 0x20, 0x20, 0xff}, false)
		vector.DrawFilledRect(screen, float32(b.X+b.Width-1), float32(b.Y), 1, float32(b.Height), color.RGBA{0x20, 0x20, 0x20, 0xff}, false)
		vector.DrawFilledRect(screen, float32(b.X), float32(b.Y+b.Height-1), float32(b.Width), 1, color.RGBA{0x20, 0x20, 0x20, 0xff}, false)

		bounds := text.BoundString(b.Font, b.Text)
		textWidth := bounds.Max.X - bounds.Min.X
		textHeight := bounds.Max.Y - bounds.Min.Y
		textX := b.X + (b.Width-textWidth)/2
		textY := b.Y + (b.Height+textHeight)/2

		text.Draw(screen, b.Text, b.Font, textX, textY, b.TextColor)
	}
}

func (b *UIButton) Contains(x, y int) bool {
	return b.Visible && x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height
}

func (b *UIButton) HandleMouseDown(x, y int) bool {
	return b.Contains(x, y)
}

func (b *UIButton) HandleMouseUp(x, y int) bool {
	if b.Contains(x, y) && b.IsPressed {
		if b.OnClick != nil {
			b.OnClick()
		}
		return true
	}
	return false
}

func (b *UIButton) IsVisible() bool             { return b.Visible }
func (b *UIButton) SetVisible(visible bool)     { b.Visible = visible }
func (b *UIButton) GetZIndex() int              { return b.ZIndex }
func (b *UIButton) SetZIndex(zIndex int)        { b.ZIndex = zIndex }
func (b *UIButton) IsStatic() bool              { return false }
func (b *UIButton) SetDraggable(draggable bool) {}
func (b *UIButton) SetPosition(x, y int) {
	b.X = x
	b.Y = y
}
