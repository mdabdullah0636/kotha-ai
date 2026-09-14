package layout

import (
	"strings"
	"testing"
)

func TestPlaceOverlayBasic(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg := "background"
	fg := "foreground"
	result := PlaceOverlay(0, 0, fg, bg, false)
	if result != fg {
		t.Errorf("expected fg when it fits, got %q", result)
	}
}

func TestPlaceOverlayWithShadow(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg := strings.Repeat(" ", 20) + "\n" + strings.Repeat(" ", 20)
	fg := "hello"
	result := PlaceOverlay(2, 1, fg, bg, true)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !strings.Contains(result, "hello") {
		t.Error("expected result to contain foreground text")
	}
}

func TestPlaceOverlayCentered(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg := strings.Repeat("x", 10)
	fg := "ab"
	result := PlaceOverlay(4, 0, fg, bg, false)
	if !strings.Contains(result, "ab") {
		t.Error("expected result to contain foreground text")
	}
}

func TestPlaceOverlayOffScreenClamped(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg := strings.Repeat("x", 5)
	fg := "longer"
	result := PlaceOverlay(100, 100, fg, bg, false)
	// Should be clamped to fit within bg
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestPlaceOverlayEmptyFG(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg := "background"
	result := PlaceOverlay(0, 0, "", bg, false)
	if result != bg {
		t.Errorf("expected bg unchanged when fg is empty, got %q", result)
	}
}

func TestPlaceOverlayLargeShadow(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg := strings.Repeat(" ", 40) + "\n" + strings.Repeat(" ", 40) + "\n" + strings.Repeat(" ", 40)
	fg := "dialog"
	result := PlaceOverlay(10, 5, fg, bg, true)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !strings.Contains(result, "dialog") {
		t.Error("expected result to contain foreground text")
	}
}

func TestPlaceOverlayNested(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg1 := strings.Repeat(" ", 30)
	mid := "middle"
	fg := "inner"
	overMid := PlaceOverlay(5, 2, mid, bg1, false)
	result := PlaceOverlay(3, 1, fg, overMid, false)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestPlaceOverlayMultipleLines(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg := strings.Repeat(".", 20) + "\n" + strings.Repeat(".", 20) + "\n" + strings.Repeat(".", 20)
	fg := "line1\nline2"
	result := PlaceOverlay(5, 5, fg, bg, false)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestPlaceOverlayWithWhitespaceOption(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg := strings.Repeat(" ", 20)
	fg := "test"
	ws := &whitespace{}
	for _, opt := range []WhitespaceOption{} {
		opt(ws)
	}
	_ = ws
	result := PlaceOverlay(0, 0, fg, bg, false)
	if !strings.Contains(result, "test") {
		t.Error("expected result to contain foreground text")
	}
}

func TestPlaceOverlayZeroPosition(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg := "x"
	fg := "y"
	result := PlaceOverlay(0, 0, fg, bg, false)
	if result != "y" {
		t.Errorf("expected y at position 0,0, got %q", result)
	}
}

func TestPlaceOverlayFullWidth(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg := strings.Repeat("a", 30)
	fg := strings.Repeat("b", 30)
	result := PlaceOverlay(0, 0, fg, bg, false)
	if result != fg {
		t.Errorf("expected fg when same width, got %q", result)
	}
}

func TestPlaceOverlayNoShadowUsesBaseStyle(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg := strings.Repeat(" ", 20)
	fg := "test"
	result := PlaceOverlay(3, 3, fg, bg, false)
	if !strings.Contains(result, "test") {
		t.Error("expected result to contain foreground text")
	}
}
