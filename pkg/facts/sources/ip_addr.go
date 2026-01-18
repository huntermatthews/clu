package sources

// Parses `ip --json addr` to collect interface names, MACs, IPv4 and IPv6 addresses.
// On empty output or non-zero exit code all fact keys are set to ParseFailMsg.

import (
	"encoding/json"
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// IpAddr gathers network interface facts.
type IpAddr struct{}

var ipAddrFacts = map[string]*types.Fact{
	"net.macs": {Name: "net.macs", Tier: types.TierTwo},
	"net.ipv4": {Name: "net.ipv4", Tier: types.TierOne},
	"net.ipv6": {Name: "net.ipv6", Tier: types.TierOne},
	"net.devs": {Name: "net.devs", Tier: types.TierOne},
}

// Provides registers fact keys produced by this source.
func (i *IpAddr) Provides(p types.Provides) {
	for name := range ipAddrFacts {
		p[name] = i
	}
}

// Requires declares external program dependency.
func (i *IpAddr) Requires(r *types.Requires) {
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
func (i *IpAddr) Parse(f *types.FactDB) {
	data, rc, _ := input.CommandRunner("ip --json addr")
	if data == "" || rc != 0 {
		for _, fact := range ipAddrFacts {
			fact.Value = types.ParseFailMsg
			f.AddFact(*fact)
		}
		return
	}

	// Decode JSON.
	var ifaces []ipAddrIface
	if err := json.Unmarshal([]byte(data), &ifaces); err != nil {
		for _, fact := range ipAddrFacts {
			fact.Value = types.ParseFailMsg
			f.AddFact(*fact)
		}
		return
	}

	for _, iface := range ifaces {
		// net.devs
		ipAddrFacts["net.devs"].Value += iface.IfName + " "

		// Addresses
		for _, addr := range iface.AddrInfo {
			switch addr.Family {
			case "inet":
				ipAddrFacts["net.ipv4"].Value += addr.Local + " "
			case "inet6":
				ipAddrFacts["net.ipv6"].Value += addr.Local + " "
			}
		}

		// MAC address (may be empty depending on interface type).
		if strings.TrimSpace(iface.Address) != "" {
			ipAddrFacts["net.macs"].Value += iface.Address + " "
		}
	}

	// Add all facts to the FactDB
	for _, fact := range ipAddrFacts {
		f.AddFact(*fact)
	}
}
