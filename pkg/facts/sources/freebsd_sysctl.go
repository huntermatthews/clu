package sources

// Parses FreeBSD `sysctl` output to extract OS and hardware facts.
// Expected format: lines with "key: value" pairs.

import (
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// FreeBSDSysctl collects OS and hardware facts from sysctl output.
type FreeBSDSysctl struct{}

var freebsdSysctlFacts = map[string]*types.Fact{
	"os.kernel.name":    {Name: "os.kernel.name", Tier: types.TierOne, Origin: "kern.ostype"},
	"os.kernel.version": {Name: "os.kernel.version", Tier: types.TierTwo, Origin: "kern.osrelease"},
	"os.hostname":       {Name: "os.hostname", Tier: types.TierOne, Origin: "kern.hostname"},
	"phy.arch":          {Name: "phy.arch", Tier: types.TierOne, Origin: "hw.machine"},
	"phy.cpu.model":     {Name: "phy.cpu.model", Tier: types.TierTwo, Origin: "hw.model"},
	"phy.cpu.cores":     {Name: "phy.cpu.cores", Tier: types.TierOne, Origin: "hw.ncpu"},
	"phy.ram":           {Name: "phy.ram", Tier: types.TierOne, Origin: "hw.physmem"},
	"run.boot_time":     {Name: "run.boot_time", Tier: types.TierTwo, Origin: "kern.boottime"},
}

// Provides registers fact keys produced by this source.
func (s *FreeBSDSysctl) Provides(p types.Provides) {
	for name := range freebsdSysctlFacts {
		p[name] = s
	}
}

// Requires declares program dependency.
func (s *FreeBSDSysctl) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "sysctl")
}

// Parse executes sysctl and extracts OS/hardware facts from output.
// Maps sysctl keys to fact values using the Origin field.
func (s *FreeBSDSysctl) Parse(f *types.FactDB) {
	// Initialize all facts to ParseFailMsg
	for _, fact := range freebsdSysctlFacts {
		fact.Value = types.ParseFailMsg
	}

	data, rc, _ := input.CommandRunner("sysctl")
	if data == "" || rc != 0 {
		for _, fact := range freebsdSysctlFacts {
			f.AddFact(*fact)
		}
		return
	}

	// Parse sysctl output: "key: value" format
	sysctlMap := make(map[string]string)
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			sysctlMap[key] = value
		}
	}

	// Map sysctl keys to facts using Origin field
	for _, fact := range freebsdSysctlFacts {
		if fact.Origin != "" {
			if v, ok := sysctlMap[fact.Origin]; ok && v != "" {
				// Special handling for RAM to add units
				if fact.Name == "phy.ram" {
					fact.Value = v + " bytes"
				} else {
					fact.Value = v
				}
			}
		}
	}

	// Add all facts to the FactDB
	for _, fact := range freebsdSysctlFacts {
		f.AddFact(*fact)
	}
}
