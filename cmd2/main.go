package main

import (
	"fmt"
	"log/slog"
	"runtime/debug"
)

var version = "unset"

func main() {
	// Create a custom handler with the Debug level
	handler := &CustomHandler{level: slog.LevelDebug}

	// Replace the default logger
	slog.SetDefault(slog.New(handler))

	// Now this will be visible
	slog.Debug("debug logging is enabled")
	slog.Debug("version: " + version)

	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Build info not available")
		return
	}

	fmt.Printf("Build Information:\n")
	fmt.Printf("  GoVersion: %s\n", info.GoVersion)                           // this
	fmt.Printf("  Path: %s\n", info.Path)                                     // not this field
	fmt.Printf("  Main Module: %s @ %s\n", info.Main.Path, info.Main.Version) // this

	fmt.Printf("\n\nDependencies:\n")
	for count, dep := range info.Deps {
		fmt.Printf("  Dependency[%d]: %s@%s\n", count, dep.Path, dep.Version)
	}

	fmt.Printf("\n\nSettings:\n")
	for _, setting := range info.Settings {
		fmt.Printf("  %s: %s\n", setting.Key, setting.Value)
	}
}

// goversion   string
// main module  string
// main version string. **** parse this!!
// deps []string of dep.Path + "@" + dep.Version.   == (devel)
// settings map[string]string of setting.Key -> setting.Value
// CGO_ENABLED:
// GOOS:
// GOARCH:
// vcs: git.   optional
// vcs.revision: 67399a5b661390d23e1c556ae8e8c7e89146a0eb opt
// vcs.time: 2026-01-08T04:11:41Z.  opt
// vcs.modified: true/false.  opt
