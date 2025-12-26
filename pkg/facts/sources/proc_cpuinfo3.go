package sources

// ProcCpuinfo3 parses /proc/cpuinfo to derive socket, core and thread counts
// as well as CPU vendor and model. It mirrors the style of other sources
// in this package: model/vendor are set to empty string when missing; cores
// and threads are set to ParseFailMsg when they cannot be computed.

import (
	"strconv"
	"strings"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/facts/types"
)

// ProcCpuinfo3 provides topology and identification facts from /proc/cpuinfo.
type ProcCpuinfo3 struct{}

var procCpuinfo3Keys = []string{
	"phy.cpu.model",
	"phy.cpu.vendor",
	"phy.cpu.cores",
	"phy.cpu.threads",
	"phy.cpu.sockets",
}

// Provides registers produced keys.
func (p *ProcCpuinfo3) Provides(pr types.Provides) {
	for _, k := range procCpuinfo3Keys {
		pr[k] = p
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
func (p *ProcCpuinfo3) Parse(f *types.Facts) {
	data, err := pkg.FileReader("/proc/cpuinfo")
	if err != nil || data == "" {
		// Mirror other parsers: set counts to ParseFailMsg to indicate failure.
		f.Set("phy.cpu.cores", types.ParseFailMsg)
		f.Set("phy.cpu.threads", types.ParseFailMsg)
		f.Set("phy.cpu.sockets", "")
		f.Set("phy.cpu.vendor", "")
		f.Set("phy.cpu.model", "")
		return
	}

	// Normalize line endings and trim trailing whitespace.
	data = strings.ReplaceAll(data, "\r\n", "\n")
	data = strings.TrimSpace(data)
	if data == "" {
		f.Set("phy.cpu.cores", types.ParseFailMsg)
		f.Set("phy.cpu.threads", types.ParseFailMsg)
		f.Set("phy.cpu.sockets", "")
		f.Set("phy.cpu.vendor", "")
		f.Set("phy.cpu.model", "")
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

	// Set facts (model/vendor may be empty string like lscpu behavior)
	f.Set("phy.cpu.model", model)
	f.Set("phy.cpu.vendor", vendor)
	f.Set("phy.cpu.cores", cores)
	f.Set("phy.cpu.threads", threads)
	f.Set("phy.cpu.sockets", sockets)
}
