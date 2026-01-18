package sources

// Go port of src/clu/sources/ipmitool.py
// Collects BMC/IPMI facts via `ipmitool mc info` and `ipmitool lan print`.
// Skips parsing unless platform is physical (facts["phy.platform"] == "physical").

import (
	"regexp"
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
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
	data, rc, _ := input.CommandRunner("ipmitool lan print")
	if data == "" || rc != 0 {
		for _, fact := range ipmitoolLanFacts {
			fact.Value = types.ParseFailMsg
			f.AddFact(*fact)
		}
		return
	}
	// Extract values
	for _, fact := range ipmitoolLanFacts {
		re := regexp.MustCompile(fact.Origin)
		m := re.FindStringSubmatch(data)
		if len(m) == 2 {
			val := strings.TrimSpace(m[1])
			if val != "" {
				fact.Value = val
			}
		}
	}
	// Add all facts to the FactDB
	for _, fact := range ipmitoolLanFacts {
		f.AddFact(*fact)
	}
}

func (i *Ipmitool) parseMcInfo(f *types.FactDB) {
	data, rc, _ := input.CommandRunner("ipmitool mc info")
	if data == "" || rc != 0 {
		for _, fact := range ipmitoolMcInfoFacts {
			fact.Value = types.ParseFailMsg
			f.AddFact(*fact)
		}
		return
	}
	for _, fact := range ipmitoolMcInfoFacts {
		re := regexp.MustCompile(fact.Origin)
		m := re.FindStringSubmatch(data)
		if len(m) == 2 {
			val := strings.TrimSpace(m[1])
			if val != "" {
				fact.Value = val
			}
		}
	}
	// Add all facts to the FactDB
	for _, fact := range ipmitoolMcInfoFacts {
		f.AddFact(*fact)
	}
}
