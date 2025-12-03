package source

// Go port of src/clu/sources/ip_addr.py
// Parses `ip --json addr` to collect interface names, MACs, IPv4 and IPv6 addresses.
// On empty output or non-zero exit code all fact keys are set to ParseFailMsg.

import (
	"encoding/json"
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
)

// IpAddr gathers network interface facts.
type IpAddr struct{}

// ipAddrKeys mirrors Python primary_keys list order.
var ipAddrKeys = []string{"net.macs", "net.ipv4", "net.ipv6", "net.devs"}

// Provides registers fact keys produced by this source.
func (i *IpAddr) Provides(p pkg.Provides) {
	for _, k := range ipAddrKeys {
		p[k] = i
	}
}

// Requires declares external program dependency.
func (i *IpAddr) Requires(r *pkg.Requires) {
	r.Programs = append(r.Programs, "ip --json addr")
}

// ipAddrIface models a subset of JSON returned by `ip --json addr`.
type ipAddrIface struct {
	IfName   string `json:"ifname"`
	Address  string `json:"address"`
	AddrInfo []struct {
		Family string `json:"family"`
		Local  string `json:"local"`
	} `json:"addr_info"`
}

// Parse executes command, decodes JSON and accumulates space-separated lists (trailing space retained for parity).
func (i *IpAddr) Parse(f *pkg.Facts) {
	data, rc := pkg.CommandRunner("ip --json addr")
	if data == "" || rc != 0 {
		for _, k := range ipAddrKeys {
			f.Set(k, ParseFailMsg)
		}
		return
	}
	// Initialize keys to empty strings (Python behavior before accumulation).
	for _, k := range ipAddrKeys {
		f.Set(k, "")
	}
	// Decode JSON.
	var ifaces []ipAddrIface
	if err := json.Unmarshal([]byte(data), &ifaces); err != nil {
		for _, k := range ipAddrKeys {
			f.Set(k, ParseFailMsg)
		}
		return
	}
	for _, iface := range ifaces {
		// net.devs
		cur, _ := f.Get("net.devs")
		f.Set("net.devs", cur+iface.IfName+" ")
		// Addresses
		for _, addr := range iface.AddrInfo {
			if addr.Family == "inet" {
				cur4, _ := f.Get("net.ipv4")
				f.Set("net.ipv4", cur4+addr.Local+" ")
			} else if addr.Family == "inet6" {
				cur6, _ := f.Get("net.ipv6")
				f.Set("net.ipv6", cur6+addr.Local+" ")
			}
		}
		// MAC address (may be empty depending on interface type).
		if strings.TrimSpace(iface.Address) != "" {
			curMac, _ := f.Get("net.macs")
			f.Set("net.macs", curMac+iface.Address+" ")
		}
	}
}
