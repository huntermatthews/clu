package facts

// FreeBSD operating system source aggregation and early fact lists.

import (
	"github.com/NHGRI/clu/pkg/facts/sources"
	"github.com/NHGRI/clu/pkg/facts/types"
)

// NewFreeBSD constructs the FreeBSD OpSys with its ordered sources and fact lists.
func NewFreeBSD() *OpSys {
	return &OpSys{
		Sources: []types.Sources{
			&sources.Uname{},
			&sources.FreeBSDSysctl{},
			&sources.Uptime{},
			&sources.Clu{},
		},
		EarlyFacts: []string{},
	}
}
