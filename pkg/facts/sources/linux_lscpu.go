package sources

// Go port of src/clu/sources/lscpu.py
// Parses `lscpu` output to derive model, vendor, cores, threads, sockets.
// Cores = cores_per_socket * sockets; Threads = threads_per_core * cores.

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// Lscpu collects CPU topology and identification facts.
type Lscpu struct{}

var lscpuFacts = map[string]*types.Fact{
	"phy.cpu.model":   {Name: "phy.cpu.model", Tier: types.TierOne},
	"phy.cpu.vendor":  {Name: "phy.cpu.vendor", Tier: types.TierTwo},
	"phy.cpu.cores":   {Name: "phy.cpu.cores", Tier: types.TierOne},
	"phy.cpu.threads": {Name: "phy.cpu.threads", Tier: types.TierTwo},
	"phy.cpu.sockets": {Name: "phy.cpu.sockets", Tier: types.TierThree},
}

// Provides registers the fact keys produced by this source.
func (l *Lscpu) Provides(p types.Provides) {
	for name := range lscpuFacts {
		p[name] = l
	}
}

// Requires declares program dependency.
func (l *Lscpu) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "lscpu")
}

// Parse executes lscpu and extracts values. On command failure it silently returns
// (matching Python behavior where no facts are set). Computation failures set
// cores/threads to ParseFailMsg.
func (l *Lscpu) Parse(f *types.FactDB) {
	data, rc, _ := input.CommandRunner("lscpu")
	if data == "" || rc != 0 {
		return
	}

	// Normalize escaped newlines ("\n") into real newlines so both raw and
	// literal-escaped inputs are handled consistently.
	data = strings.ReplaceAll(data, "\\n", "\n")

	patterns := map[string]string{
		`(?m)^\s*Model name:\s*(.+)`:            "model",
		`(?m)^\s*Vendor ID:\s*(.+)`:             "vendor",
		`(?m)^\s*Core\(s\) per socket:\s*(\d+)`: "cores_per_socket",
		`(?m)^\s*Thread\(s\) per core:\s*(\d+)`: "threads_per_core",
		`(?m)^\s*Socket\(s\):\s*(\d+)`:          "sockets",
		`(?m)^\s*CPU\(s\):\s*(\d+)`:             "cpus",
	}
	fields := map[string]string{}

	for regex, key := range patterns {
		re := regexp.MustCompile(regex)
		m := re.FindStringSubmatch(data)
		if len(m) == 2 {
			fields[key] = strings.TrimSpace(m[1])
		}
	}

	// Derive cores & threads if possible
	cores := types.ParseFailMsg
	threads := types.ParseFailMsg
	if cps, ok1 := fields["cores_per_socket"]; ok1 {
		if sockets, ok2 := fields["sockets"]; ok2 {
			if cpsInt, err1 := strconv.Atoi(cps); err1 == nil {
				if socketsInt, err2 := strconv.Atoi(sockets); err2 == nil {
					totalCores := cpsInt * socketsInt
					cores = strconv.Itoa(totalCores)
					if tpc, ok3 := fields["threads_per_core"]; ok3 {
						if tpcInt, err3 := strconv.Atoi(tpc); err3 == nil {
							threads = strconv.Itoa(tpcInt * totalCores)
						}
					}
				}
			}
		}
	}

	// Set fact values
	lscpuFacts["phy.cpu.model"].Value = fields["model"]   // may be empty
	lscpuFacts["phy.cpu.vendor"].Value = fields["vendor"] // may be empty
	lscpuFacts["phy.cpu.cores"].Value = cores
	lscpuFacts["phy.cpu.threads"].Value = threads
	lscpuFacts["phy.cpu.sockets"].Value = fields["sockets"] // may be empty

	// Add all facts to the FactDB
	for _, fact := range lscpuFacts {
		f.AddFact(*fact)
	}
}
