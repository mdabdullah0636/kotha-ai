package dialog

import "kotha/internal/tui/theme"

func setupTheme() {
	if theme.CurrentTheme() == nil {
		_ = theme.SetTheme("kotha")
	}
}
