package manager

import (
	"strings"
)

type OverlayEntry struct {
	Visible bool
	View    func() string
}

type OverlayManager struct {
	overlays map[string]*OverlayEntry
	order    []string
}

func NewOverlayManager() *OverlayManager {
	return &OverlayManager{
		overlays: make(map[string]*OverlayEntry),
	}
}

func (m *OverlayManager) Set(name string, visible bool, view func() string) {
	if _, exists := m.overlays[name]; !exists {
		m.order = append(m.order, name)
	}
	m.overlays[name] = &OverlayEntry{Visible: visible, View: view}
}

func (m *OverlayManager) Has(name string) bool {
	entry, ok := m.overlays[name]
	return ok && entry.Visible
}

func (m *OverlayManager) Hide(name string) {
	if entry, ok := m.overlays[name]; ok {
		entry.Visible = false
	}
}

func (m *OverlayManager) Show(name string) {
	if entry, ok := m.overlays[name]; ok {
		entry.Visible = true
	}
}

func (m *OverlayManager) Remove(name string) {
	delete(m.overlays, name)
	for i, n := range m.order {
		if n == name {
			m.order = append(m.order[:i], m.order[i+1:]...)
			break
		}
	}
}

func (m *OverlayManager) Clear() {
	m.overlays = make(map[string]*OverlayEntry)
	m.order = nil
}

func (m *OverlayManager) View(appView string) string {
	var overlays []string
	for _, name := range m.order {
		if entry, ok := m.overlays[name]; ok && entry.Visible {
			if content := entry.View(); content != "" {
				overlays = append(overlays, content)
			}
		}
	}
	if len(overlays) == 0 {
		return appView
	}
	return strings.Join(overlays, "") + appView
}

func (m *OverlayManager) VisibleNames() []string {
	names := make([]string, 0)
	for _, name := range m.order {
		if entry, ok := m.overlays[name]; ok && entry.Visible {
			names = append(names, name)
		}
	}
	return names
}
