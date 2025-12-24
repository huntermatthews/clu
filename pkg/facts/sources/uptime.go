package sources

import (
	"regexp"
	"strings"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/facts/types"
)

// Uptime parses system uptime via `uptime` command.
type Uptime struct{}

func (u *Uptime) Provides(p types.Provides) {
	p["run.uptime"] = u
}

func (u *Uptime) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "uptime")
}

// Match both singular 'user' and plural 'users'.
var uptimeRegex = regexp.MustCompile(`.*up *(.*) \d+ user(?:s)?,? .*`)

func (u *Uptime) Parse(f *types.Facts) {
	f.Add(types.TierOne, "run.uptime", types.ParseFailMsg)

	data, rc := pkg.CommandRunner("uptime")
	if data == "" || rc != 0 {
		return
	}

	m := uptimeRegex.FindStringSubmatch(data)
	uptime := types.ParseFailMsg
	if len(m) > 1 {
		uptime = strings.TrimSuffix(m[1], ",")
		uptime = standardizeSpaces(uptime)
	}
	f.Add(types.TierOne, "run.uptime", uptime)
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
