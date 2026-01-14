package sources

import (
	"fmt"
	"testing"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
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
	orig := input.CommandRunner
	input.CommandRunner = func(cmd string) (string, int, error) {
		return "MEMORY_DEVICE_0_SIZE=1048576\nMEMORY_DEVICE_1_SIZE=2097152\n", 0, nil
	}
	defer func() { input.CommandRunner = orig }()
	f := types.NewFactDB()
	src := &UdevadmRam{}
	src.Parse(f)
	got, _ := f.Get("phy.ram")
	// total = 3145728 bytes
	expected := input.BytesToSI(3145728)
	if got != expected {
		t.Fatalf("want %q got %q", expected, got)
	}
}

func TestUdevadmRamNoMatches(t *testing.T) {
	orig := input.CommandRunner
	input.CommandRunner = func(cmd string) (string, int, error) { return "OTHER=1", 0, nil }
	defer func() { input.CommandRunner = orig }()
	f := types.NewFactDB()
	src := &UdevadmRam{}
	src.Parse(f)
	got, _ := f.Get("phy.ram")
	expected := input.BytesToSI(0)
	if got != expected {
		t.Fatalf("expected zero formatting %q got %q", expected, got)
	}
}

func TestUdevadmRamFailure(t *testing.T) {
	orig := input.CommandRunner
	input.CommandRunner = func(cmd string) (string, int, error) { return "", 1, fmt.Errorf("fail") }
	defer func() { input.CommandRunner = orig }()
	f := types.NewFactDB()
	src := &UdevadmRam{}
	src.Parse(f)
	got, _ := f.Get("phy.ram")
	if got != types.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got %q", got)
	}
}
