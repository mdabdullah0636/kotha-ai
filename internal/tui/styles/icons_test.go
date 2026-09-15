package styles

import (
	"testing"
)

func TestIcons(t *testing.T) {
	t.Parallel()

	if KothaIcon == "" {
		t.Error("KothaIcon should not be empty")
	}
	if CheckIcon == "" {
		t.Error("CheckIcon should not be empty")
	}
	if ErrorIcon == "" {
		t.Error("ErrorIcon should not be empty")
	}
	if WarningIcon == "" {
		t.Error("WarningIcon should not be empty")
	}
	if HintIcon == "" {
		t.Error("HintIcon should not be empty")
	}
	if SpinnerIcon == "" {
		t.Error("SpinnerIcon should not be empty")
	}
	if LoadingIcon == "" {
		t.Error("LoadingIcon should not be empty")
	}
	if DocumentIcon == "" {
		t.Error("DocumentIcon should not be empty")
	}
	// InfoIcon is intentionally empty
}
