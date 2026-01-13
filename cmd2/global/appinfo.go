package global

import (
	"fmt"
	"io"
)

// AppInfo contains static information about the program
type AppInfo struct {
	Name       string
	Version    string
	License    string
	Author     string
	Repository string
}

// GetAppInfo returns the static program information with the provided version
func GetAppInfo(version string) *AppInfo {
	return &AppInfo{
		Name:       "clu",
		Version:    version,
		License:    "Non-Distributable",
		Author:     "Hunter Matthews <hunter@unix.haus>",
		Repository: "https://github.com/huntermatthews/clu",
	}
}

// Print writes program information to the provided writer
func (a *AppInfo) Print(w io.Writer) {
	fmt.Fprintf(w, "Program Information:\n")
	fmt.Fprintf(w, "  Name: %s\n", a.Name)
	fmt.Fprintf(w, "  Version: %s\n", a.Version)
	fmt.Fprintf(w, "  License: %s\n", a.License)
	fmt.Fprintf(w, "  Author: %s\n", a.Author)
	fmt.Fprintf(w, "  Repository: %s\n", a.Repository)
}
