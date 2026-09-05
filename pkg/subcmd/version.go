// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package subcmd

// Version subcommand displays version and logo information.

import (
	"fmt"

	"github.com/huntermatthews/clu/pkg/global"
	"github.com/huntermatthews/clu/pkg/input"
)

// VersionCmd implements the "version" subcommand.
type VersionCmd struct{}

func (f *VersionCmd) Run(stdout input.Stdout, stderr input.Stderr) error {
	fmt.Fprintf(stdout, "  .oooooo.   ooooo        ooooo     ooo\n")
	fmt.Fprintf(stdout, " d8P'  `Y8b  `888'        `888'     `8'\n")
	fmt.Fprintf(stdout, "888           888          888       8\n")
	fmt.Fprintf(stdout, "888           888          888       8\n")
	fmt.Fprintf(stdout, "888           888          888       8\n")
	fmt.Fprintf(stdout, "`88b    ooo   888       o  `88.    .8'\n")
	fmt.Fprintf(stdout, " `Y8bood8P'  o888ooooood8    `YbodP'\n")
	fmt.Fprintf(stdout, "\n")

	buildInfo := global.GetBuildInfo()
	if buildInfo.CompilerVersion != "" {
		fmt.Fprintf(stdout, "  version:    %s\n", buildInfo.MainVersion)
		fmt.Fprintf(stdout, "  module:     %s\n", buildInfo.MainPath)
		fmt.Fprintf(stdout, "  compiler:   %s\n", buildInfo.CompilerVersion)

		if buildInfo.VCSRevision != "" {
			fmt.Fprintf(stdout, "  vcs commit: %s\n", buildInfo.VCSRevision)
			fmt.Fprintf(stdout, "  vcs time:   %s\n", buildInfo.VCSTime)
			// fmt.Fprintf(stdout, "  vcs dirty:  %s\n", buildInfo.VCSModified)
		}

		fmt.Fprintf(stdout, "\n  libraries:\n")
		for _, dep := range buildInfo.Dependencies {
			fmt.Fprintf(stdout, "    %s\n", dep)
		}
	}

	return nil
}
