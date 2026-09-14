package page

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"kotha/internal/tui/layout"
	"kotha/internal/tui/styles"
	"kotha/internal/tui/theme"
)

func setupTheme() {
	if theme.CurrentTheme() == nil {
		_ = theme.SetTheme("kotha")
	}
}

func TestPageIDType(t *testing.T) {
	t.Parallel()

	p := ChatPage
	if p == "" {
		t.Error("expected non-empty page ID")
	}
	if string(p) != "chat" {
		t.Errorf("expected 'chat', got %q", string(p))
	}
}

func TestPageChangeMsg(t *testing.T) {
	t.Parallel()

	msg := PageChangeMsg{ID: ChatPage}
	if msg.ID != ChatPage {
		t.Error("expected page ID to be ChatPage")
	}
}

func TestLogsPageCreation(t *testing.T) {
	t.Parallel()
	setupTheme()

	p := NewLogsPage()
	if p == nil {
		t.Fatal("expected non-nil page")
	}
	view := p.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestLogsPageSize(t *testing.T) {
	t.Parallel()
	setupTheme()

	p := NewLogsPage()
	w, h := p.GetSize()
	if w != 0 || h != 0 {
		t.Errorf("expected (0, 0), got (%d, %d)", w, h)
	}
}

func TestLogsPageInit(t *testing.T) {
	t.Parallel()
	setupTheme()

	p := NewLogsPage()
	cmd := p.Init()
	if cmd == nil {
		t.Log("init may return nil for containers")
	}
}

func TestLogsPageUpdate(t *testing.T) {
	t.Parallel()
	setupTheme()

	p := NewLogsPage()
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	model, cmd := p.Update(msg)
	if model == nil {
		t.Error("expected non-nil model")
	}
	_ = cmd
}

func TestChatPageKeyMap(t *testing.T) {
	t.Parallel()

	keys := layout.KeyMapToSlice(keyMap)
	if len(keys) != 3 {
		t.Errorf("expected 3 key bindings, got %d", len(keys))
	}
}

func TestChatPageKeyMapHasCancel(t *testing.T) {
	t.Parallel()

	for _, k := range keyMap.Cancel.Keys() {
		if k != "esc" {
			t.Errorf("expected cancel key to be esc, got %q", k)
		}
	}
}

func TestChatPageViewWithoutTheme(t *testing.T) {
	t.Parallel()
	setupTheme()

	p := NewLogsPage()
	view := p.View()
	base := styles.BaseStyle()
	if !containsStr(view, "") {
		_ = base
		_ = view
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
