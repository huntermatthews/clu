package pkg

import "testing"

func TestNewFactsAndAdd(t *testing.T) {
	f := NewFacts()
	if f == nil {
		t.Fatal("NewFacts returned nil map")
	}
	f.Add("os", "linux")
	if got, ok := f["os"]; !ok || got != "linux" {
		t.Fatalf("expected os=linux got %q ok=%v", got, ok)
	}
}

func TestLiteralAndMake(t *testing.T) {
	f := Facts{}
	f.Add("a", "b")
	if f["a"] != "b" {
		t.Fatalf("expected a=b got %q", f["a"])
	}

	f2 := make(Facts)
	f2.Add("x", "y")
	if f2["x"] != "y" {
		t.Fatalf("expected x=y got %q", f2["x"])
	}
}

func TestAddOnNilPanics(t *testing.T) {
	var f Facts
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when calling Add on nil Facts")
		}
	}()
	// This should panic because f is nil
	f.Add("p", "q")
}
