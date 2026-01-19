package global

import (
	"fmt"
	"io"
	"runtime/debug"
)

// BuildInfo contains build-time information gathered from runtime/debug
type BuildInfo struct {
	CompilerVersion string
	MainPath        string
	MainVersion     string
	Dependencies    []string // Each in "path@version" format
	CGOEnabled      string
	GOOS            string
	GOARCH          string
	VCSRevision     string
	VCSTime         string
	VCSModified     string
}

// GetBuildInfo gathers all available build information from debug.ReadBuildInfo.
// Returns nil if build info is not available.
func GetBuildInfo() *BuildInfo {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}

	bi := &BuildInfo{}

	// Set build information from debug.ReadBuildInfo
	bi.CompilerVersion = info.GoVersion
	bi.MainPath = info.Main.Path
	bi.MainVersion = info.Main.Version

	// Collect dependencies
	if len(info.Deps) > 0 {
		bi.Dependencies = make([]string, len(info.Deps))
		for i, dep := range info.Deps {
			bi.Dependencies[i] = fmt.Sprintf("%s@%s", dep.Path, dep.Version)
		}
	}

	// Extract relevant settings
	bi.CGOEnabled = getSetting(info.Settings, "CGO_ENABLED")
	bi.GOOS = getSetting(info.Settings, "GOOS")
	bi.GOARCH = getSetting(info.Settings, "GOARCH")
	bi.VCSRevision = getSetting(info.Settings, "vcs.revision")
	bi.VCSTime = getSetting(info.Settings, "vcs.time")
	bi.VCSModified = getSetting(info.Settings, "vcs.modified")

	return bi
}

// getSetting is a helper function to extract a setting value by key
func getSetting(settings []debug.BuildSetting, key string) string {
	for _, setting := range settings {
		if setting.Key == key {
			return setting.Value
		}
	}
	return ""
}

// Print writes build information to the provided writer
func (b *BuildInfo) Print(w io.Writer) {
	fmt.Fprintf(w, "\nBuild Information:\n")
	fmt.Fprintf(w, "  Compiler Version: %s\n", b.CompilerVersion)
	fmt.Fprintf(w, "  Main Path: %s\n", b.MainPath)
	fmt.Fprintf(w, "  Main Version: %s\n", b.MainVersion)

	fmt.Fprintf(w, "\nDependencies:\n")
	if len(b.Dependencies) > 0 {
		for count, dep := range b.Dependencies {
			fmt.Fprintf(w, "  Dependency[%d]: %s\n", count, dep)
		}
	} else {
		fmt.Fprintf(w, "  (none)\n")
	}

	fmt.Fprintf(w, "\nBuild Settings:\n")
	fmt.Fprintf(w, "  CGO_ENABLED: %s\n", b.CGOEnabled)
	fmt.Fprintf(w, "  GOOS: %s\n", b.GOOS)
	fmt.Fprintf(w, "  GOARCH: %s\n", b.GOARCH)

	fmt.Fprintf(w, "\nVersion Control:\n")
	if b.VCSRevision != "" {
		fmt.Fprintf(w, "  vcs.revision: %s\n", b.VCSRevision)
		fmt.Fprintf(w, "  vcs.time: %s\n", b.VCSTime)
		fmt.Fprintf(w, "  vcs.modified: %s\n", b.VCSModified)
	} else {
		fmt.Fprintf(w, "  (no VCS information available)\n")
	}
}
