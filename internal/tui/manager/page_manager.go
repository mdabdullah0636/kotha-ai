package manager

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"kotha/internal/tui/page"
	"kotha/internal/tui/layout"
)

type PageManager struct {
	currentPage page.PageID
	loadedPages map[page.PageID]bool
	pages       map[page.PageID]tea.Model
}

func NewPageManager(initialPage page.PageID) *PageManager {
	return &PageManager{
		currentPage: initialPage,
		loadedPages: make(map[page.PageID]bool),
		pages:       make(map[page.PageID]tea.Model),
	}
}

func (m *PageManager) RegisterPage(id page.PageID, model tea.Model) {
	m.pages[id] = model
}

func (m *PageManager) CurrentPage() page.PageID {
	return m.currentPage
}

func (m *PageManager) CurrentModel() tea.Model {
	return m.pages[m.currentPage]
}

func (m *PageManager) Pages() map[page.PageID]tea.Model {
	return m.pages
}

func (m *PageManager) LoadPage(id page.PageID) tea.Cmd {
	if _, ok := m.loadedPages[id]; !ok {
		if model, ok := m.pages[id]; ok {
			return model.Init()
		}
		m.loadedPages[id] = true
	}
	return nil
}

func (m *PageManager) IsLoaded(id page.PageID) bool {
	return m.loadedPages[id]
}

func (m *PageManager) MoveTo(id page.PageID, isAgentBusy func() bool) (tea.Cmd, error) {
	if isAgentBusy != nil && isAgentBusy() {
		return nil, errAgentBusy
	}
	m.currentPage = id
	return m.LoadPage(id), nil
}

func (m *PageManager) SetPageSize(id page.PageID, width, height int) tea.Cmd {
	model := m.pages[id]
	if sizable, ok := model.(layout.Sizeable); ok {
		return sizable.SetSize(width, height)
	}
	return nil
}

func (m *PageManager) SizeablePages() []page.PageID {
	var result []page.PageID
	for id, model := range m.pages {
		if _, ok := model.(layout.Sizeable); ok {
			result = append(result, id)
		}
	}
	return result
}

func (m *PageManager) SizeableBindings(id page.PageID) []key.Binding {
	model := m.pages[id]
	if b, ok := model.(layout.Bindings); ok {
		return b.BindingKeys()
	}
	return nil
}

func (m *PageManager) BindingsForPage(id page.PageID) []key.Binding {
	model := m.pages[id]
	if b, ok := model.(layout.Bindings); ok {
		return b.BindingKeys()
	}
	return nil
}

var errAgentBusy = context.Canceled
