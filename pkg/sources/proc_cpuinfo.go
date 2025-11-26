package sources

// Go port of src/clu/sources/proc_cpuinfo.py
// Derives an x86_64 architecture micro-version (x86_64_vN) based on presence of
// progressively demanding CPU flags in /proc/cpuinfo. Skips non-x86_64 architectures.

import (
	"strconv"
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
)

// ProcCpuinfo determines phy.cpu.arch_version for x86_64/amd64 platforms.
type ProcCpuinfo struct{}

// Provides registers the architecture version fact.
func (p *ProcCpuinfo) Provides(pr pkg.Provides) { pr["phy.cpu.arch_version"] = p }

// Requires declares /proc/cpuinfo file dependency.
func (p *ProcCpuinfo) Requires(r *pkg.Requires) { r.Files = append(r.Files, "/proc/cpuinfo") }

// hasAllFlags checks whether every flag in required (space separated) exists in allFlags.
func hasAllFlags(required, allFlags string) bool {
	req := strings.Fields(required)
	all := strings.Fields(allFlags)
	// Build a set for quick membership tests.
	set := make(map[string]struct{}, len(all))
	for _, f := range all {
		set[f] = struct{}{}
	}
	for _, f := range req {
		if _, ok := set[f]; !ok {
			return false
		}
	}
	return true
}

// Parse implements the flag progression logic.
func (p *ProcCpuinfo) Parse(f *pkg.Facts) {
	arch, _ := f.Get("phy.arch")
	if arch != "x86_64" && arch != "amd64" {
		return // skip non-x86
	}

	versions := []string{
		"lm cmov cx8 fpu fxsr mmx syscall sse2",
		"cx16 lahf_lm popcnt sse4_1 sse4_2 ssse3",
		"avx avx2 bmi1 bmi2 f16c fma abm movbe xsave",
		"avx512f avx512bw avx512cd avx512dq avx512vl",
	}

	data, err := pkg.FileReader("/proc/cpuinfo")
	flagsField := ""
	if err == nil && data != "" {
		for _, line := range strings.Split(data, "\n") {
			if strings.HasPrefix(line, "flags") {
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
	f.Set("phy.cpu.arch_version", "x86_64_v"+strconv.Itoa(version))
}
