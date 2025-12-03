package source

import (
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
)

// SystemVersionPlist parses macOS SystemVersion.plist for OS info.
type SystemVersionPlist struct{}

var svKeys = []string{"os.name", "os.version", "os.build", "id.build_id"}

func (s *SystemVersionPlist) Provides(p pkg.Provides) {
	for _, k := range svKeys {
		p[k] = s
	}
}

func (s *SystemVersionPlist) Requires(r *pkg.Requires) {
	r.Files = append(r.Files, "/System/Library/CoreServices/SystemVersion.plist")
}

func (s *SystemVersionPlist) Parse(f *pkg.Facts) {
	if f.Contains("os.name") {
		return
	}
	path := "/System/Library/CoreServices/SystemVersion.plist"
	data, err := pkg.FileReader(path)
	if err != nil {
		for _, k := range svKeys {
			f.Set(k, ParseFailMsg)
		}
		return
	}
	content := data
	// Simple extraction for expected keys; not a full plist parser.
	name := extractTag(content, "ProductName")
	version := extractTag(content, "ProductVersion")
	build := extractTag(content, "ProductBuildVersion")
	buildID := extractTag(content, "BuildID")
	if name == "" && version == "" && build == "" && buildID == "" {
		for _, k := range svKeys {
			f.Set(k, ParseFailMsg)
		}
		return
	}
	f.Set("os.name", name)
	f.Set("os.version", version)
	f.Set("os.build", build)
	f.Set("id.build_id", buildID)
}

func extractTag(s, key string) string {
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
