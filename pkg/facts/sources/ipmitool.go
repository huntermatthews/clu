package sources

// Go port of src/clu/sources/ipmitool.py
// Collects BMC/IPMI facts via `ipmitool mc info` and `ipmitool lan print`.
// Skips parsing unless platform is physical (facts["phy.platform"] == "physical").

import (
	"regexp"
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

// Ipmitool gathers BMC firmware/network details.
type Ipmitool struct{}

var ipmitoolLanFacts = []types.Fact{
	{Name: "bmc.ipv4_source", Value: "", Origin: `(?m)^IP Address Source *: (.+)`, Tier: types.TierTwo},
	{Name: "bmc.ipv4_address", Value: "", Origin: `(?m)^IP Address *: (.+)`, Tier: types.TierOne},
	{Name: "bmc.ipv4_mask", Value: "", Origin: `(?m)^Subnet Mask *: (.+)`, Tier: types.TierTwo},
	{Name: "bmc.mac_address", Value: "", Origin: `(?m)^MAC Address *: (.+)`, Tier: types.TierOne},
}
var ipmitoolMcInfoFacts = []types.Fact{
	{Name: "bmc.firmware_version", Value: "", Origin: `(?m)^Firmware Revision *: (.+)`, Tier: types.TierOne},
	{Name: "bmc.manufacturer_id", Value: "", Origin: `(?m)^Manufacturer ID *: (.+)`, Tier: types.TierTwo},
	{Name: "bmc.manufacturer_name", Value: "", Origin: `(?m)^Manufacturer Name *: (.+)`, Tier: types.TierTwo},
}

// Provides registers all fact keys produced by this source.
func (i *Ipmitool) Provides(p types.Provides) {
	for _, k := range ipmitoolLanFacts {
		p[k.Name] = i
	}
	for _, k := range ipmitoolMcInfoFacts {
		p[k.Name] = i
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
		for _, k := range ipmitoolLanFacts {
			f.Add(k.Tier, k.Name, types.ParseFailMsg)
		}
		return
	}
	// Extract values
	for _, k := range ipmitoolLanFacts {
		re := regexp.MustCompile(k.Origin)
		m := re.FindStringSubmatch(data)
		if len(m) == 2 {
			val := strings.TrimSpace(m[1])
			if val != "" {
				f.Add(k.Tier, k.Name, val)
			}
		}
	}
}

func (i *Ipmitool) parseMcInfo(f *types.FactDB) {

	data, rc, _ := input.CommandRunner("ipmitool mc info")
	if data == "" || rc != 0 {
		for _, k := range ipmitoolMcInfoFacts {
			f.Add(k.Tier, k.Name, types.ParseFailMsg)
		}
		return
	}
	for _, k := range ipmitoolMcInfoFacts {
		re := regexp.MustCompile(k.Origin)
		m := re.FindStringSubmatch(data)
		if len(m) == 2 {
			val := strings.TrimSpace(m[1])
			if val != "" {
				f.Add(k.Tier, k.Name, val)
			}
		}
	}
}
