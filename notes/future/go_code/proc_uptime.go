package sources

// Go port of src/clu/sources/proc_uptime.py
// Reads /proc/uptime and converts the first float (seconds since boot) to human text.

import (
	"strconv"
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
	facts "github.com/huntermatthews/clu/pkg/facts"
)

// ProcUptime parses /proc/uptime for run.uptime fact.
type ProcUptime struct{}

// Provides registers run.uptime key.
func (p *ProcUptime) Provides(pr facts.Provides) { pr["run.uptime"] = p }

// Requires declares file dependency.
func (p *ProcUptime) Requires(r *facts.Requires) { r.Files = append(r.Files, "/proc/uptime") }

// Parse extracts uptime seconds and formats via SecondsToText. Failure yields ParseFailMsg.
func (p *ProcUptime) Parse(f *facts.Facts) {
	data, err := pkg.FileReader("/proc/uptime")
	if err != nil || data == "" {
		f.Set("run.uptime", ParseFailMsg)
		return
	}
	fields := strings.Fields(data)
	if len(fields) == 0 {
		f.Set("run.uptime", ParseFailMsg)
		return
	}
	// Value is a float string; convert to seconds int.
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		f.Set("run.uptime", ParseFailMsg)
		return
	}
	secs := int64(val)
	f.Set("run.uptime", pkg.SecondsToText(secs))
}
