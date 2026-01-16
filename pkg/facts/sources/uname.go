package sources

import (
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

// Uname collects simple uname-derived facts by running `uname -snrm`.
type Uname struct{}

var unameFacts = map[string]*types.Fact{
	"os.kernel.name":    {Name: "os.kernel.name", Tier: types.TierOne},
	"os.hostname":       {Name: "os.hostname", Tier: types.TierOne},
	"os.kernel.version": {Name: "os.kernel.version", Tier: types.TierTwo},
	"phy.arch":          {Name: "phy.arch", Tier: types.TierOne},
}

// This ORDER is important and matches the output order of `uname -snrm`.
var unameKeysOrder = []string{
	"os.kernel.name",
	"os.hostname",
	"os.kernel.version",
	"phy.arch",
}

// Provides registers which keys this source provides.
func (u *Uname) Provides(p types.Provides) {
	for name := range unameFacts {
		p[name] = u
	}
}

// Requires declares external programs this source depends on.
func (u *Uname) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "uname -snrm")
}

// Parse populates the provided Facts map with values from `uname -snrm`.
// If the primary key already exists this function is a no-op.
func (u *Uname) Parse(f *types.FactDB) {
	if f.Contains("os.kernel.name") {
		return
	}
	data, rc, _ := input.CommandRunner("uname -snrm")
	if data == "" || rc != 0 {
		for _, fact := range unameFacts {
			fact.Value = types.ParseFailMsg
			f.AddFact(*fact)
		}
		return
	}

	fields := strings.Fields(strings.TrimSpace(data))
	for i, k := range unameKeysOrder {
		if i < len(fields) {
			unameFacts[k].Value = fields[i]
		} else {
			unameFacts[k].Value = types.ParseFailMsg
		}
	}

	// Add all facts to the FactDB
	for _, fact := range unameFacts {
		f.AddFact(*fact)
	}
}
