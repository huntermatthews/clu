package sources_test

import (
	"errors"
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/sources"
)

func TestSysDmiProvides(t *testing.T) {
	src := &sources.SysDmi{}
	p := pkg.Provides{}
	src.Provides(p)
	keys := []string{"sys.vendor", "sys.model.family", "sys.model.name", "sys.serial_no", "sys.uuid", "sys.oem", "sys.asset_no"}
	for _, k := range keys {
		if _, ok := p[k]; !ok {
			t.Fatalf("missing provides key %s", k)
		}
	}
}

func TestSysDmiSkipNonPhysical(t *testing.T) {
	orig := pkg.FileReader
	pkg.FileReader = func(path string) (string, error) { return "Dummy", nil }
	defer func() { pkg.FileReader = orig }()
	f := pkg.NewFacts()
	f.Set("phy.platform", "virtual")
	src := &sources.SysDmi{}
	src.Parse(f)
	if _, ok := f.Get("sys.vendor"); ok {
		t.Fatalf("expected no DMI facts on non-physical platform")
	}
}

func TestSysDmiSuccess(t *testing.T) {
	// Simulate different file contents via map
	contents := map[string]string{
		"/sys/devices/virtual/dmi/id/sys_vendor":        "AcmeCorp\n",
		"/sys/devices/virtual/dmi/id/product_family":    "FamilyX\n",
		"/sys/devices/virtual/dmi/id/product_name":      "ModelY\n",
		"/sys/devices/virtual/dmi/id/product_serial":    "SN123\n",
		"/sys/devices/virtual/dmi/id/product_uuid":      "UUID-ABC\n",
		"/sys/devices/virtual/dmi/id/chassis_vendor":    "ChassisCo\n",
		"/sys/devices/virtual/dmi/id/chassis_asset_tag": "Asset42\n",
	}
	orig := pkg.FileReader
	pkg.FileReader = func(path string) (string, error) {
		if v, ok := contents[path]; ok {
			return v, nil
		}
		return "", errors.New("missing")
	}
	defer func() { pkg.FileReader = orig }()
	f := pkg.NewFacts()
	f.Set("phy.platform", "physical")
	src := &sources.SysDmi{}
	src.Parse(f)
	cases := map[string]string{
		"sys.vendor":       "AcmeCorp",
		"sys.model.family": "FamilyX",
		"sys.model.name":   "ModelY",
		"sys.serial_no":    "SN123",
		"sys.uuid":         "UUID-ABC",
		"sys.oem":          "ChassisCo",
		"sys.asset_no":     "Asset42",
	}
	for k, want := range cases {
		got, ok := f.Get(k)
		if !ok || got != want {
			t.Fatalf("%s want %q got %q ok=%v", k, want, got, ok)
		}
	}
}

func TestSysDmiMissingFiles(t *testing.T) {
	orig := pkg.FileReader
	pkg.FileReader = func(path string) (string, error) { return "", errors.New("missing") }
	defer func() { pkg.FileReader = orig }()
	f := pkg.NewFacts()
	f.Set("phy.platform", "physical")
	src := &sources.SysDmi{}
	src.Parse(f)
	got, _ := f.Get("sys.vendor")
	if got != "" {
		t.Fatalf("expected empty string for missing file got %q", got)
	}
}
