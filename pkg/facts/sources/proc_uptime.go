package sources

// Reads /proc/uptime and converts the first float (seconds since boot) to human text.

import (
	"strconv"
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

// ProcUptime parses /proc/uptime for run.uptime fact.
type ProcUptime struct{}

// Provides registers run.uptime key.
func (p *ProcUptime) Provides(pr types.Provides) { pr["run.uptime"] = p }

// Requires declares file dependency.
func (p *ProcUptime) Requires(r *types.Requires) { r.Files = append(r.Files, "/proc/uptime") }

// Parse extracts uptime seconds and formats via SecondsToText. Failure yields ParseFailMsg.
func (p *ProcUptime) Parse(f *types.FactDB) {
	data, err := input.FileReader("/proc/uptime")
	if err != nil || data == "" {
		f.Add(types.TierOne, "run.uptime", types.ParseFailMsg)
		return
	}

	fields := strings.Fields(data)
	if len(fields) == 0 {
		f.Add(types.TierOne, "run.uptime", types.ParseFailMsg)
		return
	}

	// Value is a float string; convert to seconds int.
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		f.Add(types.TierOne, "run.uptime", types.ParseFailMsg)
		return
	}

	secs := int64(val)
	f.Add(types.TierOne, "run.uptime", input.SecondsToText(secs))
}
