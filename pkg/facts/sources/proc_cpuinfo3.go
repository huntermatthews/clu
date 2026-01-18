package sources

// ProcCpuinfo3 parses /proc/cpuinfo to derive socket, core and thread counts
// as well as CPU vendor and model. It mirrors the style of other sources
// in this package: model/vendor are set to empty string when missing; cores
// and threads are set to ParseFailMsg when they cannot be computed.

import (
	"strconv"
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// ProcCpuinfo3 provides topology and identification facts from /proc/cpuinfo.
type ProcCpuinfo3 struct{}

var procCpuinfo3Facts = map[string]*types.Fact{
	"phy.cpu.model":   {Name: "phy.cpu.model", Tier: types.TierOne},
	"phy.cpu.vendor":  {Name: "phy.cpu.vendor", Tier: types.TierOne},
	"phy.cpu.cores":   {Name: "phy.cpu.cores", Tier: types.TierOne},
	"phy.cpu.threads": {Name: "phy.cpu.threads", Tier: types.TierOne},
	"phy.cpu.sockets": {Name: "phy.cpu.sockets", Tier: types.TierOne},
}

// Provides registers produced keys.
func (p *ProcCpuinfo3) Provides(pr types.Provides) {
	for name := range procCpuinfo3Facts {
		pr[name] = p
	}
}

// Requires declares file dependency.
func (p *ProcCpuinfo3) Requires(r *types.Requires) {
	r.Files = append(r.Files, "/proc/cpuinfo")
}

// Parse extracts vendor/model and computes sockets/cores/threads.
// Strategy:
//   - Split the file into processor blocks separated by blank lines.
//   - For each block, collect lowercased keys/values (e.g. "physical id", "core id").
//   - Threads = count of processor blocks containing a "processor" entry.
//   - Sockets = number of distinct "physical id" values (if present).
//   - Cores = number of unique (physical id, core id) pairs when available; if not
//     available, attempt to derive from "cpu cores" * sockets; otherwise fall back
//     to ParseFailMsg.
func (p *ProcCpuinfo3) Parse(f *types.FactDB) {
	data, err := input.FileReader("/proc/cpuinfo")
	if err != nil || data == "" {
		// Mirror other parsers: set counts to ParseFailMsg to indicate failure.
		procCpuinfo3Facts["phy.cpu.cores"].Value = types.ParseFailMsg
		procCpuinfo3Facts["phy.cpu.threads"].Value = types.ParseFailMsg
		procCpuinfo3Facts["phy.cpu.sockets"].Value = ""
		procCpuinfo3Facts["phy.cpu.vendor"].Value = ""
		procCpuinfo3Facts["phy.cpu.model"].Value = ""
		for _, fact := range procCpuinfo3Facts {
			f.AddFact(*fact)
		}
		return
	}

	// Normalize line endings and trim trailing whitespace.
	data = strings.ReplaceAll(data, "\r\n", "\n")
	data = strings.TrimSpace(data)
	if data == "" {
		procCpuinfo3Facts["phy.cpu.cores"].Value = types.ParseFailMsg
		procCpuinfo3Facts["phy.cpu.threads"].Value = types.ParseFailMsg
		procCpuinfo3Facts["phy.cpu.sockets"].Value = ""
		procCpuinfo3Facts["phy.cpu.vendor"].Value = ""
		procCpuinfo3Facts["phy.cpu.model"].Value = ""
		for _, fact := range procCpuinfo3Facts {
			f.AddFact(*fact)
		}
		return
	}

	blocks := strings.Split(data, "\n\n")

	threadsCount := 0
	socketSet := make(map[string]struct{})
	corePairSet := make(map[string]struct{})
	coreIDSet := make(map[string]struct{})

	vendor := ""
	model := ""

	// Track cpu cores field from first block if present (cores per socket)
	cpuCoresPerSocket := ""

	for _, block := range blocks {
		if strings.TrimSpace(block) == "" {
			continue
		}
		// parse lines in block
		fields := map[string]string{}
		for _, line := range strings.Split(block, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.ToLower(strings.TrimSpace(parts[0]))
			val := strings.TrimSpace(parts[1])
			fields[key] = val
		}

		if _, ok := fields["processor"]; ok {
			threadsCount++
		}

		// vendor/model pick first non-empty occurrence
		if vendor == "" {
			if v, ok := fields["vendor_id"]; ok {
				vendor = v
			}
		}
		if model == "" {
			if m, ok := fields["model name"]; ok {
				model = m
			}
		}

		if cpuCoresPerSocket == "" {
			if c, ok := fields["cpu cores"]; ok {
				cpuCoresPerSocket = c
			}
		}

		phy := ""
		if v, ok := fields["physical id"]; ok {
			phy = v
			socketSet[phy] = struct{}{}
		}
		core := ""
		if v, ok := fields["core id"]; ok {
			core = v
			coreIDSet[core] = struct{}{}
		}
		if phy != "" && core != "" {
			pair := phy + ":" + core
			corePairSet[pair] = struct{}{}
		}
	}

	// Threads is straightforward
	threads := strconv.Itoa(threadsCount)

	// Sockets: number of distinct physical ids if available
	sockets := ""
	if len(socketSet) > 0 {
		sockets = strconv.Itoa(len(socketSet))
	}

	// Cores: prefer unique physical/core pairs
	cores := types.ParseFailMsg
	if len(corePairSet) > 0 {
		cores = strconv.Itoa(len(corePairSet))
	} else if cpuCoresPerSocket != "" && sockets != "" {
		if cpsInt, err1 := strconv.Atoi(cpuCoresPerSocket); err1 == nil {
			if sockInt, err2 := strconv.Atoi(sockets); err2 == nil {
				cores = strconv.Itoa(cpsInt * sockInt)
			}
		}
	} else if len(coreIDSet) > 0 {
		cores = strconv.Itoa(len(coreIDSet))
	}

	// If sockets could not be determined but we have cores and threads, try to compute sockets
	if sockets == "" && cores != types.ParseFailMsg {
		if coresInt, err1 := strconv.Atoi(cores); err1 == nil {
			if threadsCount > 0 {
				// threads = logical CPUs; sockets = threads / (cores_per_socket) unknown
				// cannot reliably derive sockets from cores+threads without more info, leave empty
				_ = coresInt
			}
		}
	}

	// Set fact values (model/vendor may be empty string like lscpu behavior)
	procCpuinfo3Facts["phy.cpu.model"].Value = model
	procCpuinfo3Facts["phy.cpu.vendor"].Value = vendor
	procCpuinfo3Facts["phy.cpu.cores"].Value = cores
	procCpuinfo3Facts["phy.cpu.threads"].Value = threads
	procCpuinfo3Facts["phy.cpu.sockets"].Value = sockets

	// Add all facts to the FactDB
	for _, fact := range procCpuinfo3Facts {
		f.AddFact(*fact)
	}
}
