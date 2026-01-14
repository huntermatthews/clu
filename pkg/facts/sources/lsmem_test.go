package sources

import (
	"fmt"
	"testing"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

var sampleLsmem = `RANGE   SIZE   STATE
0-3     4G     online
Total online memory:     4294967296`

// TestLsmemProvides ensures key registration.
func TestLsmemProvides(t *testing.T) {
	src := &Lsmem{}
	p := types.Provides{}
	src.Provides(p)
	if _, ok := p["phy.ram"]; !ok {
		t.Fatalf("phy.ram not provided")
	}
}

// TestLsmemSuccess ensures proper parsing and conversion.
func TestLsmemSuccess(t *testing.T) {
	orig := input.CommandRunner
	input.CommandRunner = func(cmdline string) (string, int, error) { return sampleLsmem, 0, nil }
	defer func() { input.CommandRunner = orig }()
	f := types.NewFactDB()
	src := &Lsmem{}
	src.Parse(f)
	got, ok := f.Get("phy.ram")
	if !ok || got == "" {
		t.Fatalf("expected phy.ram set got %q ok=%v", got, ok)
	}
	// 4294967296 bytes -> 4.00G (assuming BytesToSI formatting)
	if got != input.BytesToSI(4294967296) {
		t.Fatalf("unexpected conversion got %q", got)
	}
}

// TestLsmemFailure simulates command failure.
func TestLsmemFailure(t *testing.T) {
	orig := input.CommandRunner
	input.CommandRunner = func(cmdline string) (string, int, error) { return "", 1, fmt.Errorf("fail") }
	defer func() { input.CommandRunner = orig }()
	f := types.NewFactDB()
	src := &Lsmem{}
	src.Parse(f)
	got, _ := f.Get("phy.ram")
	if got != types.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got %q", got)
	}
}

// TestLsmemMissingLine ensures absence of target line yields ParseFailMsg.
func TestLsmemMissingLine(t *testing.T) {
	orig := input.CommandRunner
	input.CommandRunner = func(cmdline string) (string, int, error) { return "Header\nNo total line here", 0, nil }
	defer func() { input.CommandRunner = orig }()
	f := types.NewFactDB()
	src := &Lsmem{}
	src.Parse(f)
	got, _ := f.Get("phy.ram")
	if got != types.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got %q", got)
	}
}
