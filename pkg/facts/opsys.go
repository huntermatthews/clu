package facts

// Go port of src/clu/opsys/opsys.py providing an operating system abstraction
// that aggregates fact sources and exposes combined Provides and Requires data.
// Concrete OS implementations (e.g. Darwin, Linux) can construct an OpSys with
// an ordered slice of sources plus default and early fact key lists.

import (
	"fmt"

	"github.com/huntermatthews/clu/pkg/facts/sources"
	"github.com/huntermatthews/clu/pkg/facts/types"
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
	uname := &sources.Uname{}
	facts := types.NewFacts()
	uname.Parse(facts)
	if kernel, ok := facts.Get("os.kernel.name"); ok {
		if kernel == "Darwin" {
			return NewDarwin()
		}
		fmt.Printf("Unsupported OS kernel: %s\n", kernel)
	}

	panic("unsupported operating system; only Darwin (macOS) is currently implemented")
}
