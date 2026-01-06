package sources

import (
	"strings"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// SwVers collects macOS version info via `sw_vers` program.
type SwVers struct{}

var swKeys = []string{
	"os.name",
	"os.version",
	"os.build",
}

func (s *SwVers) Provides(p types.Provides) {
	for _, k := range swKeys {
		p[k] = s
	}
}

func (s *SwVers) Requires(r *types.Requires) { r.Programs = append(r.Programs, "sw_vers") }

func (s *SwVers) Parse(f *types.Facts) {
	if f.Contains("os.name") {
		return
	}
	data, rc, _ := input.CommandRunner("sw_vers")
	if data == "" || rc != 0 {
		for _, k := range swKeys {
			f.Set(k, types.ParseFailMsg)
		}
		return
	}
	for _, line := range strings.Split(strings.TrimSpace(data), "\n") {
		if !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "ProductName":
			f.Set("os.name", value)
		case "ProductVersion":
			f.Set("os.version", value)
		case "BuildVersion":
			f.Set("os.build", value)
		}
	}
}
