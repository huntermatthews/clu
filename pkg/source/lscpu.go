package source

// Go port of src/clu/sources/lscpu.py
// Parses `lscpu` output to derive model, vendor, cores, threads, sockets.
// Cores = cores_per_socket * sockets; Threads = threads_per_core * cores.

import (
	"regexp"
	"strconv"
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
)

// Lscpu collects CPU topology and identification facts.
type Lscpu struct{}

var lscpuKeys = []string{
	"phy.cpu.model",
	"phy.cpu.vendor",
	"phy.cpu.cores",
	"phy.cpu.threads",
	"phy.cpu.sockets",
}

// Provides registers the fact keys produced by this source.
func (l *Lscpu) Provides(p pkg.Provides) {
	for _, k := range lscpuKeys {
		p[k] = l
	}
}

// Requires declares program dependency.
func (l *Lscpu) Requires(r *pkg.Requires) { r.Programs = append(r.Programs, "lscpu") }

// Parse executes lscpu and extracts values. On command failure it silently returns
// (matching Python behavior where no facts are set). Computation failures set
// cores/threads to ParseFailMsg.
func (l *Lscpu) Parse(f *pkg.Facts) {
	data, rc := pkg.CommandRunner("lscpu")
	if data == "" || rc != 0 {
		return // mimic Python: skip setting anything
	}
	patterns := map[string]string{
		`^\s*Model name:\s*(.+)`:            "model",
		`^\s*Vendor ID:\s*(.+)`:             "vendor",
		`^\s*Core\(s\) per socket:\s*(\d+)`: "cores_per_socket",
		`^\s*Thread\(s\) per core:\s*(\d+)`: "threads_per_core",
		`^\s*Socket\(s\):\s*(\d+)`:          "sockets",
		`^\s*CPU\(s\):\s*(\d+)`:             "cpus",
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
	cores := ParseFailMsg
	threads := ParseFailMsg
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
	// Set facts (nil/missing -> empty string like Python assignment of None)
	f.Set("phy.cpu.model", fields["model"])   // may be empty
	f.Set("phy.cpu.vendor", fields["vendor"]) // may be empty
	f.Set("phy.cpu.cores", cores)
	f.Set("phy.cpu.threads", threads)
	f.Set("phy.cpu.sockets", fields["sockets"]) // may be empty
}
