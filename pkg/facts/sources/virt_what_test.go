package sources

import (
	"testing"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/facts/types"
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
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmd string) (string, int) { return "\n\n", 0 }
	defer func() { pkg.CommandRunner = orig }()
	f := types.NewFacts()
	src := &VirtWhat{}
	src.Parse(f)
	got, _ := f.Get("phy.platform")
	if got != "physical" {
		t.Fatalf("expected physical got %q", got)
	}
}

func TestVirtWhatVirtualizationList(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmd string) (string, int) { return "kvm\nvmware", 0 }
	defer func() { pkg.CommandRunner = orig }()
	f := types.NewFacts()
	src := &VirtWhat{}
	src.Parse(f)
	got, _ := f.Get("phy.platform")
	if got != "kvm, vmware" {
		t.Fatalf("expected joined list got %q", got)
	}
}

func TestVirtWhatFailureRc(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmd string) (string, int) { return "", 5 }
	defer func() { pkg.CommandRunner = orig }()
	f := types.NewFacts()
	src := &VirtWhat{}
	src.Parse(f)
	got, _ := f.Get("phy.platform")
	if got != types.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got %q", got)
	}
}

func TestVirtWhatSkipIfPreset(t *testing.T) {
	orig := pkg.CommandRunner
	// Would have produced virtualization list, but should be skipped.
	pkg.CommandRunner = func(cmd string) (string, int) { return "kvm", 0 }
	defer func() { pkg.CommandRunner = orig }()
	f := types.NewFacts()
	f.Set("phy.platform", "physical")
	src := &VirtWhat{}
	src.Parse(f)
	got, _ := f.Get("phy.platform")
	if got != "physical" {
		t.Fatalf("expected preset value retained got %q", got)
	}
}
