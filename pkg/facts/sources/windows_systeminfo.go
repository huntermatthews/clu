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

var systemInfoFacts = map[string]*types.Fact{
	"os.hostname":       {Name: "os.hostname", Tier: types.TierOne},
	"os.name":           {Name: "os.name", Tier: types.TierOne},
	"os.version":        {Name: "os.version", Tier: types.TierOne},
	"os.build":          {Name: "os.build", Tier: types.TierOne},
	"os.owner":          {Name: "os.owner", Tier: types.TierTwo},
	"os.organization":   {Name: "os.organization", Tier: types.TierTwo},
	"os.install_date":   {Name: "os.install_date", Tier: types.TierTwo},
	"run.boot_time":     {Name: "run.boot_time", Tier: types.TierTwo},
	"phy.manufacturer":  {Name: "phy.manufacturer", Tier: types.TierTwo},
	"phy.model":         {Name: "phy.model", Tier: types.TierTwo},
	"phy.system_type":   {Name: "phy.system_type", Tier: types.TierTwo},
	"phy.bios_version":  {Name: "phy.bios_version", Tier: types.TierThree},
	"net.domain":        {Name: "net.domain", Tier: types.TierTwo},
	"net.network_cards": {Name: "net.network_cards", Tier: types.TierThree},
	"phy.cpu.cores":     {Name: "phy.cpu.cores", Tier: types.TierOne},
	"phy.ram":           {Name: "phy.ram", Tier: types.TierOne},
}

// Provides registers fact keys produced by this source.
func (w *WindowsSysteminfo) Provides(p types.Provides) {
	for name := range systemInfoFacts {
		p[name] = w
	}
}

// Requires declares program dependency.
func (w *WindowsSysteminfo) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "systeminfo /FO CSV")
}

// Parse executes systeminfo and extracts OS/hardware facts from CSV output.
// Maps CSV header to values and extracts relevant fields.
func (w *WindowsSysteminfo) Parse(f *types.FactDB) {

	data, rc, _ := input.CommandRunner("systeminfo /FO CSV")
	if data == "" || rc != 0 {
		for _, fact := range systemInfoFacts {
			fact.Value = types.ParseFailMsg
			f.AddFact(*fact)
		}
		return
	}

	// Parse CSV: expect header + data row
	r := csv.NewReader(strings.NewReader(data))
	records, err := r.ReadAll()
	if err != nil || len(records) < 2 {
		for _, fact := range systemInfoFacts {
			fact.Value = types.ParseFailMsg
			f.AddFact(*fact)
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
		systemInfoFacts["os.hostname"].Value = v
	} else {
		systemInfoFacts["os.hostname"].Value = types.ParseFailMsg
	}

	if v, ok := fields["OS Name"]; ok && v != "" {
		systemInfoFacts["os.name"].Value = v
	} else {
		systemInfoFacts["os.name"].Value = types.ParseFailMsg
	}

	if v, ok := fields["OS Version"]; ok && v != "" {
		systemInfoFacts["os.version"].Value = v
	} else {
		systemInfoFacts["os.version"].Value = types.ParseFailMsg
	}

	if v, ok := fields["OS Build Type"]; ok && v != "" {
		systemInfoFacts["os.build"].Value = v
	} else {
		systemInfoFacts["os.build"].Value = types.ParseFailMsg
	}

	// Total Physical Memory (in MB, convert to bytes for consistency)
	if v, ok := fields["Total Physical Memory"]; ok && v != "" {
		// Extract numeric part and format as SI string
		cleaned := strings.Fields(v)
		if len(cleaned) > 0 {
			systemInfoFacts["phy.ram"].Value = cleaned[0] + " MB"
		} else {
			systemInfoFacts["phy.ram"].Value = types.ParseFailMsg
		}
	} else {
		systemInfoFacts["phy.ram"].Value = types.ParseFailMsg
	}

	// Processor count (parse "N Processor(s) Installed" format)
	if v, ok := fields["Processor(s)"]; ok && v != "" {
		// Extract first field (count)
		parts := strings.Fields(v)
		if len(parts) > 0 {
			systemInfoFacts["phy.cpu.cores"].Value = parts[0]
		} else {
			systemInfoFacts["phy.cpu.cores"].Value = types.ParseFailMsg
		}
	} else {
		systemInfoFacts["phy.cpu.cores"].Value = types.ParseFailMsg
	}

	// Registered Owner
	if v, ok := fields["Registered Owner"]; ok && v != "" {
		systemInfoFacts["os.owner"].Value = v
	} else {
		systemInfoFacts["os.owner"].Value = types.ParseFailMsg
	}

	// Registered Organization
	if v, ok := fields["Registered Organization"]; ok && v != "" {
		systemInfoFacts["os.organization"].Value = v
	} else {
		systemInfoFacts["os.organization"].Value = types.ParseFailMsg
	}

	// Original Install Date
	if v, ok := fields["Original Install Date"]; ok && v != "" {
		systemInfoFacts["os.install_date"].Value = v
	} else {
		systemInfoFacts["os.install_date"].Value = types.ParseFailMsg
	}

	// System Boot Time
	if v, ok := fields["System Boot Time"]; ok && v != "" {
		systemInfoFacts["run.boot_time"].Value = v
	} else {
		systemInfoFacts["run.boot_time"].Value = types.ParseFailMsg
	}

	// System Manufacturer
	if v, ok := fields["System Manufacturer"]; ok && v != "" {
		systemInfoFacts["phy.manufacturer"].Value = v
	} else {
		systemInfoFacts["phy.manufacturer"].Value = types.ParseFailMsg
	}

	// System Model
	if v, ok := fields["System Model"]; ok && v != "" {
		systemInfoFacts["phy.model"].Value = v
	} else {
		systemInfoFacts["phy.model"].Value = types.ParseFailMsg
	}

	// System Type
	if v, ok := fields["System Type"]; ok && v != "" {
		systemInfoFacts["phy.system_type"].Value = v
	} else {
		systemInfoFacts["phy.system_type"].Value = types.ParseFailMsg
	}

	// BIOS Version (version, date format - extract version only)
	if v, ok := fields["BIOS Version"]; ok && v != "" {
		// Split on comma and take first part (version)
		parts := strings.SplitN(v, ",", 2)
		version := strings.TrimSpace(parts[0])
		if version != "" {
			systemInfoFacts["phy.bios_version"].Value = version
		} else {
			systemInfoFacts["phy.bios_version"].Value = types.ParseFailMsg
		}
	} else {
		systemInfoFacts["phy.bios_version"].Value = types.ParseFailMsg
	}

	// Domain
	if v, ok := fields["Domain"]; ok && v != "" {
		systemInfoFacts["net.domain"].Value = v
	} else {
		systemInfoFacts["net.domain"].Value = types.ParseFailMsg
	}

	// Network Cards
	if v, ok := fields["Network Card(s)"]; ok && v != "" {
		systemInfoFacts["net.network_cards"].Value = v
	} else {
		systemInfoFacts["net.network_cards"].Value = types.ParseFailMsg
	}

	// Add all facts to the FactDB
	for _, fact := range systemInfoFacts {
		f.AddFact(*fact)
	}
}
