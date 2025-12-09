package sources

import (
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/sources"
)

// sample outputs modeled after ipmitool output patterns
var sampleLanPrint = `IP Address Source       : DHCP
IP Address              : 192.0.2.10
Subnet Mask             : 255.255.255.0
MAC Address             : 00:11:22:33:44:55`

var sampleMcInfo = `Firmware Revision       : 3.45
Manufacturer ID         : 343
Manufacturer Name       : ExampleCorp`

// TestIpmitoolProvides ensures keys are registered.
func TestIpmitoolProvides(t *testing.T) {
	src := &sources.Ipmitool{}
	p := pkg.Provides{}
	src.Provides(p)
	keys := []string{
		"bmc.ipv4_source", "bmc.ipv4_address", "bmc.ipv4_mask", "bmc.mac_address",
		"bmc.firmware_version", "bmc.manufacturer_id", "bmc.manufacturer_name",
	}
	for _, k := range keys {
		if _, ok := p[k]; !ok {
			t.Fatalf("missing provides key %s", k)
		}
	}
}

// TestIpmitoolSkipNonPhysical verifies no facts are set when platform != physical.
func TestIpmitoolSkipNonPhysical(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmdline string) (string, int) { return sampleLanPrint, 0 }
	defer func() { pkg.CommandRunner = orig }()

	f := pkg.NewFacts()
	f.Set("phy.platform", "virtual")
	src := &sources.Ipmitool{}
	src.Parse(f)
	if v, ok := f.Get("bmc.mac_address"); ok || v != "" {
		t.Fatalf("expected no bmc facts on non-physical platform")
	}
}

// TestIpmitoolSuccess ensures proper extraction for both commands.
func TestIpmitoolSuccess(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmdline string) (string, int) {
		switch cmdline {
		case "ipmitool lan print":
			return sampleLanPrint, 0
		case "ipmitool mc info":
			return sampleMcInfo, 0
		default:
			return "", 1
		}
	}
	defer func() { pkg.CommandRunner = orig }()

	f := pkg.NewFacts()
	f.Set("phy.platform", "physical")
	src := &sources.Ipmitool{}
	src.Parse(f)

	checks := map[string]string{
		"bmc.ipv4_source":       "DHCP",
		"bmc.ipv4_address":      "192.0.2.10",
		"bmc.ipv4_mask":         "255.255.255.0",
		"bmc.mac_address":       "00:11:22:33:44:55",
		"bmc.firmware_version":  "3.45",
		"bmc.manufacturer_id":   "343",
		"bmc.manufacturer_name": "ExampleCorp",
	}
	for k, want := range checks {
		got, ok := f.Get(k)
		if !ok || got != want {
			t.Fatalf("key %s: want %q got %q (ok=%v)", k, want, got, ok)
		}
	}
}

// TestIpmitoolLanFailure ensures lan print failure assigns ParseFailMsg to lan keys only.
func TestIpmitoolLanFailure(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmdline string) (string, int) {
		if cmdline == "ipmitool lan print" {
			return "", 1 // simulate failure
		}
		if cmdline == "ipmitool mc info" {
			return sampleMcInfo, 0
		}
		return "", 1
	}
	defer func() { pkg.CommandRunner = orig }()

	f := pkg.NewFacts()
	f.Set("phy.platform", "physical")
	src := &sources.Ipmitool{}
	src.Parse(f)

	lanKeys := []string{"bmc.ipv4_source", "bmc.ipv4_address", "bmc.ipv4_mask", "bmc.mac_address"}
	for _, k := range lanKeys {
		got, _ := f.Get(k)
		if got != sources.ParseFailMsg {
			t.Fatalf("expected ParseFailMsg for %s got %q", k, got)
		}
	}
	// mc info keys should be populated
	if v, _ := f.Get("bmc.firmware_version"); v != "3.45" {
		t.Fatalf("expected firmware version got %q", v)
	}
}

// TestIpmitoolMcFailure ensures mc info failure assigns ParseFailMsg to mc info keys only.
func TestIpmitoolMcFailure(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmdline string) (string, int) {
		if cmdline == "ipmitool lan print" {
			return sampleLanPrint, 0
		}
		if cmdline == "ipmitool mc info" {
			return "", 1 // failure
		}
		return "", 1
	}
	defer func() { pkg.CommandRunner = orig }()

	f := pkg.NewFacts()
	f.Set("phy.platform", "physical")
	src := &sources.Ipmitool{}
	src.Parse(f)

	mcKeys := []string{"bmc.firmware_version", "bmc.manufacturer_id", "bmc.manufacturer_name"}
	for _, k := range mcKeys {
		got, _ := f.Get(k)
		if got != sources.ParseFailMsg {
			t.Fatalf("expected ParseFailMsg for %s got %q", k, got)
		}
	}
	// lan keys should be populated
	if v, _ := f.Get("bmc.ipv4_address"); v != "192.0.2.10" {
		t.Fatalf("expected ipv4 address got %q", v)
	}
}
