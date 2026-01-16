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

var uptimeFact = types.Fact{
	Name: "run.uptime",
	Tier: types.TierOne,
}

// Provides registers run.uptime key.
func (p *ProcUptime) Provides(pr types.Provides) { pr[uptimeFact.Name] = p }

// Requires declares file dependency.
func (p *ProcUptime) Requires(r *types.Requires) { r.Files = append(r.Files, "/proc/uptime") }

// Parse extracts uptime seconds and formats via SecondsToText. Failure yields ParseFailMsg.
func (p *ProcUptime) Parse(f *types.FactDB) {
	data, err := input.FileReader("/proc/uptime")
	if err != nil || data == "" {
		uptimeFact.Value = types.ParseFailMsg
		f.AddFact(uptimeFact)
		return
	}

	fields := strings.Fields(data)
	if len(fields) == 0 {
		uptimeFact.Value = types.ParseFailMsg
		f.AddFact(uptimeFact)
		return
	}

	// Value is a float string; convert to seconds int.
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		uptimeFact.Value = types.ParseFailMsg
		f.AddFact(uptimeFact)
		return
	}

	secs := int64(val)
	uptimeFact.Value = input.SecondsToText(secs)
	f.AddFact(uptimeFact)
}
