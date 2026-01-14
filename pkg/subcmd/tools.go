package subcmd

// Check SSA admins workstations for required, recommended and optional tools

import (
	"fmt"

	"github.com/huntermatthews/clu/pkg/input"
)

var requiredTools = []string{
	"uname",
	"sed",
}

var recommendedTools = []string{"tool3", "tool4"}
var optionalTools = []string{"tool5", "tool6"}

// ToolsCmd implements the "tools" subcommand.
type ToolsCmd struct{}

func (f *ToolsCmd) Run(stdout input.Stdout, stderr input.Stderr) error {

	fmt.Println("tool sub-cmd was here.")
	for _, tool := range requiredTools {
		path := input.CheckProgramExists(tool)
		if path == "" {
			fmt.Fprintf(stderr, "required tool %s not found: \n", tool)
		} else {
			fmt.Fprintf(stdout, "required tool %s found at %s\n", tool, path)
		}
	}
	return nil
}
