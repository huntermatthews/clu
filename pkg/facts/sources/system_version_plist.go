package sources

import (
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// SystemVersionPlist parses macOS SystemVersion.plist for OS info.
type SystemVersionPlist struct{}

var systemVersionFacts = map[string]*types.Fact{
	"os.name":     {Name: "os.name", Tier: types.TierOne, Origin: "ProductName"},
	"os.version":  {Name: "os.version", Tier: types.TierOne, Origin: "ProductVersion"},
	"os.build":    {Name: "os.build", Tier: types.TierThree, Origin: "ProductBuildVersion"},
	"id.build_id": {Name: "id.build_id", Tier: types.TierThree, Origin: "BuildID"},
}

func (s *SystemVersionPlist) Provides(p types.Provides) {
	for name := range systemVersionFacts {
		p[name] = s
	}
}

func (s *SystemVersionPlist) Requires(r *types.Requires) {
	r.Files = append(r.Files, "/System/Library/CoreServices/SystemVersion.plist")
}

func (s *SystemVersionPlist) Parse(f *types.FactDB) {
	if f.Contains("os.name") {
		return
	}

	path := "/System/Library/CoreServices/SystemVersion.plist"
	data, err := input.FileReader(path)
	if err != nil {
		// failures are tier one to make sure they are visible
		for _, fact := range systemVersionFacts {
			fact.Value = types.ParseFailMsg
			f.AddFact(*fact)
		}
		return
	}

	// Simple extraction for expected keys; not a full plist parser.
	for _, fact := range systemVersionFacts {
		fact.Value = extractTag(data, fact.Origin)
	}

	// Add all facts to the FactDB
	for _, fact := range systemVersionFacts {
		f.AddFact(*fact)
	}
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
