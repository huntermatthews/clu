package sources

import (
	"testing"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/facts/types"
)

func TestUdevadmRamProvides(t *testing.T) {
	src := &UdevadmRam{}
	p := types.NewProvides()
	src.Provides(p)
	if _, ok := p["phy.ram"]; !ok {
		t.Fatalf("missing provides key phy.ram")
	}
}

func TestUdevadmRamSuccess(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmd string) (string, int) {
		return "MEMORY_DEVICE_0_SIZE=1048576\nMEMORY_DEVICE_1_SIZE=2097152\n", 0
	}
	defer func() { pkg.CommandRunner = orig }()
	f := types.NewFacts()
	src := &UdevadmRam{}
	src.Parse(f)
	got, _ := f.Get("phy.ram")
	// total = 3145728 bytes
	expected := pkg.BytesToSI(3145728)
	if got != expected {
		t.Fatalf("want %q got %q", expected, got)
	}
}

func TestUdevadmRamNoMatches(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmd string) (string, int) { return "OTHER=1", 0 }
	defer func() { pkg.CommandRunner = orig }()
	f := types.NewFacts()
	src := &UdevadmRam{}
	src.Parse(f)
	got, _ := f.Get("phy.ram")
	expected := pkg.BytesToSI(0)
	if got != expected {
		t.Fatalf("expected zero formatting %q got %q", expected, got)
	}
}

func TestUdevadmRamFailure(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmd string) (string, int) { return "", 1 }
	defer func() { pkg.CommandRunner = orig }()
	f := types.NewFacts()
	src := &UdevadmRam{}
	src.Parse(f)
	got, _ := f.Get("phy.ram")
	if got != types.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got %q", got)
	}
}
