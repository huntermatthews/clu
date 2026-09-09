// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package unix

import (
	"strconv"
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/registry"
	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

// Uname collects simple uname-derived facts by running `uname -snrm`.
func init() {
	registry.Register(func() types.Sources {
		return &Uname{}
	}, "linux", "darwin")
}

type Uname struct{}

// Origin holds the 0-based field position in `uname -snrm` output.
var unameFacts = map[string]*types.Fact{
	"os.kernel.name":    {Name: "os.kernel.name", Tier: types.TierOne, Origin: "0"},
	"os.hostname":       {Name: "os.hostname", Tier: types.TierOne, Origin: "1"},
	"os.kernel.version": {Name: "os.kernel.version", Tier: types.TierTwo, Origin: "2"},
	"phy.arch":          {Name: "phy.arch", Tier: types.TierOne, Origin: "3"},
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

	for _, fact := range unameFacts {
		fact.Value = types.ParseFailMsg
	}
	defer f.AddAll(unameFacts)

	data, rc, _ := input.CommandRunner("uname -snrm")
	if data == "" || rc != 0 {
		return
	}

	fields := strings.Fields(strings.TrimSpace(data))
	ordered := make([]*types.Fact, len(unameFacts))
	for _, fact := range unameFacts {
		i, _ := strconv.Atoi(fact.Origin)
		ordered[i] = fact
	}
	for i, fact := range ordered {
		if i < len(fields) {
			fact.Value = fields[i]
		}
	}
}
