package sources

// Alternate ProcCpuinfo parser that implements the same Sources interface.
// This variant is tolerant of slightly different /proc/cpuinfo layouts
// (accepts both "flags" and "Features" keys) but otherwise derives
// the same `phy.cpu.arch_version` value as the original implementation.

import (
	"strconv"
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

// ProcCpuinfo2 is a drop-in alternate parser for /proc/cpuinfo.
type ProcCpuinfo2 struct{}

// Provides registers the architecture version fact.
func (p *ProcCpuinfo2) Provides(pr types.Provides) {
	pr["phy.cpu.arch_version"] = p
}

// Requires declares /proc/cpuinfo file dependency.
func (p *ProcCpuinfo2) Requires(r *types.Requires) {
	r.Files = append(r.Files, "/proc/cpuinfo")
}

// Parse implements the same flag-progression logic as ProcCpuinfo but
// tolerates alternative field names.
func (p *ProcCpuinfo2) Parse(f *types.FactDB) {
	arch, _ := f.Get("phy.arch")
	if arch != "x86_64" && arch != "amd64" {
		return
	}

	versions := []string{
		"lm cmov cx8 fpu fxsr mmx syscall sse2",
		"cx16 lahf_lm popcnt sse4_1 sse4_2 ssse3",
		"avx avx2 bmi1 bmi2 f16c fma abm movbe xsave",
		"avx512f avx512bw avx512cd avx512dq avx512vl",
	}

	data, err := input.FileReader("/proc/cpuinfo")
	flagsField := ""
	if err == nil && data != "" {
		for _, line := range strings.Split(data, "\n") {
			// Accept either "flags" or "Features" (some ARM/uarch outputs use Features).
			if strings.HasPrefix(line, "flags") || strings.HasPrefix(line, "Features") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					flagsField = strings.TrimSpace(parts[1])
				}
				break
			}
		}
	}

	version := 0
	for _, group := range versions {
		if hasAllFlags(group, flagsField) {
			version++
		} else {
			break
		}
	}
	f.Add(types.TierOne, "phy.cpu.arch_version", "x86_64_v"+strconv.Itoa(version))
}
