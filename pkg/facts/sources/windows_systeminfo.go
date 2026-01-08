package sources

// Parses Windows `systeminfo` output in CSV format to extract OS and hardware facts.
// Expected format: header row (CSV with quoted fields) followed by data row.

import (
	"encoding/csv"
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// WindowsSysteminfo collects OS and hardware facts from systeminfo output.
type WindowsSysteminfo struct{}

var systemInfoKeys = []string{
	"os.hostname",
	"os.name",
	"os.version",
	"os.build",
	"os.owner",
	"os.organization",
	"os.install_date",
	"run.boot_time",
	"phy.manufacturer",
	"phy.model",
	"phy.system_type",
	"phy.bios_version",
	"net.domain",
	"net.network_cards",
	"phy.cpu.cores",
	"phy.ram",
}

// Provides registers fact keys produced by this source.
func (w *WindowsSysteminfo) Provides(p types.Provides) {
	for _, k := range systemInfoKeys {
		p[k] = w
	}
}

// Requires declares program dependency.
func (w *WindowsSysteminfo) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "systeminfo")
}

// Parse executes systeminfo and extracts OS/hardware facts from CSV output.
// Maps CSV header to values and extracts relevant fields.
func (w *WindowsSysteminfo) Parse(f *types.Facts) {

	data, rc, _ := input.CommandRunner("systeminfo")
	if data == "" || rc != 0 {
		for _, k := range systemInfoKeys {
			f.Add(types.TierOne, k, types.ParseFailMsg)
		}
		return
	}

	// Parse CSV: expect header + data row
	r := csv.NewReader(strings.NewReader(data))
	records, err := r.ReadAll()
	if err != nil || len(records) < 2 {
		for _, k := range systemInfoKeys {
			f.Add(types.TierOne, k, types.ParseFailMsg)
		}
		return
	}

	header := records[0]
	dataRow := records[1]

	// Build field map: header -> value
	fields := make(map[string]string)
	for i, h := range header {
		if i < len(dataRow) {
			fields[h] = dataRow[i]
		}
	}

	// Extract and set facts
	if v, ok := fields["Host Name"]; ok && v != "" {
		f.Add(types.TierOne, "os.hostname", v)
	} else {
		f.Add(types.TierOne, "os.hostname", types.ParseFailMsg)
	}

	if v, ok := fields["OS Name"]; ok && v != "" {
		f.Add(types.TierOne, "os.name", v)
	} else {
		f.Add(types.TierOne, "os.name", types.ParseFailMsg)
	}

	if v, ok := fields["OS Version"]; ok && v != "" {
		f.Add(types.TierOne, "os.version", v)
	} else {
		f.Add(types.TierOne, "os.version", types.ParseFailMsg)
	}

	if v, ok := fields["OS Build Type"]; ok && v != "" {
		f.Add(types.TierOne, "os.build", v)
	} else {
		f.Add(types.TierOne, "os.build", types.ParseFailMsg)
	}

	// Total Physical Memory (in MB, convert to bytes for consistency)
	if v, ok := fields["Total Physical Memory"]; ok && v != "" {
		// Extract numeric part and format as SI string
		cleaned := strings.Fields(v)
		if len(cleaned) > 0 {
			f.Add(types.TierOne, "phy.ram", cleaned[0]+" MB")
		} else {
			f.Add(types.TierOne, "phy.ram", types.ParseFailMsg)
		}
	} else {
		f.Add(types.TierOne, "phy.ram", types.ParseFailMsg)
	}

	// Processor count (parse "N Processor(s) Installed" format)
	if v, ok := fields["Processor(s)"]; ok && v != "" {
		// Extract first field (count)
		parts := strings.Fields(v)
		if len(parts) > 0 {
			f.Add(types.TierOne, "phy.cpu.cores", parts[0])
		} else {
			f.Add(types.TierOne, "phy.cpu.cores", types.ParseFailMsg)
		}
	} else {
		f.Add(types.TierOne, "phy.cpu.cores", types.ParseFailMsg)
	}

	// Registered Owner
	if v, ok := fields["Registered Owner"]; ok && v != "" {
		f.Add(types.TierTwo, "os.owner", v)
	} else {
		f.Add(types.TierTwo, "os.owner", types.ParseFailMsg)
	}

	// Registered Organization
	if v, ok := fields["Registered Organization"]; ok && v != "" {
		f.Add(types.TierTwo, "os.organization", v)
	} else {
		f.Add(types.TierTwo, "os.organization", types.ParseFailMsg)
	}

	// Original Install Date
	if v, ok := fields["Original Install Date"]; ok && v != "" {
		f.Add(types.TierTwo, "os.install_date", v)
	} else {
		f.Add(types.TierTwo, "os.install_date", types.ParseFailMsg)
	}

	// System Boot Time
	if v, ok := fields["System Boot Time"]; ok && v != "" {
		f.Add(types.TierTwo, "run.boot_time", v)
	} else {
		f.Add(types.TierTwo, "run.boot_time", types.ParseFailMsg)
	}

	// System Manufacturer
	if v, ok := fields["System Manufacturer"]; ok && v != "" {
		f.Add(types.TierTwo, "phy.manufacturer", v)
	} else {
		f.Add(types.TierTwo, "phy.manufacturer", types.ParseFailMsg)
	}

	// System Model
	if v, ok := fields["System Model"]; ok && v != "" {
		f.Add(types.TierTwo, "phy.model", v)
	} else {
		f.Add(types.TierTwo, "phy.model", types.ParseFailMsg)
	}

	// System Type
	if v, ok := fields["System Type"]; ok && v != "" {
		f.Add(types.TierTwo, "phy.system_type", v)
	} else {
		f.Add(types.TierTwo, "phy.system_type", types.ParseFailMsg)
	}

	// BIOS Version (version, date format - extract version only)
	if v, ok := fields["BIOS Version"]; ok && v != "" {
		// Split on comma and take first part (version)
		parts := strings.SplitN(v, ",", 2)
		version := strings.TrimSpace(parts[0])
		if version != "" {
			f.Add(types.TierThree, "phy.bios_version", version)
		} else {
			f.Add(types.TierThree, "phy.bios_version", types.ParseFailMsg)
		}
	} else {
		f.Add(types.TierThree, "phy.bios_version", types.ParseFailMsg)
	}

	// Domain
	if v, ok := fields["Domain"]; ok && v != "" {
		f.Add(types.TierTwo, "net.domain", v)
	} else {
		f.Add(types.TierTwo, "net.domain", types.ParseFailMsg)
	}

	// Network Cards
	if v, ok := fields["Network Card(s)"]; ok && v != "" {
		f.Add(types.TierThree, "net.network_cards", v)
	} else {
		f.Add(types.TierThree, "net.network_cards", types.ParseFailMsg)
	}
}
