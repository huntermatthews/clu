package sources

// Go port of src/clu/sources/ipmitool.py
// Collects BMC/IPMI facts via `ipmitool mc info` and `ipmitool lan print`.
// Skips parsing unless platform is physical (facts["phy.platform"] == "physical").

import (
	"regexp"
	"strings"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/facts/types"
)

// Ipmitool gathers BMC firmware/network details.
type Ipmitool struct{}

var ipmitoolLanKeys = []string{
	"bmc.ipv4_source",
	"bmc.ipv4_address",
	"bmc.ipv4_mask",
	"bmc.mac_address",
}

var ipmitoolMcInfoKeys = []string{
	"bmc.firmware_version",
	"bmc.manufacturer_id",
	"bmc.manufacturer_name",
}

// Provides registers all fact keys produced by this source.
func (i *Ipmitool) Provides(p types.Provides) {
	for _, k := range ipmitoolLanKeys {
		p[k] = i
	}
	for _, k := range ipmitoolMcInfoKeys {
		p[k] = i
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
	patterns := map[string]string{
		`(?m)^IP Address Source *: (.+)`: "bmc.ipv4_source",
		`(?m)^IP Address *: (.+)`:        "bmc.ipv4_address",
		`(?m)^Subnet Mask *: (.+)`:       "bmc.ipv4_mask",
		`(?m)^MAC Address *: (.+)`:       "bmc.mac_address",
	}
	data, rc := pkg.CommandRunner("ipmitool lan print")
	if data == "" || rc != 0 {
		for _, k := range ipmitoolLanKeys {
			f.Set(k, types.ParseFailMsg)
		}
		return
	}
	// Extract values
	for regex, key := range patterns {
		re := regexp.MustCompile(regex)
		m := re.FindStringSubmatch(data)
		if len(m) == 2 {
			val := strings.TrimSpace(m[1])
			if val != "" {
				f.Set(key, val)
			}
		}
	}
}

func (i *Ipmitool) parseMcInfo(f *types.Facts) {
	patterns := map[string]string{
		`(?m)^Firmware Revision *: (.+)`: "bmc.firmware_version",
		`(?m)^Manufacturer ID *: (.+)`:   "bmc.manufacturer_id",
		`(?m)^Manufacturer Name *: (.+)`: "bmc.manufacturer_name",
	}
	data, rc := pkg.CommandRunner("ipmitool mc info")
	if data == "" || rc != 0 {
		for _, k := range ipmitoolMcInfoKeys {
			f.Set(k, types.ParseFailMsg)
		}
		return
	}
	for regex, key := range patterns {
		re := regexp.MustCompile(regex)
		m := re.FindStringSubmatch(data)
		if len(m) == 2 {
			val := strings.TrimSpace(m[1])
			if val != "" {
				f.Set(key, val)
			}
		}
	}
}
