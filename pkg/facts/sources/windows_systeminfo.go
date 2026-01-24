package sources

// Parses Windows `systeminfo` output in CSV format to extract OS and hardware facts.
// Expected format: header row (CSV with quoted fields) followed by data row.

import (
	"encoding/csv"
	"regexp"
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// WindowsSysteminfo collects OS and hardware facts from systeminfo output.
type WindowsSysteminfo struct{}

var systemInfoFacts = map[string]*types.Fact{
	"os.kernel.name": {Name: "os.kernel.name", Tier: types.TierThree},

	"os.hostname":       {Name: "os.hostname", Tier: types.TierOne},
	"os.name":           {Name: "os.name", Tier: types.TierOne},
	"os.version.number": {Name: "os.version.number", Tier: types.TierOne},
	"os.version.build":  {Name: "os.version.build", Tier: types.TierTwo},
	// "os.version.service_pack": {Name: "os.version.service_pack", Tier: types.TierOne},
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
	// Initialize all facts to ParseFailMsg
	for _, fact := range systemInfoFacts {
		fact.Value = types.ParseFailMsg
	}

	data, rc, _ := input.CommandRunner("systeminfo /FO CSV")
	if data == "" || rc != 0 {
		for _, fact := range systemInfoFacts {
			f.AddFact(*fact)
		}
		return
	}

	// Parse CSV: expect header + data row
	r := csv.NewReader(strings.NewReader(data))
	records, err := r.ReadAll()
	if err != nil || len(records) < 2 {
		for _, fact := range systemInfoFacts {
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

	// Extract and set facts (only update on success)
	// os.kernel.name MUST follow the pattern of other OS sources,
	// IE, this is what should have been from uname for the kernel name
	// clu expects "Windows" here, and uses that fact for various logic
	systemInfoFacts["os.kernel.name"].Value = "Windows"

	if v, ok := fields["Host Name"]; ok && v != "" {
		systemInfoFacts["os.hostname"].Value = v
	}

	if v, ok := fields["OS Name"]; ok && v != "" {
		systemInfoFacts["os.name"].Value = strings.TrimPrefix(v, "Microsoft ")
	}

	if v, ok := fields["OS Version"]; ok && v != "" {
		// Format: "10.0.19041 N/A Build 19041" or "6.1.7601 Service Pack 1 Build 7601"
		// Pattern: <version> <anything> Build <build>
		re := regexp.MustCompile(`^(\S+)\s+.*Build\s+(\S+)`)
		if matches := re.FindStringSubmatch(v); len(matches) == 3 {
			systemInfoFacts["os.version.number"].Value = matches[1]
			systemInfoFacts["os.version.build"].Value = matches[2]
		}
	}

	if v, ok := fields["OS Build Type"]; ok && v != "" {
		systemInfoFacts["os.build"].Value = v
	}

	if v, ok := fields["Total Physical Memory"]; ok && v != "" {
		if cleaned := strings.Fields(v); len(cleaned) > 0 {
			systemInfoFacts["phy.ram"].Value = cleaned[0] + " MB"
		}
	}

	if v, ok := fields["Processor(s)"]; ok && v != "" {
		if parts := strings.Fields(v); len(parts) > 0 {
			systemInfoFacts["phy.cpu.cores"].Value = parts[0]
		}
	}

	if v, ok := fields["Registered Owner"]; ok && v != "" {
		systemInfoFacts["os.owner"].Value = v
	}

	if v, ok := fields["Registered Organization"]; ok && v != "" {
		systemInfoFacts["os.organization"].Value = v
	}

	if v, ok := fields["Original Install Date"]; ok && v != "" {
		systemInfoFacts["os.install_date"].Value = v
	}

	if v, ok := fields["System Boot Time"]; ok && v != "" {
		systemInfoFacts["run.boot_time"].Value = v
	}

	if v, ok := fields["System Manufacturer"]; ok && v != "" {
		systemInfoFacts["phy.manufacturer"].Value = v
	}

	if v, ok := fields["System Model"]; ok && v != "" {
		systemInfoFacts["phy.model"].Value = v
	}

	if v, ok := fields["System Type"]; ok && v != "" {
		systemInfoFacts["phy.system_type"].Value = v
	}

	if v, ok := fields["BIOS Version"]; ok && v != "" {
		if parts := strings.SplitN(v, ",", 2); len(parts) > 0 {
			if version := strings.TrimSpace(parts[0]); version != "" {
				systemInfoFacts["phy.bios_version"].Value = version
			}
		}
	}

	if v, ok := fields["Domain"]; ok && v != "" {
		systemInfoFacts["net.domain"].Value = v
	}

	if v, ok := fields["Network Card(s)"]; ok && v != "" {
		systemInfoFacts["net.network_cards"].Value = v
	}

	// Add all facts to the FactDB
	for _, fact := range systemInfoFacts {
		f.AddFact(*fact)
	}
}
