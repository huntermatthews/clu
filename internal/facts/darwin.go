// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package facts

// Go port of src/clu/opsys/darwin.py providing the macOS (Darwin) operating system
// source aggregation and early fact lists.

import (
	"github.com/huntermatthews/clu/internal/facts/sources"
	"github.com/huntermatthews/clu/internal/facts/types"
)

// NewDarwin constructs the Darwin OpSys with its ordered sources and fact lists.
// Order matches the Python implementation for deterministic precedence.
func NewDarwin() *OpSys {
	return &OpSys{
		Sources: []types.Sources{
			&sources.Uname{},
			&sources.SystemVersion{},
			&sources.MacOSName{},
			&sources.Uptime{},
			&sources.Clu{},
		},
		EarlyFacts: []string{"os.version"},
	}
}
