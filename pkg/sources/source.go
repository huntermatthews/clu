package sources

// Go port of src/clu/sources/__init__.py defining common source interface and constants.

import (
	pkg "github.com/huntermatthews/clu/pkg"
)

const (
	ParseFailMsg   = "Unknown/Error"
	NetDisabledMsg = "Unknown - Network Queries Disabled"
)

// Source is the common interface for fact sources.
// Methods mirror Python abstract methods: provides, requires, parse.
type Source interface {
	Provides(p pkg.Provides)
	Requires(r *pkg.Requires)
	Parse(f *pkg.Facts)
}
