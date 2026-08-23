// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package subcmd

// Version subcommand displays version and logo information.

import (
	"fmt"

	"github.com/huntermatthews/clu/internal/global"
	"github.com/huntermatthews/clu/internal/input"
)

// VersionCmd implements the "version" subcommand.
type VersionCmd struct{}

func (f *VersionCmd) Run(stdout input.Stdout, stderr input.Stderr) error {
	fmt.Fprintln(stdout, "  .oooooo.   ooooo        ooooo     ooo")
	fmt.Fprintln(stdout, " d8P'  `Y8b  `888'        `888'     `8'")
	fmt.Fprintln(stdout, "888           888          888       8")
	fmt.Fprintln(stdout, "888           888          888       8")
	fmt.Fprintln(stdout, "888           888          888       8")
	fmt.Fprintln(stdout, "`88b    ooo   888       o  `88.    .8'")
	fmt.Fprintln(stdout, " `Y8bood8P'  o888ooooood8    `YbodP'")
	fmt.Fprintln(stdout, "")

	fmt.Fprintf(stdout, "  version: %s\n", global.GetVersion())

	buildInfo := global.GetBuildInfo()
	if buildInfo.CompilerVersion != "" {
		fmt.Fprintf(stdout, "  compiler version: %s\n", buildInfo.CompilerVersion)
		fmt.Fprintf(stdout, "  main path: %s\n", buildInfo.MainPath)

		fmt.Fprintf(stdout, "  program version: %s\n", buildInfo.MainVersion)
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "  libraries:")
		for _, dep := range buildInfo.Dependencies {
			fmt.Fprintf(stdout, "    %s\n", dep)
		}
		if buildInfo.VCSRevision != "" {
			fmt.Fprintln(stdout, "")
			fmt.Fprintf(stdout, "  vcs revision: %s\n", buildInfo.VCSRevision)
			fmt.Fprintf(stdout, "  vcs time:     %s\n", buildInfo.VCSTime)
			if buildInfo.VCSModified == "true" {
				fmt.Fprintf(stdout, "  vcs modified: yes\n")
			}
		}
	}

	return nil
}
