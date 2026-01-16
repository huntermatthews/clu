package sources

import (
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

// SwVers collects macOS version info via `sw_vers` program.
type SwVers struct{}

var swVersFacts = map[string]*types.Fact{
	"os.name":    {Name: "os.name", Tier: types.TierOne},
	"os.version": {Name: "os.version", Tier: types.TierOne},
	"os.build":   {Name: "os.build", Tier: types.TierThree},
}

func (s *SwVers) Provides(p types.Provides) {
	for name := range swVersFacts {
		p[name] = s
	}
}

func (s *SwVers) Requires(r *types.Requires) { r.Programs = append(r.Programs, "sw_vers") }

func (s *SwVers) Parse(f *types.FactDB) {
	if f.Contains("os.name") {
		return
	}

	data, rc, _ := input.CommandRunner("sw_vers")
	if data == "" || rc != 0 {
		for _, fact := range swVersFacts {
			fact.Value = types.ParseFailMsg
			f.AddFact(*fact)
		}
		return
	}

	for _, line := range strings.Split(strings.TrimSpace(data), "\n") {
		if !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "ProductName":
			swVersFacts["os.name"].Value = value
		case "ProductVersion":
			swVersFacts["os.version"].Value = value
		case "BuildVersion":
			swVersFacts["os.build"].Value = value
		}
	}

	// Add all facts to the FactDB
	for _, fact := range swVersFacts {
		f.AddFact(*fact)
	}
}
