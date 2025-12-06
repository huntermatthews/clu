package source

// Go port of src/clu/sources/__init__.py defining common source interface and constants.

import (
	"github.com/huntermatthews/clu/pkg/types"
)

const (
	ParseFailMsg   = "Unknown/Error"
	NetDisabledMsg = "Unknown - Network Queries Disabled"
)

// Source is the common interface for fact sources.
// Methods mirror Python abstract methods: provides, requires, parse.
type Source interface {
	Provides(p pkg.types.Provides)
	Requires(r *pkg.types.Requires)
	Parse(f *pkg.types.Facts)
}
