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

	fmt.Fprintf(stdout, "  version: %s\n", global.Version)

	buildInfo := global.GetBuildInfo()
	if buildInfo != nil {
		fmt.Fprintf(stdout, "  compiler version: %s\n", buildInfo.CompilerVersion)
		fmt.Fprintf(stdout, "  main path: %s\n", buildInfo.MainPath)
		fmt.Fprintf(stdout, "  build version: %s\n", buildInfo.MainVersion)
		fmt.Fprintf(stdout, "  app version: %s\n", global.Version)
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "  libraries:")
		for _, dep := range buildInfo.Dependencies {
			fmt.Fprintf(stdout, "    %s\n", dep)
		}
	}

	return nil
}
