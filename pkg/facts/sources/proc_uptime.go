package sources

// Reads /proc/uptime and converts the first float (seconds since boot) to human text.

import (
	"strconv"
	"strings"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/facts/types"
)

// ProcUptime parses /proc/uptime for run.uptime fact.
type ProcUptime struct{}

// Provides registers run.uptime key.
func (p *ProcUptime) Provides(pr types.Provides) { pr["run.uptime"] = p }

// Requires declares file dependency.
func (p *ProcUptime) Requires(r *types.Requires) { r.Files = append(r.Files, "/proc/uptime") }

// Parse extracts uptime seconds and formats via SecondsToText. Failure yields ParseFailMsg.
func (p *ProcUptime) Parse(f *types.Facts) {
	data, err := pkg.FileReader("/proc/uptime")
	if err != nil || data == "" {
		f.Set("run.uptime", types.ParseFailMsg)
		return
	}

	fields := strings.Fields(data)
	if len(fields) == 0 {
		f.Set("run.uptime", types.ParseFailMsg)
		return
	}

	// Value is a float string; convert to seconds int.
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		f.Set("run.uptime", types.ParseFailMsg)
		return
	}

	secs := int64(val)
	f.Set("run.uptime", pkg.SecondsToText(secs))
}
