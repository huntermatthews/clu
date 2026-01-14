package sources

import (
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/types"
)

// MacOSName derives the macOS marketing code name from the major version.
type MacOSName struct{}

func (m *MacOSName) Provides(p types.Provides) { p["os.code_name"] = m }

func (m *MacOSName) Requires(r *types.Requires) { r.Facts = append(r.Facts, "os.version") }

func (m *MacOSName) Parse(f *types.FactDB) {
	ver, ok := f.Get("os.version")
	if !ok || strings.TrimSpace(ver) == "" {
		f.Add(types.TierOne, "os.code_name", types.ParseFailMsg)
		return
	}

	major := strings.Split(ver, ".")[0]
	code := codeNameFromMajor(major)
	f.Add(types.TierOne, "os.code_name", code)
}

func codeNameFromMajor(major string) string {
	switch major {
	case "26":
		return "Tahoe"
	case "15":
		return "Sequoia"
	case "14":
		return "Sonoma"
	case "13":
		return "Ventura"
	case "12":
		return "Monterey"
	case "11":
		return "Big Sur" // big sur is OLD enough...
	default:
		return types.ParseFailMsg
	}
}
