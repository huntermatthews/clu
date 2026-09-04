// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package facts

// Constructs the Linux OpSys with its ordered fact sources and early fact lists.

import (
	"github.com/huntermatthews/clu/internal/facts/sources"
	"github.com/huntermatthews/clu/internal/facts/types"
)

// NewLinux constructs the Linux OpSys. Order mirrors Python minus unported AwsImds source.
func NewLinux() *OpSys {
	return &OpSys{
		Sources: []types.Sources{
			&sources.Clu{},
			&sources.LinuxCheckUpdate{},
			&sources.IpAddr{},
			&sources.Ipmitool{},
			&sources.Lscpu{},
			&sources.Memory{},
			&sources.Salt{},
			&sources.OsRelease{},
			&sources.ProcCpuinfo{},
			&sources.ProcUptime{},
			&sources.Selinux{},
			&sources.SysDmi{},
			&sources.Uname{},
			&sources.Platform{},
		},
		EarlyFacts: []string{
			"phy.arch",
			"phy.platform",
		},
	}
}
