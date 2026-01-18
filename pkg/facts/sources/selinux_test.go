package sources

import (
	"fmt"
	"testing"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

func TestSelinuxProvides(t *testing.T) {
	src := &Selinux{}
	p := types.Provides{}
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
		{"error", 5, types.ParseFailMsg},
	}
	orig := input.CommandRunner
	defer func() { input.CommandRunner = orig }()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			input.CommandRunner = func(cmd string) (string, int, error) {
				if cmd == "selinuxenabled" {
					if c.rc == 0 {
						return "", 0, nil
					}
					return "", c.rc, fmt.Errorf("rc %d", c.rc)
				}
				if cmd == "getenforce" {
					return "Enforcing", 0, nil
				}
				return "", 1, fmt.Errorf("fail")
			}
			f := types.NewFactDB()
			src := &Selinux{}
			src.Parse(f)
			gotEnable, _ := f.Get("os.selinux.enable")
			gotMode, _ := f.Get("os.selinux.mode")
			if gotEnable != c.want {
				t.Fatalf("enable rc=%d want %q got %q", c.rc, c.want, gotEnable)
			}
			if gotMode != "Enforcing" && gotMode != types.ParseFailMsg {
				t.Fatalf("unexpected mode value %q", gotMode)
			}
		})
	}
}

func TestSelinuxModeFailure(t *testing.T) {
	orig := input.CommandRunner
	input.CommandRunner = func(cmd string) (string, int, error) {
		if cmd == "selinuxenabled" {
			return "", 0, nil
		}
		if cmd == "getenforce" {
			return "", 1, fmt.Errorf("fail")
		}
		return "", 1, fmt.Errorf("fail")
	}
	defer func() { input.CommandRunner = orig }()
	f := types.NewFactDB()
	src := &Selinux{}
	src.Parse(f)
	mode, _ := f.Get("os.selinux.mode")
	if mode != types.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got %q", mode)
	}
}
