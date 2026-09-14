package dialog

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"kotha/internal/tui/styles"
	"kotha/internal/tui/theme"
)

type ConfirmDialog struct {
	Title       string
	Description string
	Confirmed   bool
	Width       int
}

func NewConfirmDialog(title, description string) *ConfirmDialog {
	return &ConfirmDialog{
		Title:       title,
		Description: description,
		Width:       50,
	}
}

func (c *ConfirmDialog) View() string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()

	title := baseStyle.
		Bold(true).
		Foreground(t.Primary()).
		Render(c.Title)

	desc := baseStyle.
		Foreground(t.Text()).
		Width(c.Width - 4).
		Render(c.Description)

	content := lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		desc,
		"",
		baseStyle.Render("[Y] Yes  [N] No"),
	)

	return baseStyle.Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderBackground(t.Background()).
		BorderForeground(t.TextMuted()).
		Width(c.Width).
		Render(content)
}

func (c *ConfirmDialog) Init() tea.Cmd { return nil }
func (c *ConfirmDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}
func (c *ConfirmDialog) Bindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "confirm")),
		key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "cancel")),
	}
}
func (c *ConfirmDialog) SetSize(width, height int) tea.Cmd { c.Width = width; return nil }
func (c *ConfirmDialog) GetSize() (int, int) { return c.Width, 0 }

func (c *ConfirmDialog) BindingKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "confirm")),
		key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "cancel")),
	}
}
