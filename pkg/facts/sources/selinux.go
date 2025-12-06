package source

// Go port of src/clu/sources/selinux.py
// Determines SELinux enablement and current mode using `selinuxenabled` and `getenforce`.

import (
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
)

// Selinux source collects os.selinux.enable and os.selinux.mode facts.
type Selinux struct{}

// Provides registers SELinux fact keys.
func (s *Selinux) Provides(p pkg.Provides) {
	p["os.selinux.enable"] = s
	p["os.selinux.mode"] = s
}

// Requires declares external programs used.
func (s *Selinux) Requires(r *pkg.Requires) {
	r.Programs = append(r.Programs, "selinuxenabled", "getenforce")
}

// Parse executes the required commands, mapping exit codes and output.
// selinuxenabled: rc 0 -> True, rc 1 -> False, else ParseFailMsg.
// getenforce: trimmed stdout or ParseFailMsg.
func (s *Selinux) Parse(f *pkg.Facts) {
	_, rc := pkg.CommandRunner("selinuxenabled")
	switch rc {
	case 0:
		f.Set("os.selinux.enable", "True")
	case 1:
		f.Set("os.selinux.enable", "False")
	default:
		f.Set("os.selinux.enable", ParseFailMsg)
	}

	data, _ := pkg.CommandRunner("getenforce")
	if strings.TrimSpace(data) == "" {
		f.Set("os.selinux.mode", ParseFailMsg)
	} else {
		f.Set("os.selinux.mode", strings.TrimSpace(data))
	}
}
