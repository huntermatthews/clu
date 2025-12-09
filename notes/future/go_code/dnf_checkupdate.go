package sources

// Go port of src/clu/sources/dnf_checkupdate.py
// Determines whether system package updates are required via `dnf check-update`.
// Exit code semantics (per dnf):
//   0   -> no updates available
//   100 -> updates available
//   other -> error / parsing failure

import (
	pkg "github.com/huntermatthews/clu/pkg"
	facts "github.com/huntermatthews/clu/pkg/facts"
)

// DnfCheckUpdate reports if updates are required (True/False) or an error state.
type DnfCheckUpdate struct{}

// Provides registers the fact key supplied by this source.
func (d *DnfCheckUpdate) Provides(p facts.Provides) {
	p["run.update_required"] = d
}

// Requires declares the external command needed.
func (d *DnfCheckUpdate) Requires(r *facts.Requires) {
	r.Programs = append(r.Programs, "dnf check-update")
}

// Parse executes `dnf check-update` and records whether updates are required.
// It respects a config key `net` (bool). If present and false, network queries
// are disabled and a placeholder value is stored.
func (d *DnfCheckUpdate) Parse(f *facts.Facts) {
	if !pkg.Config.NetEnabled {
		f.Set("run.update_required", NetDisabledMsg)
		return
	}

	_, rc := pkg.CommandRunner("dnf check-update")

	var value string
	switch rc {
	case 0:
		value = "False"
	case 100:
		value = "True"
	default:
		value = ParseFailMsg
	}
	f.Set("run.update_required", value)
}
