package source

// Go port of src/clu/sources/ipmitool.py
// Collects BMC/IPMI facts via `ipmitool mc info` and `ipmitool lan print`.
// Skips parsing unless platform is physical (facts["phy.platform"] == "physical").

import (
	"regexp"
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
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
func (i *Ipmitool) Provides(p pkg.Provides) {
	for _, k := range ipmitoolLanKeys {
		p[k] = i
	}
	for _, k := range ipmitoolMcInfoKeys {
		p[k] = i
	}
}

// Requires declares external program dependency.
func (i *Ipmitool) Requires(r *pkg.Requires) {
	r.Programs = append(r.Programs, "ipmitool")
}

// Parse orchestrates the two ipmitool queries if platform is physical.
func (i *Ipmitool) Parse(f *pkg.Facts) {
	platform, _ := f.Get("phy.platform")
	if platform != "physical" { // skip if not physical hardware
		return
	}
	i.parseMcInfo(f)
	i.parseLanPrint(f)
}

func (i *Ipmitool) parseLanPrint(f *pkg.Facts) {
	patterns := map[string]string{
		`^IP Address Source *: (.+)`: "bmc.ipv4_source",
		`^IP Address *: (.+)`:        "bmc.ipv4_address",
		`^Subnet Mask *: (.+)`:       "bmc.ipv4_mask",
		`^MAC Address *: (.+)`:       "bmc.mac_address",
	}
	data, rc := pkg.CommandRunner("ipmitool lan print")
	if data == "" || rc != 0 {
		for _, k := range ipmitoolLanKeys {
			f.Set(k, ParseFailMsg)
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

func (i *Ipmitool) parseMcInfo(f *pkg.Facts) {
	patterns := map[string]string{
		`^Firmware Revision *: (.+)`: "bmc.firmware_version",
		`^Manufacturer ID *: (.+)`:   "bmc.manufacturer_id",
		`^Manufacturer Name *: (.+)`: "bmc.manufacturer_name",
	}
	data, rc := pkg.CommandRunner("ipmitool mc info")
	if data == "" || rc != 0 {
		for _, k := range ipmitoolMcInfoKeys {
			f.Set(k, ParseFailMsg)
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
