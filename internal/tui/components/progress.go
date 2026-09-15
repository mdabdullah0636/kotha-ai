package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"kotha/internal/tui/theme"
)

type ProgressBar struct {
	Width   int
	Percent int
	Message string
}

func NewProgressBar(width int, message string) *ProgressBar {
	return &ProgressBar{
		Width:   width,
		Message: message,
	}
}

func (p *ProgressBar) SetPercent(percent int) {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	p.Percent = percent
}

func (p *ProgressBar) View() string {
	t := theme.CurrentTheme()
	barWidth := p.Width - 2
	if barWidth < 0 {
		barWidth = 0
	}

	filled := strings.Repeat("█", p.Percent*barWidth/100)
	empty := strings.Repeat("░", barWidth-len(filled))

	barStyle := lipgloss.NewStyle().
		Background(t.Primary()).
		Foreground(t.Background())

	emptyStyle := lipgloss.NewStyle().
		Background(t.BackgroundDarker())

	bar := barStyle.Render(filled) + emptyStyle.Render(empty)

	label := ""
	if p.Message != "" {
		label = p.Message + " "
	}
	label += fmt.Sprintf("%d%%", p.Percent)

	labelStyle := lipgloss.NewStyle().
		Foreground(t.Primary()).
		Bold(true)

	content := bar + " " + labelStyle.Render(label)

	return lipgloss.NewStyle().
		Width(p.Width).
		Padding(0, 1).
		Render(content)
}
