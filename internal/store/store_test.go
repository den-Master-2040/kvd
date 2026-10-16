package store

import (
	"testing"
)

func TestStore(t *testing.T) {
	s := New()
	s.Set("key", []byte("value"))
	v, ok := s.Get("key")
	if !ok {
		t.Error("Expected key not found")
	}
	if string(v) != "value" {
		t.Errorf("Expected 'value', got '%s'", string(v))
	}
}
