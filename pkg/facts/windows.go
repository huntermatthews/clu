// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package facts

// Windows operating system source aggregation and early fact lists.

import (
	"github.com/huntermatthews/clu/pkg/facts/sources"
	"github.com/huntermatthews/clu/pkg/facts/types"
)

// NewWindows constructs the Windows OpSys with its ordered sources and fact lists.
func NewWindows() *OpSys {
	return &OpSys{
		Sources: []types.Sources{
			&sources.WindowsSysteminfo{},
			&sources.Clu{},
		},
		EarlyFacts: []string{},
	}
}
