package sources

// macOS update check via `softwareupdate --list`
// If your system is current, the output will state: "No new software available".

import (
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/global"
	"github.com/NHGRI/clu/pkg/input"
)

// MacOSCheckUpdate reports if updates are required (True/False) or an error state.
type MacOSCheckUpdate struct{}

var macUpdateFact = types.Fact{
	Name: "run.update_required",
	Tier: types.TierOne,
}

// Provides registers the fact key supplied by this source.
func (m *MacOSCheckUpdate) Provides(p types.Provides) {
	p[macUpdateFact.Name] = m
}

// Requires declares the external command needed.
func (m *MacOSCheckUpdate) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "softwareupdate --list")
}

// Parse executes `softwareupdate --list` and records whether updates are required.
// It respects a config key `net` (bool). If present and false, network queries
// are disabled and a placeholder value is stored.
func (m *MacOSCheckUpdate) Parse(f *types.FactDB) {
	if !global.Config.NetEnabled {
		macUpdateFact.Value = types.NetDisabledMsg
		f.AddFact(macUpdateFact)
		return
	}

	output, _, err := input.CommandRunner("softwareupdate --list")
	if err != nil {
		macUpdateFact.Value = types.ParseFailMsg
		f.AddFact(macUpdateFact)
		return
	}

	// If output contains "No new software available", no updates are needed
	if strings.Contains(output, "No new software available") {
		macUpdateFact.Value = "False"
	} else {
		macUpdateFact.Value = "True"
	}

	f.AddFact(macUpdateFact)
}

// TODO: Add Homebrew update check
// Refresh the package list:
// brew update
// See outdated packages:
// brew outdated
// This command specifically lists only the packages that have a newer version available.
