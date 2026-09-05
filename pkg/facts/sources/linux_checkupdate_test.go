// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package sources

import (
	"testing"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/global"
	"github.com/huntermatthews/clu/pkg/input"
)

// TestDnfCheckUpdateProvides ensures the source registers its fact key.
func TestDnfCheckUpdateProvides(t *testing.T) {
	src := &LinuxCheckUpdate{}
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
			global.Config.NetEnabled = true
			defer func() { global.Config.NetEnabled = false }()

			// Stub command runner.
			orig := input.CommandRunner
			input.CommandRunner = func(cmdline string) (string, int, error) {
				return c.want, c.rc, nil
			}
			defer func() { input.CommandRunner = orig }()

			f := types.NewFactDB()
			// LinuxCheckUpdate.Parse routes on os.distro.family; set fedora to exercise dnfParse.
			f.AddFact(types.Fact{Name: "os.distro.family", Value: "fedora", Tier: types.TierTwo})
			src := &LinuxCheckUpdate{}
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
	orig := input.CommandRunner
	input.CommandRunner = func(cmdline string) (string, int, error) { return "", 100, nil }
	defer func() { input.CommandRunner = orig }()

	f := types.NewFactDB()
	src := &LinuxCheckUpdate{}
	src.Parse(f)
	got, _ := f.Get("run.update_required")
	if got != types.NetDisabledMsg {
		t.Fatalf("expected NetDisabledMsg got %q", got)
	}
}
