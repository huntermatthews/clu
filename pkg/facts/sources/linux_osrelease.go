package sources

// Go port of src/clu/sources/os_release.py
// Extracts os.distro.name (ID) and os.distro.version (VERSION_ID) from /etc/os-release.

import (
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// OsRelease parses the /etc/os-release file for distro identification.
type OsRelease struct{}

var osReleaseFacts = map[string]*types.Fact{
	"os.distro.name":    {Name: "os.distro.name", Tier: types.TierOne, Origin: "ID"},
	"os.distro.version": {Name: "os.distro.version", Tier: types.TierOne, Origin: "VERSION_ID"},
	"os.distro.family":  {Name: "os.distro.family", Tier: types.TierTwo, Origin: "ID_LIKE"},
}

// Provides registers distro fact keys.
func (o *OsRelease) Provides(p types.Provides) {
	for name := range osReleaseFacts {
		p[name] = o
	}
}

// Requires declares the file dependency.
func (o *OsRelease) Requires(r *types.Requires) { r.Files = append(r.Files, "/etc/os-release") }

func (o *OsRelease) Parse(f *types.FactDB) {
	// Initialize all facts to ParseFailMsg
	for _, fact := range osReleaseFacts {
		fact.Value = types.ParseFailMsg
	}

	data, err := input.FileReader("/etc/os-release")
	if err != nil || data == "" {
		for _, fact := range osReleaseFacts {
			f.AddFact(*fact)
		}
		return
	}

	// Parse os-release output: "KEY=value" or "KEY=\"value\"" format
	osReleaseMap := make(map[string]string)
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
			osReleaseMap[key] = value
		}
	}

	// Map os-release keys to facts using Origin field
	for _, fact := range osReleaseFacts {
		if v, ok := osReleaseMap[fact.Origin]; ok && v != "" {
			fact.Value = v
		}
	}

	// Override os.distro.family with cleaned/canonical family name
	osReleaseFacts["os.distro.family"].Value = determineDistroFamily(osReleaseMap)

	// Add all facts to the FactDB
	for _, fact := range osReleaseFacts {
		f.AddFact(*fact)
	}
}

// determineDistroFamily determines the canonical distribution family.
// Examines ID_LIKE and falls back to ID if ID_LIKE is empty.
// Returns the most fundamental upstream distribution family.
func determineDistroFamily(osReleaseMap map[string]string) string {
	idLike := strings.ToLower(osReleaseMap["ID_LIKE"])

	var family string
	// Check ID_LIKE for known families (order matters - check most specific first)
	switch {
	case strings.Contains(idLike, "debian"):
		family = "debian"
	case strings.Contains(idLike, "fedora"):
		family = "fedora"
	case strings.Contains(idLike, "arch"):
		family = "arch"
	case idLike == "":
		// If ID_LIKE is empty, use ID as the family (foundational distros)
		family = osReleaseMap["ID_LIKE"]
	default:
		// Return ID_LIKE as-is if no pattern matched
		family = osReleaseMap["ID_LIKE"]
	}

	return family
}
