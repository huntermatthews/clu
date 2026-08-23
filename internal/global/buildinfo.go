// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package global

import (
	"fmt"
	"runtime/debug"
)

// BuildInfo contains build-time information gathered from runtime/debug
type BuildInfo struct {
	CompilerVersion string
	MainPath        string
	MainVersion     string
	Dependencies    []string // Each in "path@version" format
	VCSRevision     string
	VCSTime         string
	VCSModified     string
}

// GetBuildInfo gathers all available build information from debug.ReadBuildInfo.
// Returns an empty *BuildInfo if build info is unavailable.
func GetBuildInfo() *BuildInfo {
	bi := &BuildInfo{}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return bi
	}

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

	// Extract VCS settings
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

// GetVersion returns the main module version from build info, or "unset" if
// build info is unavailable (e.g. when running via `go run`).
func GetVersion() string {
	if v := GetBuildInfo().MainVersion; v != "" {
		return v
	}
	return "unset"
}
