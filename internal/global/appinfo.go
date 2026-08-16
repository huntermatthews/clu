// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package global

import (
	"fmt"
	"io"
)

var Version = "unset"

// AppInfo contains static information about the program
type AppInfo struct {
	Name       string
	Version    string
	License    string
	Author     string
	Repository string
}

// GetAppInfo returns the static program information
func GetAppInfo() *AppInfo {
	return &AppInfo{
		Name:       "clu",
		Version:    Version,
		License:    "LGPL-2.1-only",
		Author:     "Hunter Matthews <hunter@huntermatthews.net>",
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
