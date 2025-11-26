package sources_test

import (
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/sources"
)

// TestDnfCheckUpdateProvides ensures the source registers its fact key.
func TestDnfCheckUpdateProvides(t *testing.T) {
	src := &sources.DnfCheckUpdate{}
	p := pkg.Provides{}
	src.Provides(p)
	if _, ok := p["run.update_required"]; !ok {
		t.Fatalf("expected run.update_required to be provided")
	}
}

// TestDnfCheckUpdateParseExitCodes validates rc -> value mapping.
func TestDnfCheckUpdateParseExitCodes(t *testing.T) {
	cases := []struct {
		name string
		rc   int
		want string
	}{
		{"no_updates", 0, "False"},
		{"updates", 100, "True"},
		{"error", 5, sources.ParseFailMsg},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Ensure net not disabled.
			pkg.SetConfig(map[string]interface{}{"net": true})

			// Stub command runner.
			orig := pkg.CommandRunner
			pkg.CommandRunner = func(cmdline string) (string, int) { return "", c.rc }
			defer func() { pkg.CommandRunner = orig }()

			f := pkg.NewFacts()
			src := &sources.DnfCheckUpdate{}
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
	pkg.SetConfig(map[string]interface{}{"net": false})

	// Stub command runner that would otherwise signal updates.
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmdline string) (string, int) { return "", 100 }
	defer func() { pkg.CommandRunner = orig }()

	f := pkg.NewFacts()
	src := &sources.DnfCheckUpdate{}
	src.Parse(f)
	got, _ := f.Get("run.update_required")
	if got != sources.NetDisabledMsg {
		t.Fatalf("expected NetDisabledMsg got %q", got)
	}
}
