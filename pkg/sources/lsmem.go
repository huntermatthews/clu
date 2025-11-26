package sources

// Go port of src/clu/sources/lsmem.py
// Parses `lsmem --summary --bytes` to extract total online memory and convert to SI string.

import (
	"strconv"
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
)

// Lsmem collects physical RAM size.
type Lsmem struct{}

// Provides registers the RAM fact key.
func (l *Lsmem) Provides(p pkg.Provides) { p["phy.ram"] = l }

// Requires declares dependency on lsmem command.
func (l *Lsmem) Requires(r *pkg.Requires) { r.Programs = append(r.Programs, "lsmem --summary --bytes") }

// Parse extracts total online memory from command output. On failure sets ParseFailMsg.
func (l *Lsmem) Parse(f *pkg.Facts) {
	data, rc := pkg.CommandRunner("lsmem --summary --bytes")
	if data == "" || rc != 0 {
		f.Set("phy.ram", ParseFailMsg)
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
		f.Set("phy.ram", ParseFailMsg)
		return
	}
	// Convert using BytesToSI; input expects a float64 number of bytes.
	if v, err := strconv.ParseFloat(byteCount, 64); err == nil {
		f.Set("phy.ram", pkg.BytesToSI(v))
	} else {
		f.Set("phy.ram", ParseFailMsg)
	}
}
