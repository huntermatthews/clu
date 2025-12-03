package source

import (
	"regexp"
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
)

// Uptime parses system uptime via `uptime` command.
type Uptime struct{}

func (u *Uptime) Provides(p pkg.Provides) { p["run.uptime"] = u }

func (u *Uptime) Requires(r *pkg.Requires) { r.Programs = append(r.Programs, "uptime") }

// Match both singular 'user' and plural 'users'.
var uptimeRegex = regexp.MustCompile(`.*up *(.*) \d+ users? .*`)

func (u *Uptime) Parse(f *pkg.Facts) {
	data, rc := pkg.CommandRunner("uptime")
	if data == "" || rc != 0 {
		f.Set("run.uptime", ParseFailMsg)
		return
	}
	m := uptimeRegex.FindStringSubmatch(data)
	uptime := ParseFailMsg
	if len(m) > 1 {
		uptime = strings.TrimSuffix(m[1], ",")
	}
	f.Set("run.uptime", uptime)
}
