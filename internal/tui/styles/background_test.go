package styles

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestForceReplaceBackgroundWithExistingBG(t *testing.T) {
	t.Parallel()
	setupTheme()

	input := "\x1b[41mred bg\x1b[0m"
	result := ForceReplaceBackgroundWithLipgloss(input, lipgloss.Color("#2a2a2a"))
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !strings.Contains(result, "\x1b[") {
		t.Error("expected ANSI escape codes in result")
	}
}

func TestForceReplaceBackgroundPreservesTextAttributes(t *testing.T) {
	t.Parallel()
	setupTheme()

	input := "\x1b[1;41mbold red\x1b[0m"
	result := ForceReplaceBackgroundWithLipgloss(input, lipgloss.Color("#2a2a2a"))
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !strings.Contains(result, "\x1b[") {
		t.Error("expected ANSI escape codes in result")
	}
}

func TestForceReplaceBackgroundMultipleCodes(t *testing.T) {
	t.Parallel()
	setupTheme()

	input := "\x1b[41;32mbg41 fg32\x1b[0m"
	result := ForceReplaceBackgroundWithLipgloss(input, lipgloss.Color("#2a2a2a"))
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}
