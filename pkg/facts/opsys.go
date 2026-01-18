package facts

// Go port of src/clu/opsys/opsys.py providing an operating system abstraction
// that aggregates fact sources and exposes combined Provides and Requires data.
// Concrete OS implementations (e.g. Darwin, Linux) can construct an OpSys with
// an ordered slice of sources plus default and early fact key lists.

import (
	"github.com/NHGRI/clu/pkg/facts/sources"
	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// OpSys aggregates a set of fact sources for an operating system.
// DefaultFacts and EarlyFacts mirror the Python methods returning ordered key lists.
type OpSys struct {
	Sources      []types.Sources
	DefaultFacts []string
	EarlyFacts   []string
}

// New constructs a new OpSys instance.
func NewOpSys(s []types.Sources, defaults []string, early []string) *OpSys {
	return &OpSys{Sources: s, DefaultFacts: defaults, EarlyFacts: early}
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

// GetDefaultFacts returns the list of default fact keys (may be empty).
func (o *OpSys) GetDefaultFacts() []string {
	return append([]string{}, o.DefaultFacts...)
}

// GetEarlyFacts returns the list of early fact keys (may be empty).
func (o *OpSys) GetEarlyFacts() []string {
	return append([]string{}, o.EarlyFacts...)
}

// OpSysFactory replicates Python opsys_factory minimal logic using runtime.GOOS.
func OpSysFactory() *OpSys {
	// we look for cmd.exe to determine if we are on Windows
	// ver was my first choice, but its a builtin
	if input.ProgramChecker("cmd.exe") != "" {
		return NewWindows()
	} else {
		// If its not Windows, use Uname to determine OS.
		uname := &sources.Uname{}
		facts := types.NewFacts()
		uname.Parse(facts)
		kernel, ok := facts.Get("os.kernel.name")
		if !ok {
			panic("unable to determine OS kernel name")
		}

	switch kernel {
	case "Darwin":
		return NewDarwin()
	case "Linux":
		return NewLinux()
	default:
		panic("unsupported operating system; got " + kernel)
		switch kernel {
		case "Darwin":
			return NewDarwin()
		case "Linux":
			return NewLinux()
		default:
			panic("unsupported operating system; got " + kernel)

		}
	}
}
