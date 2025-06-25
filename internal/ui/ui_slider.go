package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var _ Element = (*UISlider)(nil)

type UISlider struct {
	X, Y          int
	Width, Height int
	Value         float64
	OnChange      func(value float64)

	isDragging bool
	visible    bool
	zIndex     int
	tags       Tag
}

func NewUISlider(x, y, width, height int, initialValue float64) *UISlider {
	return &UISlider{
		X:          x,
		Y:          y,
		Width:      width,
		Height:     height,
		Value:      initialValue,
		visible:    true,
		isDragging: false,
	}
}

func (s *UISlider) Update() {
	if !s.visible {
		return
	}

	if s.isDragging {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			mx, _ := ebiten.CursorPosition()
			// Calculate new value based on mouse position
			relativeX := float64(mx - (s.X - s.Width/2))
			newValue := relativeX / float64(s.Width)

			// Clamp value between 0 and 1
			if newValue < 0 {
				newValue = 0
			} else if newValue > 1 {
				newValue = 1
			}

			if newValue != s.Value {
				s.Value = newValue
				if s.OnChange != nil {
					s.OnChange(s.Value)
				}
			}
		} else {
			s.isDragging = false
		}
	}
}

func (s *UISlider) Draw(screen *ebiten.Image) {
	if !s.visible {
		return
	}

	// Draw slider track
	trackColor := color.RGBA{100, 100, 100, 255}
	vector.DrawFilledRect(screen, float32(s.X-s.Width/2), float32(s.Y-s.Height/2),
		float32(s.Width), float32(s.Height), trackColor, false)

	// Draw slider fill
	fillWidth := float32(s.Width) * float32(s.Value)
	fillColor := color.RGBA{50, 150, 250, 255}
	vector.DrawFilledRect(screen, float32(s.X-s.Width/2), float32(s.Y-s.Height/2),
		fillWidth, float32(s.Height), fillColor, false)

	// Draw slider handle
	handleX := float32(s.X-s.Width/2) + fillWidth
	handleColor := color.RGBA{200, 200, 200, 255}
	vector.DrawFilledCircle(screen, handleX, float32(s.Y), float32(s.Height), handleColor, false)
}

func (s *UISlider) HandleMouseDown(x, y int) bool {
	if !s.visible {
		return false
	}

	// Check if click is within slider bounds (with some extra height for easier interaction)
	if x >= s.X-s.Width/2 && x <= s.X+s.Width/2 &&
		y >= s.Y-s.Height && y <= s.Y+s.Height {
		s.isDragging = true

		// Update value immediately on click
		relativeX := float64(x - (s.X - s.Width/2))
		newValue := relativeX / float64(s.Width)

		// Clamp value between 0 and 1
		if newValue < 0 {
			newValue = 0
		} else if newValue > 1 {
			newValue = 1
		}

		if newValue != s.Value {
			s.Value = newValue
			if s.OnChange != nil {
				s.OnChange(s.Value)
			}
		}

		return true
	}

	return false
}

func (s *UISlider) HandleMouseUp(x, y int) bool {
	if s.isDragging {
		s.isDragging = false
		return true
	}
	return false
}

func (s *UISlider) IsVisible() bool {
	return s.visible
}

func (s *UISlider) SetVisible(visible bool) {
	s.visible = visible
}

func (s *UISlider) ZIndex() int {
	return s.zIndex
}

func (s *UISlider) SetZIndex(z int) {
	s.zIndex = z
}

func (s *UISlider) GetTags() Tag {
	return s.tags
}

func (s *UISlider) SetTags(tags Tag) {
	s.tags = tags
}

// Additional methods to implement the Element interface
func (s *UISlider) Contains(x, y int) bool {
	if !s.visible {
		return false
	}

	// Check if point is within slider bounds (with some extra height for easier interaction)
	return x >= s.X-s.Width/2 && x <= s.X+s.Width/2 &&
		y >= s.Y-s.Height && y <= s.Y+s.Height
}

func (s *UISlider) SetPosition(x, y int) {
	s.X = x
	s.Y = y
}

func (s *UISlider) GetZIndex() int {
	return s.zIndex
}

func (s *UISlider) IsStatic() bool {
	return false
}

func (s *UISlider) SetDraggable(bool) {
	// Sliders aren't draggable as a whole, only the handle is interactive
}
