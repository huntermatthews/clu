package sources

import (
	"strings"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/facts/types"
)

// SystemVersionPlist parses macOS SystemVersion.plist for OS info.
type SystemVersionPlist struct{}

var svKeys = []string{"os.name", "os.version", "os.build", "id.build_id"}

func (s *SystemVersionPlist) Provides(p types.Provides) {
	for _, k := range svKeys {
		p[k] = s
	}
}

func (s *SystemVersionPlist) Requires(r *types.Requires) {
	r.Files = append(r.Files, "/System/Library/CoreServices/SystemVersion.plist")
}

func (s *SystemVersionPlist) Parse(f *types.Facts) {
	if f.Contains("os.name") {
		return
	}

	for _, k := range svKeys {
		// failures are tier one to make sure they are visible
		f.Add(types.TierOne, k, types.ParseFailMsg)
	}

	path := "/System/Library/CoreServices/SystemVersion.plist"
	data, err := pkg.FileReader(path)
	if err != nil {
		return
	}

	// Simple extraction for expected keys; not a full plist parser.
	name := extractTag(data, "ProductName")
	version := extractTag(data, "ProductVersion")
	build := extractTag(data, "ProductBuildVersion")
	buildID := extractTag(data, "BuildID")

	f.Add(types.TierOne, "os.name", name)
	f.Add(types.TierOne, "os.version", version)
	f.Add(types.TierThree, "os.build", build)
	f.Add(types.TierThree, "id.build_id", buildID)
}

func extractTag(s, key string) string {
	// FIXME: This is a very naive xml plist parser; consider using a proper xml  library. (AI did this)
	// Look for <key>key</key><string>value</string>
	keyTag := "<key>" + key + "</key>"
	idx := strings.Index(s, keyTag)
	if idx == -1 {
		return ""
	}
	rest := s[idx+len(keyTag):]
	start := strings.Index(rest, "<string>")
	end := strings.Index(rest, "</string>")
	if start == -1 || end == -1 || end <= start+len("<string>") {
		return ""
	}
	return rest[start+len("<string>") : end]
}
