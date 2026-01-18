package sources

// Determines SELinux enablement and current mode using `selinuxenabled` and `getenforce`.

import (
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// Selinux source collects os.selinux.enable and os.selinux.mode facts.
type Selinux struct{}

var selinuxFacts = map[string]*types.Fact{
	"os.selinux.enable": {Name: "os.selinux.enable", Tier: types.TierOne},
	"os.selinux.mode":   {Name: "os.selinux.mode", Tier: types.TierOne},
}

// Provides registers SELinux fact keys.
func (s *Selinux) Provides(p types.Provides) {
	for name := range selinuxFacts {
		p[name] = s
	}
}

// Requires declares external programs used.
func (s *Selinux) Requires(r *types.Requires) {
	r.Programs = append(r.Programs, "selinuxenabled", "getenforce")
}

// Parse executes the required commands, mapping exit codes and output.
// selinuxenabled: rc 0 -> True, rc 1 -> False, else ParseFailMsg.
// getenforce: trimmed stdout or ParseFailMsg.
func (s *Selinux) Parse(f *types.FactDB) {
	_, rc, _ := input.CommandRunner("selinuxenabled")
	switch rc {
	case 0:
		selinuxFacts["os.selinux.enable"].Value = "True"
	case 1:
		selinuxFacts["os.selinux.enable"].Value = "False"
	default:
		selinuxFacts["os.selinux.enable"].Value = types.ParseFailMsg
	}

	data, _, _ := input.CommandRunner("getenforce")
	if strings.TrimSpace(data) == "" {
		selinuxFacts["os.selinux.mode"].Value = types.ParseFailMsg
	} else {
		selinuxFacts["os.selinux.mode"].Value = strings.TrimSpace(data)
	}

	// Add all facts to the FactDB
	for _, fact := range selinuxFacts {
		f.AddFact(*fact)
	}
}
