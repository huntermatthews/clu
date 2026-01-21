package subcmd

// Check SSA admins workstations for required, recommended and optional tools

import (
	"fmt"
	"text/tabwriter"

	"github.com/NHGRI/clu/pkg/input"
	"github.com/NHGRI/clu/pkg/tools"
)

var specs = []ToolSpec{
	{"python3", Required, tools.DefaultVersionParser},
	{"jq", Required, tools.DefaultVersionParser},
	{"fish", Required, tools.DefaultVersionParser},
	{"packer", Required, tools.DefaultVersionParser},
	{"terraform", Required, tools.DefaultVersionParser},
	{"bash", Required, tools.DefaultVersionParser},
	{"go", Required, tools.DefaultVersionParser},
	{"sed", Required, tools.DefaultVersionParser}, // special - gnu vs bsd
	{"awk", Required, tools.DefaultVersionParser}, // special - gnu vs bsd
	{"pip3", Required, tools.DefaultVersionParser},
	{"curl", Required, tools.DefaultVersionParser},
	{"pass", Required, tools.VersionPassParser},
	{"gpg", Required, tools.DefaultVersionParser},
	{"ssh", Required, tools.VersionOpensshParser},
	{"git", Required, tools.DefaultVersionParser},
	{"make", Required, tools.DefaultVersionParser},

	{"yq", Recommended, tools.DefaultVersionParser},
	{"uv", Recommended, tools.DefaultVersionParser},
	{"salt", Recommended, tools.DefaultVersionParser},
	{"zsh", Recommended, tools.DefaultVersionParser},
	{"pkcs15-tool", Recommended, tools.DefaultVersionParser},
	{"pkcs11-tool", Recommended, tools.DefaultVersionParser},
	{"just", Recommended, tools.DefaultVersionParser},

	{"lua", Optional, tools.DefaultVersionParser},
	{"stow", Optional, tools.DefaultVersionParser},
	{"lmod", Optional, tools.DefaultVersionParser},
	{"pipx", Optional, tools.DefaultVersionParser},
	{"keyring", Optional, tools.DefaultVersionParser},
	{"age", Optional, tools.DefaultVersionParser}, // maybe one day
	{"vpn", Optional, tools.DefaultVersionParser},
}

type ToolCategory string

// Declare the possible constant values for ToolCategory
const (
	All         ToolCategory = "all"
	Required    ToolCategory = "req"
	Recommended ToolCategory = "rec"
	Optional    ToolCategory = "opt"
)

type VersionParser func(path string) string

type ToolSpec struct {
	Name         string
	Category     ToolCategory
	ParseVersion VersionParser
}

// ToolsCmd implements the "tools" subcommand.
type ToolsCmd struct {
	Category ToolCategory `name:"cat" enum:"all,req,rec,opt" default:"req" help:"Tool category to check: all, required (req), recommended (rec), or optional (opt)."`
}

func (f *ToolsCmd) Run(stdout input.Stdout, stderr input.Stderr) error {
	w := tabwriter.NewWriter(stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Tool\tVersion\tPath")
	fmt.Fprintln(w, "====\t=======\t====")

	for _, spec := range specs {
		if f.Category == All || spec.Category == f.Category {
			processTool(spec, w)
		}
	}
	w.Flush()
	return nil
}

// processTool checks if a tool exists and reports its status
func processTool(spec ToolSpec, w *tabwriter.Writer) {
	path := input.CheckProgramExists(spec.Name)
	if path == "" {
		fmt.Fprintf(w, "%s\t%s\t%s\n", spec.Name, "-", "NOT FOUND")
	} else {
		version := spec.ParseVersion(path)
		fmt.Fprintf(w, "%s\t%s\t%s\n", spec.Name, version, path)
	}
}
