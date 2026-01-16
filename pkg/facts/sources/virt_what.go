package sources

// Go port of src/clu/sources/virt_what.py
// Uses `virt-what` to determine virtualization platform(s); empty output implies physical.

import (
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

// VirtWhat determines phy.platform if not already set.
type VirtWhat struct{}

var platformFact = types.Fact{
	Name: "phy.platform",
	Tier: types.TierOne,
}

// Provides registers phy.platform key.
func (v *VirtWhat) Provides(p types.Provides) { p[platformFact.Name] = v }

// Requires declares program dependency.
func (v *VirtWhat) Requires(r *types.Requires) { r.Programs = append(r.Programs, "virt-what") }

// Parse sets phy.platform unless already present. Non-zero rc -> types.ParseFailMsg.
// Multiple lines are joined by ", "; empty result -> "physical".
func (v *VirtWhat) Parse(f *types.FactDB) {
	if f.Contains("phy.platform") {
		return
	}

	data, rc, _ := input.CommandRunner("virt-what")
	if rc != 0 {
		platformFact.Value = types.ParseFailMsg
		f.AddFact(platformFact)
		return
	}

	data = strings.TrimSpace(data)
	if data == "" {
		platformFact.Value = "physical"
	} else {
		// Replace newlines with comma-space
		platformFact.Value = strings.ReplaceAll(data, "\n", ", ")
	}
	f.AddFact(platformFact)
}
