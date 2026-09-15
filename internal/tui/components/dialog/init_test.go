package dialog

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInitDialogView(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewInitDialogCmp()
	view := d.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
	if !containsStr(view, "Initialize Project") {
		t.Error("expected view to contain title")
	}
}

func TestInitDialogInit(t *testing.T) {
	t.Parallel()

	d := NewInitDialogCmp()
	cmd := d.Init()
	if cmd != nil {
		t.Error("expected nil init cmd")
	}
}

func TestInitDialogSelectedDefaults(t *testing.T) {
	t.Parallel()

	d := NewInitDialogCmp()
	if d.selected != 0 {
		t.Errorf("expected default selected 0, got %d", d.selected)
	}
}

func TestInitDialogUpdateKeyMsg(t *testing.T) {
	t.Parallel()

	d := NewInitDialogCmp()
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	model, cmd := d.Update(msg)
	if model == nil {
		t.Error("expected non-nil model")
	}
	if cmd == nil {
		t.Log("cmd may be nil")
	}
}

func TestInitDialogUpdateWindowSize(t *testing.T) {
	t.Parallel()

	d := NewInitDialogCmp()
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	model, cmd := d.Update(msg)
	if model == nil {
		t.Error("expected non-nil model")
	}
	if cmd != nil {
		t.Log("cmd may be nil")
	}
}

func TestInitDialogUpdateGenericMsg(t *testing.T) {
	t.Parallel()

	d := NewInitDialogCmp()
	model, cmd := d.Update(tea.Msg("generic"))
	if model == nil {
		t.Error("expected non-nil model")
	}
	if cmd != nil {
		t.Error("expected nil cmd for generic msg")
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
