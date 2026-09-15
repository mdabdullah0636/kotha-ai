package dialog

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCommandDialogView(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewCommandDialogCmp()
	view := d.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
	if !containsStr(view, "Commands") {
		t.Error("expected view to contain title")
	}
}

func TestCommandDialogInit(t *testing.T) {
	t.Parallel()

	d := NewCommandDialogCmp()
	cmd := d.Init()
	if cmd == nil {
		t.Log("init may return nil for list")
	}
}

func TestCommandDialogSetCommands(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewCommandDialogCmp()
	cmds := []Command{
		{ID: "1", Title: "Command 1", Description: "First command"},
		{ID: "2", Title: "Command 2", Description: ""},
	}
	d.SetCommands(cmds)

	items := d.(*commandDialogCmp).listView.GetItems()
	if len(items) != 2 {
		t.Errorf("expected 2 commands, got %d", len(items))
	}
	if items[0].Title != "Command 1" {
		t.Errorf("expected 'Command 1', got %q", items[0].Title)
	}
}

func TestCommandDialogSetCommandsEmpty(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewCommandDialogCmp()
	d.SetCommands([]Command{})

	items := d.(*commandDialogCmp).listView.GetItems()
	if len(items) != 0 {
		t.Errorf("expected 0 commands, got %d", len(items))
	}
}

func TestCommandDialogUpdateEscape(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewCommandDialogCmp()
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := d.Update(msg)
	if cmd == nil {
		t.Error("expected CloseCommandDialogMsg cmd")
	}
}

func TestCommandDialogUpdateEnterEmpty(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewCommandDialogCmp()
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := d.Update(msg)
	if cmd != nil {
		t.Error("expected nil cmd when no items selected")
	}
}

func TestCommandDialogUpdateGenericMsg(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewCommandDialogCmp()
	_, cmd := d.Update(tea.Msg("generic"))
	if cmd != nil {
		t.Error("expected nil cmd for generic msg")
	}
}

func TestCommandDialogBindingKeys(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewCommandDialogCmp()
	keys := d.BindingKeys()
	if len(keys) == 0 {
		t.Error("expected non-empty bindings")
	}
}

func TestCommandRender(t *testing.T) {
	t.Parallel()
	setupTheme()

	cmd := Command{ID: "1", Title: "Test Command", Description: "A test"}
	result := cmd.Render(true, 40)
	if result == "" {
		t.Error("expected non-empty render")
	}

	result2 := cmd.Render(false, 40)
	if result2 == "" {
		t.Error("expected non-empty render for unselected")
	}
}

func TestCommandRenderNoDescription(t *testing.T) {
	t.Parallel()
	setupTheme()

	cmd := Command{ID: "1", Title: "Test", Description: ""}
	result := cmd.Render(true, 40)
	if result == "" {
		t.Error("expected non-empty render")
	}
}
