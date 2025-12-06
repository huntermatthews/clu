package source_test

import (
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/sources"
)

func TestSelinuxProvides(t *testing.T) {
	src := &sources.Selinux{}
	p := pkg.Provides{}
	src.Provides(p)
	for _, k := range []string{"os.selinux.enable", "os.selinux.mode"} {
		if _, ok := p[k]; !ok {
			t.Fatalf("missing provides key %s", k)
		}
	}
}

func TestSelinuxEnableMapping(t *testing.T) {
	cases := []struct {
		name string
		rc   int
		want string
	}{
		{"enabled", 0, "True"},
		{"disabled", 1, "False"},
		{"error", 5, sources.ParseFailMsg},
	}
	orig := pkg.CommandRunner
	defer func() { pkg.CommandRunner = orig }()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pkg.CommandRunner = func(cmd string) (string, int) {
				if cmd == "selinuxenabled" {
					return "", c.rc
				}
				if cmd == "getenforce" {
					return "Enforcing", 0
				}
				return "", 1
			}
			f := pkg.NewFacts()
			src := &sources.Selinux{}
			src.Parse(f)
			gotEnable, _ := f.Get("os.selinux.enable")
			gotMode, _ := f.Get("os.selinux.mode")
			if gotEnable != c.want {
				t.Fatalf("enable rc=%d want %q got %q", c.rc, c.want, gotEnable)
			}
			if gotMode != "Enforcing" && gotMode != sources.ParseFailMsg { // depending on mapping
				t.Fatalf("unexpected mode value %q", gotMode)
			}
		})
	}
}

func TestSelinuxModeFailure(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmd string) (string, int) {
		if cmd == "selinuxenabled" {
			return "", 0
		}
		if cmd == "getenforce" {
			return "", 1
		}
		return "", 1
	}
	defer func() { pkg.CommandRunner = orig }()
	f := pkg.NewFacts()
	src := &sources.Selinux{}
	src.Parse(f)
	mode, _ := f.Get("os.selinux.mode")
	if mode != sources.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got %q", mode)
	}
}
