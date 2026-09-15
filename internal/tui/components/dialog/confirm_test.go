package dialog

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestConfirmDialogView(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewConfirmDialog("Are you sure?", "This action cannot be undone.")
	view := d.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
	if !contains(view, "Are you sure?") {
		t.Error("expected view to contain title")
	}
}

func TestConfirmDialogInit(t *testing.T) {
	t.Parallel()

	d := NewConfirmDialog("Title", "Desc")
	cmd := d.Init()
	if cmd != nil {
		t.Error("expected nil init cmd")
	}
}

func TestConfirmDialogUpdate(t *testing.T) {
	t.Parallel()

	d := NewConfirmDialog("Title", "Desc")
	_, cmd := d.Update(tea.KeyMsg{})
	if cmd != nil {
		t.Error("expected nil cmd for generic msg")
	}
}

func TestConfirmDialogBindings(t *testing.T) {
	t.Parallel()

	d := NewConfirmDialog("Title", "Desc")
	keys := d.BindingKeys()
	if len(keys) != 2 {
		t.Errorf("expected 2 bindings, got %d", len(keys))
	}
}

func TestConfirmDialogSetSize(t *testing.T) {
	t.Parallel()

	d := NewConfirmDialog("Title", "Desc")
	d.SetSize(80, 24)
	if d.Width != 80 {
		t.Errorf("expected width 80, got %d", d.Width)
	}
	w, h := d.GetSize()
	if w != 80 || h != 0 {
		t.Errorf("expected (80, 0), got (%d, %d)", w, h)
	}
}

func TestConfirmDialogDefaultWidth(t *testing.T) {
	t.Parallel()

	d := NewConfirmDialog("Title", "Desc")
	if d.Width != 50 {
		t.Errorf("expected default width 50, got %d", d.Width)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
