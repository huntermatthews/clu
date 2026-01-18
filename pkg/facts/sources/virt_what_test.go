package sources

import (
	"fmt"
	"testing"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

func TestVirtWhatProvides(t *testing.T) {
	src := &VirtWhat{}
	p := types.NewProvides()
	src.Provides(p)
	if _, ok := p["phy.platform"]; !ok {
		t.Fatalf("missing provides key phy.platform")
	}
}

func TestVirtWhatPhysicalFallback(t *testing.T) {
	orig := input.CommandRunner
	input.CommandRunner = func(cmd string) (string, int, error) { return "\n\n", 0, nil }
	defer func() { input.CommandRunner = orig }()
	f := types.NewFactDB()
	src := &VirtWhat{}
	src.Parse(f)
	got, _ := f.Get("phy.platform")
	if got != "physical" {
		t.Fatalf("expected physical got %q", got)
	}
}

func TestVirtWhatVirtualizationList(t *testing.T) {
	orig := input.CommandRunner
	input.CommandRunner = func(cmd string) (string, int, error) { return "kvm\nvmware", 0, nil }
	defer func() { input.CommandRunner = orig }()
	f := types.NewFactDB()
	src := &VirtWhat{}
	src.Parse(f)
	got, _ := f.Get("phy.platform")
	if got != "kvm, vmware" {
		t.Fatalf("expected joined list got %q", got)
	}
}

func TestVirtWhatFailureRc(t *testing.T) {
	orig := input.CommandRunner
	input.CommandRunner = func(cmd string) (string, int, error) { return "", 5, fmt.Errorf("fail") }
	defer func() { input.CommandRunner = orig }()
	f := types.NewFactDB()
	src := &VirtWhat{}
	src.Parse(f)
	got, _ := f.Get("phy.platform")
	if got != types.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got %q", got)
	}
}

func TestVirtWhatSkipIfPreset(t *testing.T) {
	orig := input.CommandRunner
	// Would have produced virtualization list, but should be skipped.
	input.CommandRunner = func(cmd string) (string, int, error) { return "kvm", 0, nil }
	defer func() { input.CommandRunner = orig }()
	f := types.NewFactDB()
	f.AddFact(types.Fact{Name: "phy.platform", Value: "physical"})
	src := &VirtWhat{}
	src.Parse(f)
	got, _ := f.Get("phy.platform")
	if got != "physical" {
		t.Fatalf("expected preset value retained got %q", got)
	}
}
