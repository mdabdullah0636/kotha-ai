package dialog

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestThemeDialogView(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewThemeDialogCmp()
	view := d.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
}

func TestThemeDialogInit(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewThemeDialogCmp()
	cmd := d.Init()
	if cmd != nil {
		t.Error("expected nil init cmd")
	}
}

func TestThemeDialogInitThemes(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewThemeDialogCmp()
	d.Init()
	cmp := d.(*themeDialogCmp)
	if len(cmp.themes) == 0 {
		t.Error("expected at least one theme")
	}
}

func TestThemeDialogUpdateUp(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewThemeDialogCmp()
	d.Init()
	cmp := d.(*themeDialogCmp)
	initial := cmp.selectedIdx

	msg := tea.KeyMsg{Type: tea.KeyUp}
	_, _ = d.Update(msg)

	if cmp.selectedIdx != initial && cmp.selectedIdx != initial-1 {
		t.Errorf("expected selectedIdx to change from %d, got %d", initial, cmp.selectedIdx)
	}
}

func TestThemeDialogUpdateDown(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewThemeDialogCmp()
	d.Init()
	cmp := d.(*themeDialogCmp)
	maxIdx := len(cmp.themes) - 1
	if maxIdx < 0 {
		maxIdx = 0
	}
	initial := cmp.selectedIdx

	msg := tea.KeyMsg{Type: tea.KeyDown}
	_, _ = d.Update(msg)

	expected := initial + 1
	if expected > maxIdx {
		expected = maxIdx
	}
	if cmp.selectedIdx != expected {
		t.Errorf("expected selectedIdx %d, got %d", expected, cmp.selectedIdx)
	}
}

func TestThemeDialogUpdateEscape(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewThemeDialogCmp()
	d.Init()
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := d.Update(msg)
	if cmd == nil {
		t.Error("expected CloseThemeDialogMsg cmd")
	}
}

func TestThemeDialogUpdateWindowSize(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewThemeDialogCmp()
	d.Init()
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	_, _ = d.Update(msg)
	cmp := d.(*themeDialogCmp)
	if cmp.width != 120 || cmp.height != 40 {
		t.Errorf("expected (120, 40), got (%d, %d)", cmp.width, cmp.height)
	}
}

func TestThemeDialogUpdateGenericMsg(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewThemeDialogCmp()
	d.Init()
	_, cmd := d.Update(tea.Msg("generic"))
	if cmd != nil {
		t.Error("expected nil cmd for generic msg")
	}
}

func TestThemeDialogBindingKeys(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewThemeDialogCmp()
	keys := d.BindingKeys()
	if len(keys) == 0 {
		t.Error("expected non-empty bindings")
	}
}

func TestThemeDialogViewNoThemes(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewThemeDialogCmp()
	cmp := d.(*themeDialogCmp)
	cmp.themes = []string{}
	view := d.View()
	if view == "" {
		t.Fatal("expected non-empty view even with no themes")
	}
}
