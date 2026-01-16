package sources

// Go port of src/clu/sources/udevadm_ram.py
// Parses `udevadm info --path /devices/virtual/dmi/id` to aggregate MEMORY_DEVICE_*_SIZE values.

import (
	"regexp"
	"strconv"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

// UdevadmRam collects total physical RAM size from udevadm output.
type UdevadmRam struct{}

var ramFactUdevadm = types.Fact{
	Name: "phy.ram",
	Tier: types.TierOne,
}

// Provides registers phy.ram key.
func (u *UdevadmRam) Provides(p types.Provides) { p[ramFactUdevadm.Name] = u }

// Requires declares dependency on udevadm command.
func (u *UdevadmRam) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "udevadm info --path /devices/virtual/dmi/id")
}

// Parse executes command, extracts MEMORY_DEVICE_x_SIZE numbers, sums them and converts to SI string.
// On failure sets ParseFailMsg.
func (u *UdevadmRam) Parse(f *types.FactDB) {
	data, rc, _ := input.CommandRunner("udevadm info --path /devices/virtual/dmi/id")
	if data == "" || rc != 0 {
		ramFactUdevadm.Value = types.ParseFailMsg
		f.AddFact(ramFactUdevadm)
		return
	}
	re := regexp.MustCompile(`MEMORY_DEVICE_\d+_SIZE=(\d+)`)
	matches := re.FindAllStringSubmatch(data, -1)
	var total float64
	for _, m := range matches {
		if len(m) == 2 {
			if v, err := strconv.ParseFloat(m[1], 64); err == nil {
				total += v
			}
		}
	}
	// If no matches, total remains 0 -> BytesToSI will format "0.0 B" (consistent with Python result for sum([])=0).
	ramFactUdevadm.Value = input.BytesToSI(total)
	f.AddFact(ramFactUdevadm)
}
