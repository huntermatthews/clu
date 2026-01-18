package sources

import (
	"regexp"
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// Uptime parses system uptime via `uptime` command.
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
	data, rc, _ := input.CommandRunner("uptime")
	if data == "" || rc != 0 {
		uptimeFactCmd.Value = types.ParseFailMsg
		f.AddFact(uptimeFactCmd)
		return
	}

	m := uptimeRegex.FindStringSubmatch(data)
	if len(m) > 1 {
		uptime := strings.TrimSuffix(m[1], ",")
		uptimeFactCmd.Value = strings.Join(strings.Fields(uptime), " ")
	} else {
		uptimeFactCmd.Value = types.ParseFailMsg
	}
	f.AddFact(uptimeFactCmd)
}
