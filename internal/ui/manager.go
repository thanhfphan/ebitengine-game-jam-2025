package ui

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type Manager struct {
	Elements      []Element
	TopZIndex     int
	ActiveElement Element
}

func NewManager() *Manager {
	return &Manager{
		Elements:      make([]Element, 0),
		TopZIndex:     0,
		ActiveElement: nil,
	}
}

func (m *Manager) Update() {
	for _, e := range m.Elements {
		e.Update()
	}
}

func (m *Manager) Draw(screen *ebiten.Image) {
	sort.Slice(m.Elements, func(i, j int) bool {
		return m.Elements[i].GetZIndex() < m.Elements[j].GetZIndex()
	})

	for _, e := range m.Elements {
		if e.IsVisible() {
			e.Draw(screen)
		}
	}
}

func (m *Manager) AddElement(element Element) {
	m.Elements = append(m.Elements, element)
	m.BringToFront(element)
}

func (m *Manager) RemoveElement(element Element) {
	for i, e := range m.Elements {
		if e == element {
			m.Elements = append(m.Elements[:i], m.Elements[i+1:]...)
			break
		}
	}
}

func (m *Manager) HandleMouseDown(x, y int) bool {
	for i := len(m.Elements) - 1; i >= 0; i-- {
		element := m.Elements[i]
		if element.IsVisible() && element.Contains(x, y) {
			if element.HandleMouseDown(x, y) {
				if !element.IsStatic() {
					m.BringToFront(element)
				}
				m.ActiveElement = element
				return true
			}
		}
	}
	m.ActiveElement = nil
	return false
}

func (m *Manager) HandleMouseUp(x, y int) bool {
	if m.ActiveElement != nil {
		return m.ActiveElement.HandleMouseUp(x, y)
	}

	for i := len(m.Elements) - 1; i >= 0; i-- {
		element := m.Elements[i]
		if element.IsVisible() && element.HandleMouseUp(x, y) {
			return true
		}
	}
	return false
}

func (m *Manager) BringToFront(element Element) {
	m.TopZIndex++
	element.SetZIndex(m.TopZIndex)
}
