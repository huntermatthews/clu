// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package linux

// Detects presence of a sentinel file /no_salt and records existence and reason.

import (
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/registry"
	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

// Salt reports whether the /no_salt file exists and its content as a reason.
func init() {
	registry.Register(func() types.Sources {
		return &Salt{}
	}, "linux")
}

type Salt struct{}

var saltFacts = map[string]*types.Fact{
	"salt.no_salt.exists": {Name: "salt.no_salt.exists", Tier: types.TierTwo},
	"salt.no_salt.reason": {Name: "salt.no_salt.reason", Tier: types.TierOne},
	"salt.no_salt.age":    {Name: "salt.no_salt.age", Tier: types.TierThree},
	"salt.highstate.age":  {Name: "salt.highstate.age", Tier: types.TierTwo},
}

// Provides registers fact keys produced by this source.
func (n *Salt) Provides(p types.Provides) {
	for name := range saltFacts {
		p[name] = n
	}
}

// Requires declares file dependency.
func (n *Salt) Requires(r *types.Requires) {
	r.Files = append(r.Files, "/no_salt")
	r.Files = append(r.Files, "/var/run/salt/minion/latest_salt_highstate")
}

// Parse reads /no_salt optionally. Defaults exists to False and reason to n/a;
// overwrites on success.
func (n *Salt) Parse(f *types.FactDB) {
	saltFacts["salt.no_salt.exists"].Value = "False"
	saltFacts["salt.no_salt.reason"].Value = "n/a"
	saltFacts["salt.no_salt.age"].Value = "n/a"
	saltFacts["salt.highstate.age"].Value = "n/a"
	defer f.AddAll(saltFacts)

	data, err := input.FileReader("/no_salt")
	if err != nil || strings.TrimSpace(data) == "" {
		return
	}
	saltFacts["salt.no_salt.exists"].Value = "True"
	saltFacts["salt.no_salt.reason"].Value = strings.TrimSpace(data)

	if mtime, err := input.FileAgeReader("/no_salt"); err == nil {
		saltFacts["salt.no_salt.age"].Value = mtime.Format("2006-01-02 15:04:05Z07:00")
	}

	if mtime, err := input.FileAgeReader("/var/run/salt/minion/latest_salt_highstate"); err == nil {
		saltFacts["salt.highstate.age"].Value = mtime.Format("2006-01-02 15:04:05Z07:00")
	}
}
