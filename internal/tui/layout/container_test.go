package layout

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type testModel struct{}

func (testModel) Init() tea.Cmd                       { return nil }
func (testModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return testModel{}, nil }
func (testModel) View() string                        { return "test" }
func (testModel) SetSize(int, int) tea.Cmd            { return nil }
func (testModel) GetSize() (int, int)                 { return 10, 10 }

func TestNewContainer(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := NewContainer(testModel{})
	if c == nil {
		t.Fatal("expected non-nil container")
	}
	if _, h := c.GetSize(); h != 0 {
		t.Errorf("expected default height 0, got %d", h)
	}
}

func TestContainerWithPadding(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := NewContainer(testModel{}, WithPadding(1, 2, 3, 4))
	c.SetSize(20, 20)
	view := c.View()
	if !strings.Contains(view, "test") {
		t.Error("expected view to contain model content")
	}
}

func TestContainerWithPaddingAll(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := NewContainer(testModel{}, WithPaddingAll(2))
	c.SetSize(20, 20)
	view := c.View()
	if !strings.Contains(view, "test") {
		t.Error("expected view to contain model content")
	}
}

func TestContainerWithBorder(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := NewContainer(testModel{}, WithBorderAll())
	c.SetSize(20, 20)
	view := c.View()
	if !strings.Contains(view, "test") {
		t.Error("expected view to contain model content")
	}
}

func TestContainerWithRoundedBorder(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := NewContainer(testModel{}, WithRoundedBorder())
	c.SetSize(20, 20)
	view := c.View()
	if !strings.Contains(view, "test") {
		t.Error("expected view to contain model content")
	}
}

func TestContainerWithBorderStyle(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := NewContainer(testModel{}, WithBorderStyle(lipgloss.DoubleBorder()))
	c.SetSize(20, 20)
	view := c.View()
	if !strings.Contains(view, "test") {
		t.Error("expected view to contain model content")
	}
}

func TestContainerGetSize(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := NewContainer(testModel{})
	c.SetSize(50, 30)
	w, h := c.GetSize()
	if w != 50 || h != 30 {
		t.Errorf("expected (50, 30), got (%d, %d)", w, h)
	}
}

func TestContainerSetSizeWithContent(t *testing.T) {
	t.Parallel()
	setupTheme()

	inner := NewContainer(testModel{})
	c := NewContainer(inner, WithPadding(1, 1, 1, 1))
	c.SetSize(20, 20)
	// The outer container sets its own size, which triggers SetSize on the inner container
	// adjusted for padding
}

func TestContainerBindingKeys(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := NewContainer(testModel{})
	keys := c.BindingKeys()
	if keys == nil {
		t.Error("expected non-nil bindings")
	}
}

func TestContainerViewWithBorderAndPadding(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := NewContainer(testModel{}, WithBorderAll(), WithPadding(1, 1, 1, 1))
	c.SetSize(30, 30)
	view := c.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestContainerFocusBlur(t *testing.T) {
	t.Parallel()

	c := NewContainer(testModel{})
	_, ok := c.(interface{ Focus() tea.Cmd })
	if ok {
		t.Log("container implements Focus")
	}
	_, ok = c.(interface{ Blur() tea.Cmd })
	if ok {
		t.Log("container implements Blur")
	}
}
