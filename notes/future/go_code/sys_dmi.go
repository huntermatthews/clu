package sources

// Go port of src/clu/sources/sys_dmi.py
// Reads selected DMI/SMBIOS identification files under /sys/devices/virtual/dmi/id
// when platform is physical, populating system vendor/model/serial/uuid/etc.

import (
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
	types "github.com/huntermatthews/clu/pkg/facts/types"
)

// SysDmi collects hardware identity facts from sysfs.
type SysDmi struct{}

var sysDmiMap = map[string]string{
	"sys.vendor":       "/sys/devices/virtual/dmi/id/sys_vendor",
	"sys.model.family": "/sys/devices/virtual/dmi/id/product_family",
	"sys.model.name":   "/sys/devices/virtual/dmi/id/product_name",
	"sys.serial_no":    "/sys/devices/virtual/dmi/id/product_serial",
	"sys.uuid":         "/sys/devices/virtual/dmi/id/product_uuid",
	"sys.oem":          "/sys/devices/virtual/dmi/id/chassis_vendor",
	"sys.asset_no":     "/sys/devices/virtual/dmi/id/chassis_asset_tag",
}

// Provides registers all DMI fact keys.
func (s *SysDmi) Provides(p types.Provides) {
	for k := range sysDmiMap {
		p[k] = s
	}
}

// Requires declares file dependencies for all DMI paths.
func (s *SysDmi) Requires(r *types.Requires) {
	for _, path := range sysDmiMap {
		r.Files = append(r.Files, path)
	}
}

// Parse populates facts only if phy.platform == "physical". Missing/empty files yield empty string values.
func (s *SysDmi) Parse(f *types.Facts) {
	platform, _ := f.Get("phy.platform")
	if platform != "physical" {
		return
	}
	for key, path := range sysDmiMap {
		data, err := pkg.FileReader(path)
		if err != nil || data == "" {
			f.Set(key, "")
			continue
		}
		f.Set(key, strings.TrimSpace(data))
	}
}
