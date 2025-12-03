package source

// Go port of src/clu/sources/virt_what.py
// Uses `virt-what` to determine virtualization platform(s); empty output implies physical.

import (
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
)

// VirtWhat determines phy.platform if not already set.
type VirtWhat struct{}

// Provides registers phy.platform key.
func (v *VirtWhat) Provides(p pkg.Provides) { p["phy.platform"] = v }

// Requires declares program dependency.
func (v *VirtWhat) Requires(r *pkg.Requires) { r.Programs = append(r.Programs, "virt-what") }

// Parse sets phy.platform unless already present. Non-zero rc -> ParseFailMsg.
// Multiple lines are joined by ", "; empty result -> "physical".
func (v *VirtWhat) Parse(f *pkg.Facts) {
	if f.Contains("phy.platform") {
		return
	}
	data, rc := pkg.CommandRunner("virt-what")
	if rc != 0 {
		f.Set("phy.platform", ParseFailMsg)
		return
	}
	data = strings.TrimSpace(data)
	if data == "" {
		f.Set("phy.platform", "physical")
		return
	}
	// Replace newlines with comma-space
	data = strings.ReplaceAll(data, "\n", ", ")
	f.Set("phy.platform", data)
}
