// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package sources

// Go port of src/clu/sources/sys_dmi.py
// Reads selected DMI/SMBIOS identification files under /sys/devices/virtual/dmi/id
// when platform is physical, populating system vendor/model/serial/uuid/etc.

import (
	"strings"

	"github.com/huntermatthews/clu/internal/facts/types"
	"github.com/huntermatthews/clu/internal/input"
)

// SysDmi collects hardware identity facts from sysfs.
type SysDmi struct{}

var sysDmiFacts = map[string]*types.Fact{
	"sys.vendor":       {Name: "sys.vendor", Tier: types.TierOne, Origin: "/sys/devices/virtual/dmi/id/sys_vendor"},
	"sys.model.family": {Name: "sys.model.family", Tier: types.TierOne, Origin: "/sys/devices/virtual/dmi/id/product_family"},
	"sys.model.name":   {Name: "sys.model.name", Tier: types.TierOne, Origin: "/sys/devices/virtual/dmi/id/product_name"},
	"sys.serial_no":    {Name: "sys.serial_no", Tier: types.TierOne, Origin: "/sys/devices/virtual/dmi/id/product_serial"},
	"sys.uuid":         {Name: "sys.uuid", Tier: types.TierOne, Origin: "/sys/devices/virtual/dmi/id/product_uuid"},
	"sys.oem":          {Name: "sys.oem", Tier: types.TierOne, Origin: "/sys/devices/virtual/dmi/id/chassis_vendor"},
	"sys.asset_no":     {Name: "sys.asset_no", Tier: types.TierOne, Origin: "/sys/devices/virtual/dmi/id/chassis_asset_tag"},
}

// Provides registers all DMI fact keys.
func (s *SysDmi) Provides(p types.Provides) {
	for name := range sysDmiFacts {
		p[name] = s
	}
}

// Requires declares file dependencies for all DMI paths.
func (s *SysDmi) Requires(r *types.Requires) {
	for _, fact := range sysDmiFacts {
		r.Files = append(r.Files, fact.Origin)
	}
}

// Parse populates facts only if phy.platform == "physical". Missing/empty files yield empty string values.
// No ParseFailMsg pre-set is needed: the loop below unconditionally assigns every fact before AddAll fires,
// so there is no early-return path that could leave a fact at a stale or zero value.
func (s *SysDmi) Parse(f *types.FactDB) {
	platform, _ := f.Get("phy.platform")
	if platform != "physical" {
		return
	}

	defer f.AddAll(sysDmiFacts)

	for _, fact := range sysDmiFacts {
		data, err := input.FileReader(fact.Origin)
		if err != nil || data == "" {
			fact.Value = ""
		} else {
			fact.Value = strings.TrimSpace(data)
		}
	}
}
