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

	macosUpdatesAvailable, err := hasMacosUpdatesAvailable()
	if err != nil {
		macUpdateFact.Value = types.ParseFailMsg
		f.AddFact(macUpdateFact)
		return
	}

	brewUpdatesAvailable, err := hasBrewUpdatesAvailable()
	if err != nil {
		macUpdateFact.Value = types.ParseFailMsg
		f.AddFact(macUpdateFact)
		return
	}

	if macosUpdatesAvailable || brewUpdatesAvailable {
		macUpdateFact.Value = "True"
	} else {
		macUpdateFact.Value = "False"
	}

	f.AddFact(macUpdateFact)
}

// Check software update alone and return a boolean
func hasMacosUpdatesAvailable() (bool, error) {
	output, _, err := input.CommandRunner("softwareupdate --list")
	if err != nil {
		return false, err
	}

	// If output contains "No new software available", no updates are needed
	if strings.Contains(output, "No new software available") {
		return false, nil
	} else {
		return true, nil
	}
}

// Check software update alone and return a boolean
func hasBrewUpdatesAvailable() (bool, error) {
	// First update Homebrew metadata
	_, rc, err := input.CommandRunner("brew update")
	if err != nil || rc != 0 {
		return false, err
	}

	// Check for outdated packages
	output, _, err := input.CommandRunner("brew outdated")
	if err != nil {
		return false, err
	}

	// If output is empty, no updates are needed
	if strings.TrimSpace(output) == "" {
		return false, nil
	} else {
		return true, nil
	}
}


