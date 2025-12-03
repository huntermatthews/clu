package source

// Go port of src/clu/sources/os_release.py
// Extracts os.distro.name (ID) and os.distro.version (VERSION_ID) from /etc/os-release.

import (
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
)

// OsRelease parses the /etc/os-release file for distro identification.
type OsRelease struct{}

// Provides registers distro fact keys.
func (o *OsRelease) Provides(p pkg.Provides) {
	p["os.distro.name"] = o
	p["os.distro.version"] = o
}

// Requires declares the file dependency.
func (o *OsRelease) Requires(r *pkg.Requires) { r.Files = append(r.Files, "/etc/os-release") }

// Parse reads /etc/os-release, setting ParseFailMsg on failure. Only ID and VERSION_ID are used.
func (o *OsRelease) Parse(f *pkg.Facts) {
	data, err := pkg.FileReader("/etc/os-release")
	if err != nil || data == "" { // treat empty/error uniformly
		f.Set("os.distro.name", ParseFailMsg)
		f.Set("os.distro.version", ParseFailMsg)
		return
	}
	var id, version string
	for _, line := range strings.Split(data, "\n") {
		if !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.Trim(parts[0], " \"")
		value := strings.Trim(parts[1], " \"")
		switch key {
		case "ID":
			id = value
		case "VERSION_ID":
			version = value
		}
	}
	if id == "" {
		id = ParseFailMsg
	}
	if version == "" {
		version = ParseFailMsg
	}
	f.Set("os.distro.name", id)
	f.Set("os.distro.version", version)
}
