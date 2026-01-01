package sources

import (
	"testing"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/global"
)

// TestDnfCheckUpdateProvides ensures the source registers its fact key.
func TestDnfCheckUpdateProvides(t *testing.T) {
	src := &DnfCheckUpdate{}
	p := types.NewProvides()
	src.Provides(p)
	if _, ok := p["run.update_required"]; !ok {
		t.Fatalf("expected run.update_required to be provided")
	}
}

// TestDnfCheckUpdateParseExitCodes validates rc -> value mapping.
func TestDnfCheckUpdateParseExitCodes(t *testing.T) {
	// t.Skip("TODO: fix expectations vs NetDisabledMsg")
	cases := []struct {
		name string
		rc   int
		want string
	}{
		{"no_updates", 0, "False"},
		{"updates", 100, "True"},
		{"error", 5, types.ParseFailMsg},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Ensure net not disabled.
			global.CluConfig.NetEnabled = true
			defer func() { global.CluConfig.NetEnabled = false }()

			// Stub command runner.
			orig := pkg.CommandRunner
			pkg.CommandRunner = func(cmdline string) (string, int, error) {
				return c.want, c.rc, nil
			}
			defer func() { pkg.CommandRunner = orig }()

			f := types.NewFacts()
			src := &DnfCheckUpdate{}
			src.Parse(f)
			got, _ := f.Get("run.update_required")
			if got != c.want {
				t.Fatalf("rc %d: expected %q got %q", c.rc, c.want, got)
			}
		})
	}
}

// TestDnfCheckUpdateNetDisabled confirms net gating short-circuits parse.
func TestDnfCheckUpdateNetDisabled(t *testing.T) {

	// Stub command runner that would otherwise signal updates.
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmdline string) (string, int, error) { return "", 100, nil }
	defer func() { pkg.CommandRunner = orig }()

	f := types.NewFacts()
	src := &DnfCheckUpdate{}
	src.Parse(f)
	got, _ := f.Get("run.update_required")
	if got != types.NetDisabledMsg {
		t.Fatalf("expected NetDisabledMsg got %q", got)
	}
}
