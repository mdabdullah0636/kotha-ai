package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOverlayManager_SetAndView(t *testing.T) {
	m := NewOverlayManager()

	m.Set("test", true, func() string { return "overlay" })
	result := m.View("background")
	assert.Contains(t, result, "overlay")
	assert.Contains(t, result, "background")
}

func TestOverlayManager_Hide(t *testing.T) {
	m := NewOverlayManager()

	m.Set("test", true, func() string { return "overlay" })
	m.Hide("test")
	result := m.View("background")
	assert.Equal(t, "background", result)
}

func TestOverlayManager_Remove(t *testing.T) {
	m := NewOverlayManager()

	m.Set("test", true, func() string { return "overlay" })
	m.Remove("test")
	assert.False(t, m.Has("test"))
	assert.Empty(t, m.VisibleNames())
}

func TestOverlayManager_Multiple(t *testing.T) {
	m := NewOverlayManager()

	m.Set("first", true, func() string { return "FIRST" })
	m.Set("second", true, func() string { return "SECOND" })
	result := m.View("background")
	assert.Contains(t, result, "FIRST")
	assert.Contains(t, result, "SECOND")
}

func TestOverlayManager_VisibleNames(t *testing.T) {
	m := NewOverlayManager()

	m.Set("visible", true, func() string { return "x" })
	m.Set("hidden", false, func() string { return "y" })

	names := m.VisibleNames()
	assert.Equal(t, []string{"visible"}, names)
}

func TestOverlayManager_Clear(t *testing.T) {
	m := NewOverlayManager()

	m.Set("test", true, func() string { return "x" })
	m.Clear()
	assert.Empty(t, m.VisibleNames())
}

func TestOverlayManager_EmptyView(t *testing.T) {
	m := NewOverlayManager()
	result := m.View("background")
	assert.Equal(t, "background", result)
}

func TestOverlayManager_OverlaysWithEmptyView(t *testing.T) {
	m := NewOverlayManager()

	m.Set("first", true, func() string { return "FIRST" })
	m.Set("second", true, func() string { return "" })
	result := m.View("background")
	assert.Contains(t, result, "FIRST")
	assert.NotContains(t, result, "SECOND")
}

func TestOverlayManager_NonOverlayHeight(t *testing.T) {
	m := NewOverlayManager()
	m.Set("test", true, func() string { return "x" })

	result := m.View("")
	assert.NotEmpty(t, result)
}

func TestOverlayManager_OrderPreserved(t *testing.T) {
	m := NewOverlayManager()

	m.Set("a", true, func() string { return "A" })
	m.Set("b", true, func() string { return "B" })
	m.Set("c", true, func() string { return "C" })

	names := m.VisibleNames()
	require.Len(t, names, 3)
	assert.Equal(t, "a", names[0])
	assert.Equal(t, "b", names[1])
	assert.Equal(t, "c", names[2])
}

func TestOverlayManager_HideThenShow(t *testing.T) {
	m := NewOverlayManager()

	m.Set("test", true, func() string { return "overlay" })
	m.Hide("test")
	assert.False(t, m.Has("test"))

	m.Show("test")
	assert.True(t, m.Has("test"))
	result := m.View("background")
	assert.Contains(t, result, "overlay")
}

func TestOverlayManager_Toggle(t *testing.T) {
	m := NewOverlayManager()

	m.Set("test", true, func() string { return "overlay" })
	m.Hide("test")
	m.Hide("test")
	assert.False(t, m.Has("test"))
}
