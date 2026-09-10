// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package facts

// Provide an operating system abstraction that aggregates fact sources and exposes combined
// Provides and Requires data. Concrete OS implementations (e.g. Darwin, Linux) can construct an
// OpSys with an ordered slice of sources plus early fact key list.

import (
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/registry"
	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"

	_ "github.com/huntermatthews/clu/pkg/facts/darwin"
	_ "github.com/huntermatthews/clu/pkg/facts/linux"
	_ "github.com/huntermatthews/clu/pkg/facts/unix"
	_ "github.com/huntermatthews/clu/pkg/facts/windows"
)

// OpSys aggregates a set of fact sources for an operating system.
// EarlyFacts mirrors the Python method returning ordered key list.
type OpSys struct {
	Sources    []types.Sources
	EarlyFacts []string
}

// Provides builds a provider map by invoking Provides on each source in order.
func (o *OpSys) Provides() types.Provides {
	provs := types.Provides{}
	for _, src := range o.Sources {
		if src != nil {
			src.Provides(provs)
		}
	}
	return provs
}

// Requires aggregates program/file/API/fact requirements from all sources.
func (o *OpSys) Requires() *types.Requires {
	reqs := types.NewRequires()
	for _, src := range o.Sources {
		if src != nil {
			src.Requires(reqs)
		}
	}
	return reqs
}

// GetEarlyFacts returns the list of early fact keys (may be empty).
func (o *OpSys) GetEarlyFacts() []string {
	return append([]string{}, o.EarlyFacts...)
}

// OpSysFactory replicates Python opsys_factory minimal logic using runtime.GOOS.
func newOpSys(system string, earlyFacts []string) *OpSys {
	return &OpSys{
		Sources:    registry.GetSources(system),
		EarlyFacts: earlyFacts,
	}
}

func OpSysFactory() *OpSys {
	// we look for cmd.exe to determine if we are on Windows
	// ver was my first choice, but its a builtin
	if input.ProgramChecker("cmd.exe") != "" {
		return NewWindows()
	}

	kernel, status, err := input.CommandRunner("uname")
	if err != nil || status != 0 {
		panic("unable to determine OS kernel name")
	}

	kernel = strings.TrimSpace(kernel)
	if kernel == "" {
		panic("unable to determine OS kernel name")
	}

	switch kernel {
	case "Darwin":
		return NewDarwin()
	case "Linux":
		return NewLinux()
	default:
		panic("unsupported operating system; got " + kernel)
	}
}

// NewDarwin constructs the Darwin OpSys and its early fact list.
func NewDarwin() *OpSys {
	return newOpSys("darwin", []string{"os.version"})
}

// NewLinux constructs the Linux OpSys and its early fact list.
func NewLinux() *OpSys {
	return newOpSys("linux", []string{
		"phy.arch",
		"phy.platform",
	})
}

// NewWindows constructs the Windows OpSys and its early fact list.
func NewWindows() *OpSys {
	return newOpSys("windows", []string{})
}
