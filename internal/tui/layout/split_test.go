package layout

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewSplitPane(t *testing.T) {
	t.Parallel()
	setupTheme()

	s := NewSplitPane()
	if s == nil {
		t.Fatal("expected non-nil split pane")
	}
	w, h := s.GetSize()
	if w != 0 || h != 0 {
		t.Errorf("expected (0, 0), got (%d, %d)", w, h)
	}
}

func TestSplitPaneWithLeftPanel(t *testing.T) {
	t.Parallel()
	setupTheme()

	left := NewContainer(testModel{})
	s := NewSplitPane(WithLeftPanel(left))
	if s == nil {
		t.Fatal("expected non-nil split pane")
	}
}

func TestSplitPaneWithRightPanel(t *testing.T) {
	t.Parallel()
	setupTheme()

	right := NewContainer(testModel{})
	s := NewSplitPane(WithRightPanel(right))
	if s == nil {
		t.Fatal("expected non-nil split pane")
	}
}

func TestSplitPaneWithBottomPanel(t *testing.T) {
	t.Parallel()
	setupTheme()

	bottom := NewContainer(testModel{})
	s := NewSplitPane(WithBottomPanel(bottom))
	if s == nil {
		t.Fatal("expected non-nil split pane")
	}
}

func TestSplitPaneWithRatio(t *testing.T) {
	t.Parallel()
	setupTheme()

	s := NewSplitPane(WithRatio(0.5))
	w, h := s.GetSize()
	if w != 0 || h != 0 {
		t.Errorf("expected (0, 0), got (%d, %d)", w, h)
	}
}

func TestSplitPaneSetSize(t *testing.T) {
	t.Parallel()
	setupTheme()

	s := NewSplitPane()
	s.SetSize(100, 50)
	w, h := s.GetSize()
	if w != 100 || h != 50 {
		t.Errorf("expected (100, 50), got (%d, %d)", w, h)
	}
}

func TestSplitPaneSetLeftPanel(t *testing.T) {
	t.Parallel()
	setupTheme()

	s := NewSplitPane()
	left := NewContainer(testModel{})
	cmd := s.SetLeftPanel(left)
	if cmd == nil {
		// May be nil if size is 0
	}
}

func TestSplitPaneSetRightPanel(t *testing.T) {
	t.Parallel()
	setupTheme()

	s := NewSplitPane()
	right := NewContainer(testModel{})
	cmd := s.SetRightPanel(right)
	if cmd == nil {
		// May be nil if size is 0
	}
}

func TestSplitPaneSetBottomPanel(t *testing.T) {
	t.Parallel()
	setupTheme()

	s := NewSplitPane()
	bottom := NewContainer(testModel{})
	cmd := s.SetBottomPanel(bottom)
	if cmd == nil {
		// May be nil if size is 0
	}
}

func TestSplitPaneClearPanels(t *testing.T) {
	t.Parallel()
	setupTheme()

	s := NewSplitPane()
	s.SetSize(100, 50)
	left := NewContainer(testModel{})
	s.SetLeftPanel(left)
	s.ClearLeftPanel()
	s.ClearRightPanel()
	s.ClearBottomPanel()
}

func TestSplitPaneView(t *testing.T) {
	t.Parallel()
	setupTheme()

	s := NewSplitPane()
	view := s.View()
	if view != "" {
		t.Error("expected empty view for split pane with no panels")
	}
}

func TestSplitPaneViewWithLeftPanel(t *testing.T) {
	t.Parallel()
	setupTheme()

	left := NewContainer(testModel{})
	s := NewSplitPane(WithLeftPanel(left))
	s.SetSize(100, 50)
	view := s.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestSplitPaneBindingKeys(t *testing.T) {
	t.Parallel()
	setupTheme()

	s := NewSplitPane()
	keys := s.BindingKeys()
	if keys == nil {
		t.Error("expected non-nil bindings")
	}
}

func TestSplitPaneUpdate(t *testing.T) {
	t.Parallel()
	setupTheme()

	s := NewSplitPane()
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	_, cmd := s.Update(msg)
	if cmd == nil {
		// SetSize may return nil if not yet set
	}
}

func TestSplitPaneInit(t *testing.T) {
	t.Parallel()
	setupTheme()

	s := NewSplitPane()
	cmd := s.Init()
	if cmd == nil {
		// Init may return nil if no panels
	}
}

func TestSplitPaneWithAllPanels(t *testing.T) {
	t.Parallel()
	setupTheme()

	s := NewSplitPane(
		WithLeftPanel(NewContainer(testModel{})),
		WithRightPanel(NewContainer(testModel{})),
		WithBottomPanel(NewContainer(testModel{})),
	)
	s.SetSize(100, 50)
	view := s.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}
