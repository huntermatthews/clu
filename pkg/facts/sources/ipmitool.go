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

type factThing struct {
	name  string
	value string
	regex string
	tier  types.Tier
}

var ipmitoolLanFacts = []factThing{
	{"bmc.ipv4_source", "", `(?m)^IP Address Source *: (.+)`, types.TierTwo},
	{"bmc.ipv4_address", "", `(?m)^IP Address *: (.+)`, types.TierOne},
	{"bmc.ipv4_mask", "", `(?m)^Subnet Mask *: (.+)`, types.TierTwo},
	{"bmc.mac_address", "", `(?m)^MAC Address *: (.+)`, types.TierOne},
}
var ipmitoolMcInfoFacts = []factThing{
	{"bmc.firmware_version", "", `(?m)^Firmware Revision *: (.+)`, types.TierOne},
	{"bmc.manufacturer_id", "", `(?m)^Manufacturer ID *: (.+)`, types.TierTwo},
	{"bmc.manufacturer_name", "", `(?m)^Manufacturer Name *: (.+)`, types.TierTwo},
}

// Provides registers all fact keys produced by this source.
func (i *Ipmitool) Provides(p types.Provides) {
	for _, k := range ipmitoolLanFacts {
		p[k.name] = i
	}
	for _, k := range ipmitoolMcInfoFacts {
		p[k.name] = i
	}
}

// Requires declares external program dependency.
func (i *Ipmitool) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "ipmitool")
}

// Parse orchestrates the two ipmitool queries if platform is physical.
func (i *Ipmitool) Parse(f *types.Facts) {
	platform, _ := f.Get("phy.platform")
	if platform != "physical" { // skip if not physical hardware
		return
	}
	i.parseMcInfo(f)
	i.parseLanPrint(f)
}

func (i *Ipmitool) parseLanPrint(f *types.Facts) {

	data, rc, _ := input.CommandRunner("ipmitool lan print")
	if data == "" || rc != 0 {
		for _, k := range ipmitoolLanFacts {
			f.Add(k.tier, k.name, types.ParseFailMsg)
		}
		return
	}
	// Extract values
	for _, k := range ipmitoolLanFacts {
		re := regexp.MustCompile(k.regex)
		m := re.FindStringSubmatch(data)
		if len(m) == 2 {
			val := strings.TrimSpace(m[1])
			if val != "" {
				f.Add(k.tier, k.name, val)
			}
		}
	}
}

func (i *Ipmitool) parseMcInfo(f *types.Facts) {

	data, rc, _ := input.CommandRunner("ipmitool mc info")
	if data == "" || rc != 0 {
		for _, k := range ipmitoolMcInfoFacts {
			f.Add(k.tier, k.name, types.ParseFailMsg)
		}
		return
	}
	for _, k := range ipmitoolMcInfoFacts {
		re := regexp.MustCompile(k.regex)
		m := re.FindStringSubmatch(data)
		if len(m) == 2 {
			val := strings.TrimSpace(m[1])
			if val != "" {
				f.Add(k.tier, k.name, val)
			}
		}
	}
}
