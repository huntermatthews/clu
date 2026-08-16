// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package sources

import (
	"strings"

	"github.com/huntermatthews/clu/internal/facts/types"
	"github.com/huntermatthews/clu/internal/input"
)

// SwVers collects macOS version info via `sw_vers` program.
type SwVers struct{}

var swVersFacts = map[string]*types.Fact{
	"os.name":    {Name: "os.name", Tier: types.TierOne, Origin: "ProductName"},
	"os.version": {Name: "os.version", Tier: types.TierOne, Origin: "ProductVersion"},
	"os.build":   {Name: "os.build", Tier: types.TierThree, Origin: "BuildVersion"},
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

	for _, fact := range swVersFacts {
		fact.Value = types.ParseFailMsg
	}
	defer f.AddAll(swVersFacts)

	data, rc, _ := input.CommandRunner("sw_vers")
	if data == "" || rc != 0 {
		return
	}

	// Build a lookup from sw_vers key name -> fact pointer.
	byOrigin := make(map[string]*types.Fact, len(swVersFacts))
	for _, fact := range swVersFacts {
		byOrigin[fact.Origin] = fact
	}

	for _, line := range strings.Split(strings.TrimSpace(data), "\n") {
		if !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])
		if fact, ok := byOrigin[key]; ok {
			fact.Value = strings.TrimSpace(parts[1])
		}
	}
}
