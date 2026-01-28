package sources

import (
	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/global"
	"github.com/NHGRI/clu/pkg/input"
)

// LinuxCheckUpdate reports if updates are required (True/False) or an error state.
type LinuxCheckUpdate struct{}

var updateFact = types.Fact{
	Name: "run.update_required",
	Tier: types.TierOne,
}

// parseFunc is a function type for distro-specific update parsers.
type parseFunc func(*LinuxCheckUpdate, *types.FactDB)

// distroParserMap maps distribution families to their update check functions.
var distroParserMap = map[string]parseFunc{
	"fedora":   (*LinuxCheckUpdate).dnfParse,
	"debian":   (*LinuxCheckUpdate).aptParse,
	"arch":     (*LinuxCheckUpdate).pacmanParse,
	"opensuse": (*LinuxCheckUpdate).zypperParse,
}

// Provides registers the fact key supplied by this source.
func (d *LinuxCheckUpdate) Provides(p types.Provides) {
	p[updateFact.Name] = d
}

// Requires declares the external command needed.
func (d *LinuxCheckUpdate) Requires(r *types.Requires) {
	// fix this by just checking for each program in a loop
	// if we find one, we add it to the list

	// r.Programs = append(r.Programs, "dnf check-update")

}

// Parse executes `dnf check-update` and records whether updates are required.
// It respects a config key `net` (bool). If present and false, network queries
// are disabled and a placeholder value is stored.
func (d *LinuxCheckUpdate) Parse(f *types.FactDB) {
	if !global.Config.NetEnabled {
		updateFact.Value = types.NetDisabledMsg
		f.AddFact(updateFact)
		return
	}

	// Determine distro family to select appropriate package manager.
	distroFamily, _ := f.Get("os.distro.family")

	if parseFunc, ok := distroParserMap[distroFamily]; ok {
		parseFunc(d, f)
	} else {
		updateFact.Value = types.ParseFailMsg
		f.AddFact(updateFact)
	}
}

func (d *LinuxCheckUpdate) dnfParse(f *types.FactDB) {

	_, rc, _ := input.CommandRunner("dnf check-update")

	// Determines whether system package updates are required via `dnf check-update`.
	// Exit code semantics (per dnf):
	//   0   -> no updates available
	//   100 -> updates available
	//   other -> error / parsing failure
	switch rc {
	case 0:
		updateFact.Value = "False"
	case 100:
		updateFact.Value = "True"
	default:
		updateFact.Value = types.ParseFailMsg
	}
	f.AddFact(updateFact)
}

func (d *LinuxCheckUpdate) aptParse(f *types.FactDB) {
	// TODO: Implement apt-get update check
	updateFact.Value = types.ParseFailMsg
	f.AddFact(updateFact)
}

func (d *LinuxCheckUpdate) pacmanParse(f *types.FactDB) {
	// TODO: Implement pacman update check
	updateFact.Value = types.ParseFailMsg
	f.AddFact(updateFact)
}

func (d *LinuxCheckUpdate) zypperParse(f *types.FactDB) {
	// TODO: Implement zypper update check
	updateFact.Value = types.ParseFailMsg
	f.AddFact(updateFact)
}
