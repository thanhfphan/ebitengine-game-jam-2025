package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thanhfphan/ebitengj2025/internal/view"
	"golang.org/x/image/font"
)

var _ Element = (*UIBotHand)(nil)

type UIBotHand struct {
	X, Y          int
	Width, Height int
	CardCount     int
	CardUI        *UIImage

	visible bool
	zIndex  int
	font    font.Face
}

func NewUIBotHand(x, y, width, height int, font font.Face) *UIBotHand {
	return &UIBotHand{
		X:       x,
		Y:       y,
		Width:   width,
		Height:  height,
		visible: true,
		zIndex:  0,
		font:    font,
	}
}

func (h *UIBotHand) Update() {
	// No dynamic updates needed for bot hands
}

func (h *UIBotHand) Draw(screen *ebiten.Image) {
	if !h.visible || h.CardUI == nil || h.CardCount <= 0 {
		return
	}

	op := &ebiten.DrawImageOptions{}
	cardImg := h.CardUI.Image

	scaleX := float64(h.Width) / float64(cardImg.Bounds().Dx())
	scaleY := float64(h.Height) / float64(cardImg.Bounds().Dy())

	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(h.X), float64(h.Y))

	screen.DrawImage(cardImg, op)

	borderColor := color.RGBA{255, 255, 255, 255} // White border
	vector.StrokeRect(screen, float32(h.X), float32(h.Y), float32(h.Width), float32(h.Height), 1, borderColor, false)

	if h.CardCount > 0 {
		countText := fmt.Sprintf("%d", h.CardCount)
		face := h.font
		textWidth := font.MeasureString(face, countText).Ceil()
		metrics := face.Metrics()
		textHeight := (metrics.Ascent + metrics.Descent).Ceil()

		// Center the text within the card
		textX := h.X + (h.Width-textWidth)/2
		textY := h.Y + (h.Height-textHeight)/2 + metrics.Ascent.Ceil() // shift down by ascent

		text.Draw(screen, countText, h.font, textX, textY, color.White)
	}
}

func (h *UIBotHand) UpdateCards(cards []view.Card, cardBackImage *ebiten.Image) {
	h.CardCount = len(cards)

	// Create or update the UIImage for the card back
	if h.CardUI == nil {
		h.CardUI = NewUIImage(h.X, h.Y, cardBackImage)
	}
}

func (h *UIBotHand) Contains(x, y int) bool {
	// Bot hands don't need to be interactive
	return false
}

func (h *UIBotHand) HandleMouseDown(x, y int) bool {
	// Bot hands don't need to be interactive
	return false
}

func (h *UIBotHand) HandleMouseUp(x, y int) bool {
	// Bot hands don't need to be interactive
	return false
}

func (h *UIBotHand) IsVisible() bool             { return h.visible }
func (h *UIBotHand) SetVisible(v bool)           { h.visible = v }
func (h *UIBotHand) GetZIndex() int              { return h.zIndex }
func (h *UIBotHand) SetZIndex(z int)             { h.zIndex = z }
func (h *UIBotHand) IsStatic() bool              { return true }
func (h *UIBotHand) SetDraggable(draggable bool) {}
func (h *UIBotHand) SetPosition(x, y int) {
	h.X = x
	h.Y = y
}
