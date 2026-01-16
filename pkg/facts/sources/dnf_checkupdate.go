package sources

// Determines whether system package updates are required via `dnf check-update`.
// Exit code semantics (per dnf):
//   0   -> no updates available
//   100 -> updates available
//   other -> error / parsing failure

import (
	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/global"
	"github.com/huntermatthews/clu/pkg/input"
)

// DnfCheckUpdate reports if updates are required (True/False) or an error state.
type DnfCheckUpdate struct{}

var updateFact = types.Fact{
	Name: "run.update_required",
	Tier: types.TierOne,
}

// Provides registers the fact key supplied by this source.
func (d *DnfCheckUpdate) Provides(p types.Provides) {
	p[updateFact.Name] = d
}

// Requires declares the external command needed.
func (d *DnfCheckUpdate) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "dnf check-update")
}

// Parse executes `dnf check-update` and records whether updates are required.
// It respects a config key `net` (bool). If present and false, network queries
// are disabled and a placeholder value is stored.
func (d *DnfCheckUpdate) Parse(f *types.FactDB) {
	if !global.Config.NetEnabled {
		updateFact.Value = types.NetDisabledMsg
		f.AddFact(updateFact)
		return
	}

	_, rc, _ := input.CommandRunner("dnf check-update")

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
