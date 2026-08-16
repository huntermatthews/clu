// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package sources

// Go port of src/clu/sources/ipmitool.py
// Collects BMC/IPMI facts via `ipmitool mc info` and `ipmitool lan print`.
// Skips parsing unless platform is physical (facts["phy.platform"] == "physical").

import (
	"regexp"
	"strings"

	"github.com/huntermatthews/clu/internal/facts/types"
	"github.com/huntermatthews/clu/internal/input"
)

// Ipmitool gathers BMC firmware/network details.
type Ipmitool struct{}

var ipmitoolLanFacts = map[string]*types.Fact{
	"bmc.ipv4_source":  {Name: "bmc.ipv4_source", Origin: `(?m)^IP Address Source *: (.+)`, Tier: types.TierTwo},
	"bmc.ipv4_address": {Name: "bmc.ipv4_address", Origin: `(?m)^IP Address *: (.+)`, Tier: types.TierOne},
	"bmc.ipv4_mask":    {Name: "bmc.ipv4_mask", Origin: `(?m)^Subnet Mask *: (.+)`, Tier: types.TierTwo},
	"bmc.mac_address":  {Name: "bmc.mac_address", Origin: `(?m)^MAC Address *: (.+)`, Tier: types.TierOne},
}
var ipmitoolMcInfoFacts = map[string]*types.Fact{
	"bmc.firmware_version":  {Name: "bmc.firmware_version", Origin: `(?m)^Firmware Revision *: (.+)`, Tier: types.TierOne},
	"bmc.manufacturer_id":   {Name: "bmc.manufacturer_id", Origin: `(?m)^Manufacturer ID *: (.+)`, Tier: types.TierTwo},
	"bmc.manufacturer_name": {Name: "bmc.manufacturer_name", Origin: `(?m)^Manufacturer Name *: (.+)`, Tier: types.TierTwo},
}

// Provides registers all fact keys produced by this source.
func (i *Ipmitool) Provides(p types.Provides) {
	for name := range ipmitoolLanFacts {
		p[name] = i
	}
	for name := range ipmitoolMcInfoFacts {
		p[name] = i
	}
}

// Requires declares external program dependency.
func (i *Ipmitool) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "ipmitool")
}

// Parse orchestrates the two ipmitool queries if platform is physical.
func (i *Ipmitool) Parse(f *types.FactDB) {
	platform, _ := f.Get("phy.platform")
	if platform != "physical" { // skip if not physical hardware
		return
	}
	i.parseMcInfo(f)
	i.parseLanPrint(f)
}

func (i *Ipmitool) parseLanPrint(f *types.FactDB) {
	for _, fact := range ipmitoolLanFacts {
		fact.Value = types.ParseFailMsg
	}
	defer f.AddAll(ipmitoolLanFacts)

	data, rc, _ := input.CommandRunner("ipmitool lan print")
	if data == "" || rc != 0 {
		return
	}
	for _, fact := range ipmitoolLanFacts {
		re := regexp.MustCompile(fact.Origin)
		if m := re.FindStringSubmatch(data); len(m) == 2 {
			if val := strings.TrimSpace(m[1]); val != "" {
				fact.Value = val
			}
		}
	}
}

func (i *Ipmitool) parseMcInfo(f *types.FactDB) {
	for _, fact := range ipmitoolMcInfoFacts {
		fact.Value = types.ParseFailMsg
	}
	defer f.AddAll(ipmitoolMcInfoFacts)

	data, rc, _ := input.CommandRunner("ipmitool mc info")
	if data == "" || rc != 0 {
		return
	}
	for _, fact := range ipmitoolMcInfoFacts {
		re := regexp.MustCompile(fact.Origin)
		if m := re.FindStringSubmatch(data); len(m) == 2 {
			if val := strings.TrimSpace(m[1]); val != "" {
				fact.Value = val
			}
		}
	}
}
