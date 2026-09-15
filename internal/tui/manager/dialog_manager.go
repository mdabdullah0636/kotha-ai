package manager

import (
	"errors"
	"slices"

	"github.com/charmbracelet/lipgloss"
	"kotha/internal/tui/layout"
)

var ErrDialogLimit = errors.New("dialog limit reached")

type DialogConfig struct {
	Name     string
	View     func() string
	SetSize  func(width, height int)
	Bindings func() []interface{}
}

type DialogManager struct {
	active        map[string]bool
	configs       map[string]*DialogConfig
	maxConcurrent int
}

func NewDialogManager(maxConcurrent int) *DialogManager {
	if maxConcurrent <= 0 {
		maxConcurrent = 3
	}
	return &DialogManager{
		active:        make(map[string]bool),
		configs:       make(map[string]*DialogConfig),
		maxConcurrent: maxConcurrent,
	}
}

func (m *DialogManager) Register(name string, config DialogConfig) {
	m.configs[name] = &config
}

func (m *DialogManager) Show(name string) error {
	if m.active[name] {
		return nil
	}
	if m.Count() >= m.maxConcurrent {
		return errDialogLimitReached
	}
	m.active[name] = true
	return nil
}

func (m *DialogManager) Hide(name string) {
	delete(m.active, name)
}

func (m *DialogManager) IsActive(name string) bool {
	return m.active[name]
}

func (m *DialogManager) Count() int {
	return len(m.active)
}

func (m *DialogManager) ActiveNames() []string {
	names := make([]string, 0, len(m.active))
	for name := range m.active {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func (m *DialogManager) HasAny() bool {
	return len(m.active) > 0
}

func (m *DialogManager) View(appView string) string {
	for _, name := range m.ActiveNames() {
		config, ok := m.configs[name]
		if !ok {
			continue
		}
		overlay := config.View()
		if overlay == "" {
			continue
		}
		row := lipgloss.Height(appView) / 2
		row -= lipgloss.Height(overlay) / 2
		col := lipgloss.Width(appView) / 2
		col -= lipgloss.Width(overlay) / 2
		appView = layout.PlaceOverlay(col, row, overlay, appView, true)
	}
	return appView
}

func (m *DialogManager) SetSizeAll(width, height int) {
	for name := range m.active {
		if config, ok := m.configs[name]; ok && config.SetSize != nil {
			config.SetSize(width, height)
		}
	}
}

func (m *DialogManager) BindingsForActive() []interface{} {
	var result []interface{}
	for name := range m.active {
		if config, ok := m.configs[name]; ok && config.Bindings != nil {
			result = append(result, config.Bindings()...)
		}
	}
	return result
}

var errDialogLimitReached = ErrDialogLimit
