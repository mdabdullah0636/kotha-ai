package styles

import (
	"testing"
)

func TestBaseStyle(t *testing.T) {
	t.Parallel()
	setupTheme()

	result := BaseStyle().Render("test")
	if result == "" {
		t.Fatal("expected non-empty render")
	}
}

func TestRegular(t *testing.T) {
	t.Parallel()
	setupTheme()

	result := Regular().Render("test")
	if result != "test" {
		t.Errorf("expected 'test', got %q", result)
	}
}

func TestMuted(t *testing.T) {
	t.Parallel()
	setupTheme()

	result := Muted().Render("test")
	if result == "" {
		t.Error("expected non-empty render")
	}
}

func TestBold(t *testing.T) {
	t.Parallel()
	setupTheme()

	result := Bold().Render("test")
	if result == "" {
		t.Error("expected non-empty render")
	}
}

func TestPadded(t *testing.T) {
	t.Parallel()
	setupTheme()

	result := Padded().Render("test")
	if result == "" {
		t.Error("expected non-empty render")
	}
}

func TestBorder(t *testing.T) {
	t.Parallel()
	setupTheme()

	result := Border().Render("test")
	if result == "" {
		t.Error("expected non-empty render")
	}
}

func TestThickBorder(t *testing.T) {
	t.Parallel()
	setupTheme()

	result := ThickBorder().Render("test")
	if result == "" {
		t.Error("expected non-empty render")
	}
}

func TestDoubleBorder(t *testing.T) {
	t.Parallel()
	setupTheme()

	result := DoubleBorder().Render("test")
	if result == "" {
		t.Error("expected non-empty render")
	}
}

func TestFocusedBorder(t *testing.T) {
	t.Parallel()
	setupTheme()

	result := FocusedBorder().Render("test")
	if result == "" {
		t.Error("expected non-empty render")
	}
}

func TestDimBorder(t *testing.T) {
	t.Parallel()
	setupTheme()

	result := DimBorder().Render("test")
	if result == "" {
		t.Error("expected non-empty render")
	}
}

func TestPrimaryColor(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := PrimaryColor()
	if c.Dark == "" && c.Light == "" {
		t.Error("expected non-empty color")
	}
}

func TestSecondaryColor(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := SecondaryColor()
	if c.Dark == "" && c.Light == "" {
		t.Error("expected non-empty color")
	}
}

func TestAccentColor(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := AccentColor()
	if c.Dark == "" && c.Light == "" {
		t.Error("expected non-empty color")
	}
}

func TestErrorColor(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := ErrorColor()
	if c.Dark == "" && c.Light == "" {
		t.Error("expected non-empty color")
	}
}

func TestWarningColor(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := WarningColor()
	if c.Dark == "" && c.Light == "" {
		t.Error("expected non-empty color")
	}
}

func TestSuccessColor(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := SuccessColor()
	if c.Dark == "" && c.Light == "" {
		t.Error("expected non-empty color")
	}
}

func TestInfoColor(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := InfoColor()
	if c.Dark == "" && c.Light == "" {
		t.Error("expected non-empty color")
	}
}

func TestTextColor(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := TextColor()
	if c.Dark == "" && c.Light == "" {
		t.Error("expected non-empty color")
	}
}

func TestTextMutedColor(t *testing.T) {
	t.Parallel()
	setupTheme()

	c := TextMutedColor()
	if c.Dark == "" && c.Light == "" {
		t.Error("expected non-empty color")
	}
}

func TestBackgroundColors(t *testing.T) {
	t.Parallel()
	setupTheme()

	bg := BackgroundColor()
	if bg.Dark == "" && bg.Light == "" {
		t.Error("expected non-empty color")
	}
	bg2 := BackgroundSecondaryColor()
	if bg2.Dark == "" && bg2.Light == "" {
		t.Error("expected non-empty color")
	}
	bg3 := BackgroundDarkerColor()
	if bg3.Dark == "" && bg3.Light == "" {
		t.Error("expected non-empty color")
	}
}

func TestBorderColors(t *testing.T) {
	t.Parallel()
	setupTheme()

	bn := BorderNormalColor()
	if bn.Dark == "" && bn.Light == "" {
		t.Error("expected non-empty color")
	}
	bf := BorderFocusedColor()
	if bf.Dark == "" && bf.Light == "" {
		t.Error("expected non-empty color")
	}
	bd := BorderDimColor()
	if bd.Dark == "" && bd.Light == "" {
		t.Error("expected non-empty color")
	}
}
