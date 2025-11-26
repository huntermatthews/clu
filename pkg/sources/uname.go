package sources

import (
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
)

// Uname collects simple uname-derived facts by running `uname -snrm`.
type Uname struct{}

var unameKeys = []string{
	"os.kernel.name",
	"os.hostname",
	"os.kernel.version",
	"phy.arch",
}

// Provides registers which keys this source provides.
func (u *Uname) Provides(p pkg.Provides) {
	for _, k := range unameKeys {
		p[k] = u
	}
}

// Requires declares external programs this source depends on.
func (u *Uname) Requires(r *pkg.Requires) {
	r.Programs = append(r.Programs, "uname -snrm")
}

// Parse populates the provided Facts map with values from `uname -snrm`.
// If the primary key already exists this function is a no-op.
func (u *Uname) Parse(f *pkg.Facts) {
	if f.Contains("os.kernel.name") {
		return
	}
	data, rc := pkg.CommandRunner("uname -snrm")
	if data == "" || rc != 0 {
		for _, k := range unameKeys {
			f.Set(k, ParseFailMsg)
		}
		return
	}
	fields := strings.Fields(strings.TrimSpace(data))
	for i, k := range unameKeys {
		if i < len(fields) {
			f.Set(k, fields[i])
		} else {
			f.Set(k, ParseFailMsg)
		}
	}
}
