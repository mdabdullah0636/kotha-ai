package layout

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

type testKeyMap struct {
	A key.Binding
	B key.Binding
}

func TestKeyMapToSlice(t *testing.T) {
	t.Parallel()

	tm := testKeyMap{
		A: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "action a")),
		B: key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "action b")),
	}
	result := KeyMapToSlice(tm)
	if len(result) != 2 {
		t.Errorf("expected 2 bindings, got %d", len(result))
	}
	if !reflect.DeepEqual(result[0], tm.A) {
		t.Errorf("first binding mismatch")
	}
	if !reflect.DeepEqual(result[1], tm.B) {
		t.Errorf("second binding mismatch")
	}
}

func TestKeyMapToSliceEmpty(t *testing.T) {
	t.Parallel()

	type emptyKeyMap struct{}
	result := KeyMapToSlice(emptyKeyMap{})
	if len(result) != 0 {
		t.Errorf("expected 0 bindings, got %d", len(result))
	}
}

func TestKeyMapToSliceNonStruct(t *testing.T) {
	t.Parallel()

	result := KeyMapToSlice("not a struct")
	if result != nil {
		t.Errorf("expected nil for non-struct input, got %d bindings", len(result))
	}
}

func TestKeyMapToSliceSingleField(t *testing.T) {
	t.Parallel()

	type singleKeyMap struct {
		Only key.Binding
	}
	binding := key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "single"))
	result := KeyMapToSlice(singleKeyMap{Only: binding})
	if len(result) != 1 {
		t.Errorf("expected 1 binding, got %d", len(result))
	}
	if !reflect.DeepEqual(result[0], binding) {
		t.Errorf("binding mismatch")
	}
}
