// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package subcmd

// Output an embedded copy of the man page.
// macos/linux: | man -l -

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/huntermatthews/clu/pkg/input"
)

//go:embed clu.1
var manPage string

// Help implements the "tools" subcommand.
type HelpCmd struct{}

func (f *HelpCmd) Run(stdout input.Stdout, stderr input.Stderr) error {

	// Format man page - try mandoc (macOS/BSD) first, then groff (Linux)
	formatter := "mandoc"
	args := []string{"-Tutf8"}
	if _, err := exec.LookPath("mandoc"); err != nil {
		formatter = "groff"
		args = []string{"-man", "-Tutf8"}
	}

	formatCmd := exec.Command(formatter, args...)
	formatCmd.Stdin = bytes.NewReader([]byte(manPage))
	var formatStderr bytes.Buffer
	formatCmd.Stderr = &formatStderr

	// Get the formatted output
	formattedOutput, err := formatCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to format man page: %w (stderr: %s)", err, formatStderr.String())
	}

	// Determine which pager to use
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "less"
	}

	// Pipe through the pager; pass -R to less so ANSI escape sequences
	// (bold, underline, etc.) are rendered rather than printed literally.
	var pagerArgs []string
	if pager == "less" || strings.HasSuffix(pager, "/less") {
		pagerArgs = []string{"-R"}
	}
	pagerCmd := exec.Command(pager, pagerArgs...)
	pagerCmd.Stdin = bytes.NewReader(formattedOutput)
	pagerCmd.Stdout = stdout
	pagerCmd.Stderr = stderr

	if err := pagerCmd.Run(); err != nil {
		return fmt.Errorf("failed to run pager: %w", err)
	}

	return nil
}

func init() {
	Register(Registration{
		Name:    "help",
		Help:    "Show embedded man page.",
		Command: &HelpCmd{},
	})
}
