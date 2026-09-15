package dialog

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestQuitDialogView(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewQuitCmp()
	view := d.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
	if !containsStr(view, "Are you sure you want to quit?") {
		t.Error("expected view to contain question")
	}
}

func TestQuitDialogInit(t *testing.T) {
	t.Parallel()

	d := NewQuitCmp()
	cmd := d.Init()
	if cmd != nil {
		t.Error("expected nil init cmd")
	}
}

func TestQuitDialogDefaultSelectedNo(t *testing.T) {
	t.Parallel()

	d := NewQuitCmp()
	cmp := d.(*quitDialogCmp)
	if !cmp.selectedNo {
		t.Error("expected selectedNo to default to true")
	}
}

func TestQuitDialogUpdateToggle(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewQuitCmp()
	cmp := d.(*quitDialogCmp)
	initial := cmp.selectedNo

	msg := tea.KeyMsg{Type: tea.KeyLeft}
	_, _ = d.Update(msg)

	if cmp.selectedNo == initial {
		t.Error("expected selectedNo to toggle")
	}
}

func TestQuitDialogUpdateYes(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewQuitCmp()
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
	_, cmd := d.Update(msg)
	if cmd == nil {
		t.Error("expected tea.Quit cmd")
	}
}

func TestQuitDialogUpdateNo(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewQuitCmp()
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	_, cmd := d.Update(msg)
	if cmd == nil {
		t.Error("expected CloseQuitMsg cmd")
	}
}

func TestQuitDialogUpdateEscape(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewQuitCmp()
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	model, cmd := d.Update(msg)
	if cmd != nil {
		t.Error("expected nil cmd for Escape (not bound in quit dialog)")
	}
	if model == nil {
		t.Error("expected non-nil model")
	}
}

func TestQuitDialogUpdateEnterOnYes(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewQuitCmp()
	cmp := d.(*quitDialogCmp)
	cmp.selectedNo = false

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := d.Update(msg)
	if cmd == nil {
		t.Error("expected tea.Quit when yes is selected")
	}
}

func TestQuitDialogUpdateEnterOnNo(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewQuitCmp()
	cmp := d.(*quitDialogCmp)
	cmp.selectedNo = true

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := d.Update(msg)
	if cmd == nil {
		t.Error("expected CloseQuitMsg when no is selected")
	}
}

func TestQuitDialogBindingKeys(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewQuitCmp()
	keys := d.BindingKeys()
	if len(keys) == 0 {
		t.Error("expected non-empty bindings")
	}
}

func TestQuitDialogUpdateGenericMsg(t *testing.T) {
	t.Parallel()
	setupTheme()

	d := NewQuitCmp()
	model, cmd := d.Update(tea.Msg("generic"))
	if model == nil {
		t.Error("expected non-nil model")
	}
	if cmd != nil {
		t.Error("expected nil cmd for generic msg")
	}
}
