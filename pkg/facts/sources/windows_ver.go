package sources

// Uses 'ver' command to get Windows version information.

import (
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// WindowsVer determines os.version on Windows using the 'ver' command.
type WindowsVer struct{}

// Provides registers os.version key.
func (v *WindowsVer) Provides(p types.Provides) { p["os.version"] = v }

// Requires declares program dependency.
// 'ver' is a shell builtin, so we check for 'cmd' existence.
func (v *WindowsVer) Requires(r *types.Requires) { r.Programs = append(r.Programs, "cmd") }

// Parse sets os.version. Non-zero rc -> types.ParseFailMsg.
// Multiple lines are joined by ", "; empty result -> "physical".
func (v *WindowsVer) Parse(f *types.Facts) {

	// We can assume windows since ver exists at all - even if the parse fails.
	f.Add(types.TierOne, "os.name", "Windows")

	data, rc, _ := input.CommandRunner("ver")
	if rc != 0 {
		f.Add(types.TierOne, "os.version", types.ParseFailMsg)
		return
	}

	// Trim the newline at the end.
	data = strings.TrimSpace(data)

	// Parse "Microsoft Windows [Version X.X.XXXXX.XXXX]"
	// Extract version number from between brackets
	start := strings.Index(data, "[Version ")
	end := strings.Index(data, "]")

	if start != -1 && end != -1 && end > start {
		version := data[start+9 : end] // Skip "[Version "
		f.Add(types.TierOne, "os.version", version)
	} else {
		// Fallback: use full output if format doesn't match
		f.Add(types.TierOne, "os.version", data)
	}
}
