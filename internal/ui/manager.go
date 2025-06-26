package ui

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type Manager struct {
	elements      []Element
	topzindex     int
	activeElement Element
	needSort      bool
}

func NewManager() *Manager {
	return &Manager{
		elements:      make([]Element, 0),
		topzindex:     0,
		activeElement: nil,
		needSort:      false,
	}
}

func (m *Manager) Update() {
	for _, e := range m.elements {
		if e.IsVisible() {
			e.Update()
		}
	}
}

func (m *Manager) Draw(screen *ebiten.Image) {
	if m.needSort {
		sort.Slice(m.elements, func(i, j int) bool {
			return m.elements[i].GetZIndex() < m.elements[j].GetZIndex()
		})
		m.needSort = false
	}

	for _, e := range m.elements {
		if e.IsVisible() {
			e.Draw(screen)
		}
	}
}

func (m *Manager) AddElement(element Element) {
	m.elements = append(m.elements, element)
	m.BringToFront(element)
}

func (m *Manager) RemoveElement(element Element) {
	for i, e := range m.elements {
		if e == element {
			m.elements = append(m.elements[:i], m.elements[i+1:]...)
			m.activeElement = nil
			break
		}
	}
}

func (m *Manager) HandleMouseDown(x, y int) bool {
	for i := len(m.elements) - 1; i >= 0; i-- {
		element := m.elements[i]
		if element.IsVisible() && element.Contains(x, y) {
			if element.HandleMouseDown(x, y) {
				if !element.IsStatic() {
					m.BringToFront(element)
				}
				m.activeElement = element
				return true
			}
		}
	}
	m.activeElement = nil
	return false
}

func (m *Manager) HandleMouseUp(x, y int) bool {
	if m.activeElement != nil {
		return m.activeElement.HandleMouseUp(x, y)
	}

	for i := len(m.elements) - 1; i >= 0; i-- {
		element := m.elements[i]
		if element.IsVisible() && element.HandleMouseUp(x, y) {
			return true
		}
	}
	return false
}

func (m *Manager) BringToFront(element Element) {
	m.topzindex++
	element.SetZIndex(m.topzindex)
	m.needSort = true
}

func (m *Manager) Clear() {
	m.elements = make([]Element, 0)
	m.topzindex = 0
	m.activeElement = nil
	m.needSort = false
}
