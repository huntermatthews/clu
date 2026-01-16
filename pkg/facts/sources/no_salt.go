package sources

// Detects presence of a sentinel file /no_salt and records existence and reason.

import (
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

// NoSalt reports whether the /no_salt file exists and its content as a reason.
type NoSalt struct{}

var noSaltFacts = map[string]*types.Fact{
	"salt.no_salt.exists": {Name: "salt.no_salt.exists", Tier: types.TierOne},
	"salt.no_salt.reason": {Name: "salt.no_salt.reason", Tier: types.TierOne},
}

// Provides registers fact keys produced by this source.
func (n *NoSalt) Provides(p types.Provides) {
	for name := range noSaltFacts {
		p[name] = n
	}
}

// Requires declares file dependency.
func (n *NoSalt) Requires(r *types.Requires) {
	r.Files = append(r.Files, "/no_salt")
}

// Parse reads /no_salt optionally. If missing or empty sets exists False and returns.
// Otherwise sets exists True and reason to trimmed file content.
func (n *NoSalt) Parse(f *types.FactDB) {
	// Use DI FileReader (returns content,error). Treat any error or empty
	// trimmed content as missing (optional semantics).
	data, err := input.FileReader("/no_salt")
	if err != nil || strings.TrimSpace(data) == "" {
		noSaltFacts["salt.no_salt.exists"].Value = "False"
		f.AddFact(*noSaltFacts["salt.no_salt.exists"])
		return
	}
	noSaltFacts["salt.no_salt.exists"].Value = "True"
	noSaltFacts["salt.no_salt.reason"].Value = strings.TrimSpace(data)

	// Add all facts to the FactDB
	for _, fact := range noSaltFacts {
		f.AddFact(*fact)
	}
}
