package sources

// Go port of src/clu/sources/lsmem.py
// Parses `lsmem --summary --bytes` to extract total online memory and convert to SI string.

import (
	"strconv"
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// Lsmem collects physical RAM size.
type Lsmem struct{}

var ramFact = types.Fact{
	Name: "phy.ram",
	Tier: types.TierOne,
}

// Provides registers the RAM fact key.
func (l *Lsmem) Provides(p types.Provides) { p[ramFact.Name] = l }

// Requires declares dependency on lsmem command.
func (l *Lsmem) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "lsmem --summary --bytes")
}

// Parse extracts total online memory from command output. On failure sets ParseFailMsg.
func (l *Lsmem) Parse(f *types.FactDB) {
	data, rc, _ := input.CommandRunner("lsmem --summary --bytes")
	if data == "" || rc != 0 {
		ramFact.Value = types.ParseFailMsg
		f.AddFact(ramFact)
		return
	}

	var byteCount string
	for _, line := range strings.Split(data, "\n") {
		if strings.HasPrefix(line, "Total online memory") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				byteCount = strings.TrimSpace(parts[1])
			}
			break
		}
	}

	if byteCount == "" {
		ramFact.Value = types.ParseFailMsg
		f.AddFact(ramFact)
		return
	}

	// Convert using BytesToSI; input expects a float64 number of bytes.
	if v, err := strconv.ParseFloat(byteCount, 64); err == nil {
		ramFact.Value = input.BytesToSI(v)
	} else {
		ramFact.Value = types.ParseFailMsg
	}
	f.AddFact(ramFact)
}
