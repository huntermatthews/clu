// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package sources

// Go port of src/clu/sources/virt_what.py
// Uses `virt-what` to determine virtualization platform(s); empty output implies physical.

import (
	"strings"

	"github.com/huntermatthews/clu/internal/facts/types"
	"github.com/huntermatthews/clu/internal/input"
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

	platformFact.Value = types.ParseFailMsg
	defer func() { f.AddFact(platformFact) }()

	data, rc, _ := input.CommandRunner("virt-what")
	if rc != 0 {
		return
	}

	data = strings.TrimSpace(data)
	if data == "" {
		platformFact.Value = "physical"
	} else {
		platformFact.Value = strings.ReplaceAll(data, "\n", ", ")
	}
}
