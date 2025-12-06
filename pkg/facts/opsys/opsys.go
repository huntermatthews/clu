package opsys

// Go port of src/clu/opsys/opsys.py providing an operating system abstraction
// that aggregates fact sources and exposes combined Provides and Requires data.
// Concrete OS implementations (e.g. Darwin, Linux) can construct an OpSys with
// an ordered slice of sources plus default and early fact key lists.

import (
	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/source"
)

// OpSys aggregates a set of fact sources for an operating system.
// DefaultFacts and EarlyFacts mirror the Python methods returning ordered key lists.
type OpSys struct {
	Sources      []source.Source
	DefaultFacts []string
	EarlyFacts   []string
}

// New constructs a new OpSys instance.
func NewOpSys(s []source.Source, defaults, early []string) *OpSys {
	return &OpSys{Sources: s, DefaultFacts: defaults, EarlyFacts: early}
}

// Provides builds a provider map by invoking Provides on each source in order.
func (o *OpSys) Provides() pkg.types.Provides {
	provs := pkg.types.Provides{}
	for _, src := range o.Sources {
		if src != nil {
			src.Provides(provs)
		}
	}
	return provs
}

// Requires aggregates program/file/API/fact requirements from all sources.
func (o *OpSys) Requires() *pkg.types.Requires {
	reqs := pkg.types.NewRequires()
	for _, src := range o.Sources {
		if src != nil {
			src.Requires(reqs)
		}
	}
	return reqs
}

// GetDefaultFacts returns the list of default fact keys (may be empty).
func (o *OpSys) GetDefaultFacts() []string { return append([]string{}, o.DefaultFacts...) }

// GetEarlyFacts returns the list of early fact keys (may be empty).
func (o *OpSys) GetEarlyFacts() []string { return append([]string{}, o.EarlyFacts...) }
