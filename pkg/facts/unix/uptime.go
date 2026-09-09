// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package unix

import (
	"regexp"
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/registry"
	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

// Uptime parses system uptime via `uptime` command.
func init() {
	registry.Register(func() types.Sources {
		return &Uptime{}
	}, "darwin")
}

type Uptime struct{}

var uptimeFactCmd = types.Fact{
	Name: "run.uptime",
	Tier: types.TierOne,
}

func (u *Uptime) Provides(p types.Provides) {
	p[uptimeFactCmd.Name] = u
}

func (u *Uptime) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "uptime")
}

// Match both singular 'user' and plural 'users'.
var uptimeRegex = regexp.MustCompile(`.*up *(.*) \d+ user(?:s)?,? .*`)

func (u *Uptime) Parse(f *types.FactDB) {
	uptimeFactCmd.Value = types.ParseFailMsg
	defer func() { f.AddFact(uptimeFactCmd) }()

	data, rc, _ := input.CommandRunner("uptime")
	if data == "" || rc != 0 {
		return
	}

	m := uptimeRegex.FindStringSubmatch(data)
	if len(m) > 1 {
		uptime := strings.TrimSuffix(m[1], ",")
		uptimeFactCmd.Value = strings.Join(strings.Fields(uptime), " ")
	}
}
