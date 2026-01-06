package sources

// Determines SELinux enablement and current mode using `selinuxenabled` and `getenforce`.

import (
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// Selinux source collects os.selinux.enable and os.selinux.mode facts.
type Selinux struct{}

// Provides registers SELinux fact keys.
func (s *Selinux) Provides(p types.Provides) {
	p["os.selinux.enable"] = s
	p["os.selinux.mode"] = s
}

// Requires declares external programs used.
func (s *Selinux) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "selinuxenabled", "getenforce")
}

// Parse executes the required commands, mapping exit codes and output.
// selinuxenabled: rc 0 -> True, rc 1 -> False, else ParseFailMsg.
// getenforce: trimmed stdout or ParseFailMsg.
func (s *Selinux) Parse(f *types.Facts) {
	_, rc, _ := input.CommandRunner("selinuxenabled")
	switch rc {
	case 0:
		f.Set("os.selinux.enable", "True")
	case 1:
		f.Set("os.selinux.enable", "False")
	default:
		f.Set("os.selinux.enable", types.ParseFailMsg)
	}

	data, _, _ := input.CommandRunner("getenforce")
	if strings.TrimSpace(data) == "" {
		f.Set("os.selinux.mode", types.ParseFailMsg)
	} else {
		f.Set("os.selinux.mode", strings.TrimSpace(data))
	}
}
