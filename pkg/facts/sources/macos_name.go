package source

import (
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
)

// MacOSName derives the macOS marketing code name from the major version.
type MacOSName struct{}

func (m *MacOSName) Provides(p pkg.Provides) { p["os.code_name"] = m }

func (m *MacOSName) Requires(r *pkg.Requires) { r.Facts = append(r.Facts, "os.version") }

func (m *MacOSName) Parse(f *pkg.Facts) {
	ver, ok := f.Get("os.version")
	if !ok || strings.TrimSpace(ver) == "" {
		f.Set("os.code_name", ParseFailMsg)
		return
	}
	major := strings.Split(ver, ".")[0]
	code := codeNameFromMajor(major)
	f.Set("os.code_name", code)
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
		return "Big Sur"
	default:
		return ParseFailMsg
	}
}
