package source_test

import (
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/sources"
)

// Representative lscpu output snippet.
var sampleLscpu = `Architecture:                    x86_64
CPU(s):                          16
Model name:                      Intel(R) Xeon(R) CPU E5-2670 v3 @ 2.30GHz
Thread(s) per core:              2
Core(s) per socket:              8
Socket(s):                       1
Vendor ID:                       GenuineIntel`

// TestLscpuProvides checks keys registration.
func TestLscpuProvides(t *testing.T) {
	src := &sources.Lscpu{}
	p := pkg.Provides{}
	src.Provides(p)
	keys := []string{"phy.cpu.model", "phy.cpu.vendor", "phy.cpu.cores", "phy.cpu.threads", "phy.cpu.sockets"}
	for _, k := range keys {
		if _, ok := p[k]; !ok {
			t.Fatalf("missing provides key %s", k)
		}
	}
}

// TestLscpuSuccess verifies parsing and derived computations.
func TestLscpuSuccess(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmdline string) (string, int) { return sampleLscpu, 0 }
	defer func() { pkg.CommandRunner = orig }()
	f := pkg.NewFacts()
	src := &sources.Lscpu{}
	src.Parse(f)
	cases := map[string]string{
		"phy.cpu.model":   "Intel(R) Xeon(R) CPU E5-2670 v3 @ 2.30GHz",
		"phy.cpu.vendor":  "GenuineIntel",
		"phy.cpu.cores":   "8",
		"phy.cpu.threads": "16",
		"phy.cpu.sockets": "1",
	}
	for k, want := range cases {
		got, ok := f.Get(k)
		if !ok || got != want {
			t.Fatalf("%s want %q got %q (ok=%v)", k, want, got, ok)
		}
	}
}

// TestLscpuFailure ensures no facts set on command failure.
func TestLscpuFailure(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmdline string) (string, int) { return "", 1 }
	defer func() { pkg.CommandRunner = orig }()
	f := pkg.NewFacts()
	src := &sources.Lscpu{}
	src.Parse(f)
	if _, ok := f.Get("phy.cpu.model"); ok {
		t.Fatalf("expected no facts after failure")
	}
}

// TestLscpuPartialMissing verifies derived computations fallback when numeric fields missing.
func TestLscpuPartialMissing(t *testing.T) {
	// Output missing threads_per_core -> threads should be ParseFailMsg.
	partial := `Model name: Intel(R) Sample CPU\nCore(s) per socket: 4\nSocket(s): 2\nVendor ID: VendorX`
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmdline string) (string, int) { return partial, 0 }
	defer func() { pkg.CommandRunner = orig }()
	f := pkg.NewFacts()
	src := &sources.Lscpu{}
	src.Parse(f)
	cores, _ := f.Get("phy.cpu.cores")
	threads, _ := f.Get("phy.cpu.threads")
	if cores != "8" { // 4 * 2
		t.Fatalf("expected cores=8 got %s", cores)
	}
	if threads != sources.ParseFailMsg {
		t.Fatalf("expected threads ParseFailMsg got %s", threads)
	}
}
