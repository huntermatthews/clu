package sources_test

import (
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/sources"
)

var sampleLsmem = `RANGE   SIZE   STATE
0-3     4G     online
Total online memory:     4294967296`

// TestLsmemProvides ensures key registration.
func TestLsmemProvides(t *testing.T) {
	src := &sources.Lsmem{}
	p := pkg.Provides{}
	src.Provides(p)
	if _, ok := p["phy.ram"]; !ok {
		t.Fatalf("phy.ram not provided")
	}
}

// TestLsmemSuccess ensures proper parsing and conversion.
func TestLsmemSuccess(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmdline string) (string, int) { return sampleLsmem, 0 }
	defer func() { pkg.CommandRunner = orig }()
	f := pkg.NewFacts()
	src := &sources.Lsmem{}
	src.Parse(f)
	got, ok := f.Get("phy.ram")
	if !ok || got == "" {
		t.Fatalf("expected phy.ram set got %q ok=%v", got, ok)
	}
	// 4294967296 bytes -> 4.00G (assuming BytesToSI formatting)
	if got != pkg.BytesToSI(4294967296) {
		t.Fatalf("unexpected conversion got %q", got)
	}
}

// TestLsmemFailure simulates command failure.
func TestLsmemFailure(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmdline string) (string, int) { return "", 1 }
	defer func() { pkg.CommandRunner = orig }()
	f := pkg.NewFacts()
	src := &sources.Lsmem{}
	src.Parse(f)
	got, _ := f.Get("phy.ram")
	if got != sources.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got %q", got)
	}
}

// TestLsmemMissingLine ensures absence of target line yields ParseFailMsg.
func TestLsmemMissingLine(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmdline string) (string, int) { return "Header\nNo total line here", 0 }
	defer func() { pkg.CommandRunner = orig }()
	f := pkg.NewFacts()
	src := &sources.Lsmem{}
	src.Parse(f)
	got, _ := f.Get("phy.ram")
	if got != sources.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got %q", got)
	}
}
