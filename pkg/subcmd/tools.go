package subcmd

// Check SSA admins workstations for required, recommended and optional tools

import (
	"fmt"
	"text/tabwriter"

	"github.com/NHGRI/clu/pkg/input"
	"github.com/NHGRI/clu/pkg/tools"
)

var specs = []ToolSpec{
	{"python3", Required, tools.VersionSimpleParser},
	{"jq", Required, tools.VersionSimpleParser},
	{"fish", Required, tools.VersionSimpleParser},
	{"packer", Required, tools.VersionSimpleParser},
	{"terraform", Required, tools.VersionSimpleParser},
	{"bash", Required, tools.VersionSimpleParser},
	{"go", Required, tools.VersionSubcmdParser},
	{"sed", Required, tools.VersionSimpleParser}, // special - gnu vs bsd
	{"awk", Required, tools.VersionSimpleParser}, // special - gnu vs bsd
	{"pip3", Required, tools.VersionSimpleParser},
	{"curl", Required, tools.VersionSimpleParser},
	{"pass", Required, tools.VersionPassParser},
	{"gpg", Required, tools.VersionSimpleParser},
	{"ssh", Required, tools.VersionOpensshParser},
	{"git", Required, tools.VersionSimpleParser},
	{"make", Required, tools.VersionSimpleParser},

	{"yq", Recommended, tools.VersionSimpleParser},
	{"uv", Recommended, tools.VersionSimpleParser},
	{"salt", Recommended, tools.VersionSimpleParser},
	{"zsh", Recommended, tools.VersionSimpleParser},
	{"pkcs15-tool", Recommended, tools.VersionSimpleParser},
	{"pkcs11-tool", Recommended, tools.VersionSimpleParser},
	{"just", Recommended, tools.VersionSimpleParser},

	{"lua", Optional, tools.VersionSimpleParser},
	{"stow", Optional, tools.VersionSimpleParser},
	{"lmod", Optional, tools.VersionSimpleParser},
	{"pipx", Optional, tools.VersionSimpleParser},
	{"keyring", Optional, tools.VersionSimpleParser},
	{"age", Optional, tools.VersionSimpleParser}, // maybe one day
	{"vpn", Optional, tools.VersionSimpleParser},
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
	fmt.Fprintln(w, "TOOL\tVERSION\tPATH")

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
