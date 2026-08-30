// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package sources

// Linux memory fact collection. Memory tries supported collection strategies in
// preference order so a missing or malformed primary tool does not prevent a
// usable phy.ram fact from being reported.

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/huntermatthews/clu/internal/facts/types"
	"github.com/huntermatthews/clu/internal/input"
)

const (
	memoryFactName = "phy.ram"
	lsmemCommand   = "lsmem --summary --bytes"
	udevadmCommand = "udevadm info --path /devices/virtual/dmi/id"
)

var (
	memoryFact = types.Fact{
		Name: memoryFactName,
		Tier: types.TierOne,
	}
	memorySources = []memoryCollector{
		&lsmem{},
		&udevadmRam{},
		// Append ProcMeminfo here when it is implemented.
	}
)

// memoryCollector is an internal strategy for deriving phy.ram. Unlike a
// fact source, it has no independent Provides or Requires contract.
type memoryCollector interface {
	Parse(*types.FactDB)
}

// Memory collects physical RAM using the first successful Linux collector.
type Memory struct{}

// Provides registers the sole phy.ram provider. The individual collectors are
// implementation details and must not be independently selected.
func (m *Memory) Provides(p types.Provides) { p[memoryFact.Name] = m }

// Requires declares every collector command so the collector subcommand can
// capture all candidate outputs. Fact gathering still uses them in preference
// order and stops after the first successful result.
func (m *Memory) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, lsmemCommand, udevadmCommand)
}

// Parse tries memory collectors in preference order. A collector succeeds only
// when it stores a non-empty, non-error phy.ram value.
func (m *Memory) Parse(f *types.FactDB) {
	for _, source := range memorySources {
		source.Parse(f)
		if value, ok := f.Get(memoryFactName); ok && value != "" && value != types.ParseFailMsg {
			return
		}
	}
}

// lsmem parses `lsmem --summary --bytes` to extract total online memory.
type lsmem struct{}

// Parse extracts total online memory from command output. On failure it stores
// ParseFailMsg so Memory can try the next collector.
func (l *lsmem) Parse(f *types.FactDB) {
	memoryFact.Value = types.ParseFailMsg
	defer func() { f.AddFact(memoryFact) }()

	data, rc, _ := input.CommandRunner(lsmemCommand)
	if data == "" || rc != 0 {
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
		return
	}

	bytes, err := strconv.ParseFloat(byteCount, 64)
	if err != nil || bytes <= 0 {
		return
	}
	memoryFact.Value = input.BytesToSI(bytes)
}

// udevadmRam parses DMI memory-device sizes from
// `udevadm info --path /devices/virtual/dmi/id`.
type udevadmRam struct{}

// Parse sums non-zero MEMORY_DEVICE_*_SIZE values. Missing or malformed output
// is a parse failure rather than a valid zero-byte RAM result.
func (u *udevadmRam) Parse(f *types.FactDB) {
	memoryFact.Value = types.ParseFailMsg
	defer func() { f.AddFact(memoryFact) }()

	data, rc, _ := input.CommandRunner(udevadmCommand)
	if data == "" || rc != 0 {
		return
	}

	re := regexp.MustCompile(`MEMORY_DEVICE_\d+_SIZE=(\d+)`)
	matches := re.FindAllStringSubmatch(data, -1)
	var total float64
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		bytes, err := strconv.ParseFloat(match[1], 64)
		if err == nil {
			total += bytes
		}
	}

	if total <= 0 {
		return
	}
	memoryFact.Value = input.BytesToSI(total)
}
