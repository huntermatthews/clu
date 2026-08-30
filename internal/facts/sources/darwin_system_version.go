// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package sources

// Darwin system-version fact collection. SystemVersion prefers the system
// plist and falls back to sw_vers when the plist cannot provide a complete
// OS name, version, and build.

import (
	"strings"

	"github.com/huntermatthews/clu/internal/facts/types"
	"github.com/huntermatthews/clu/internal/input"
)

const (
	systemVersionPlistPath = "/System/Library/CoreServices/SystemVersion.plist"
	swVersCommand          = "sw_vers"
)

var darwinSystemVersionFacts = map[string]*types.Fact{
	"os.name":     {Name: "os.name", Tier: types.TierOne, Origin: "ProductName"},
	"os.version":  {Name: "os.version", Tier: types.TierOne, Origin: "ProductVersion"},
	"os.build":    {Name: "os.build", Tier: types.TierThree, Origin: "ProductBuildVersion"},
	"id.build_id": {Name: "id.build_id", Tier: types.TierThree, Origin: "BuildID"},
}

// darwinSystemVersionCollector is an internal strategy for deriving Darwin
// system-version facts. Only SystemVersion is exposed as a fact source.
type darwinSystemVersionCollector interface {
	Parse(*types.FactDB)
}

var darwinSystemVersionSources = []darwinSystemVersionCollector{
	&systemVersionPlist{},
	&swVers{},
}

// SystemVersion collects Darwin version facts using the first complete source.
type SystemVersion struct{}

// Provides registers the sole provider for Darwin system-version facts.
func (s *SystemVersion) Provides(p types.Provides) {
	for name := range darwinSystemVersionFacts {
		p[name] = s
	}
}

// Requires declares every collector dependency so the collector subcommand can
// capture both the preferred plist and the sw_vers fallback output.
func (s *SystemVersion) Requires(r *types.Requires) {
	r.Files = append(r.Files, systemVersionPlistPath)
	r.Programs = append(r.Programs, swVersCommand)
}

// Parse tries version collectors in preference order. A collector succeeds
// when it supplies all facts that both collectors can provide. id.build_id is
// plist-specific and remains ParseFailMsg when sw_vers is used as the fallback.
func (s *SystemVersion) Parse(f *types.FactDB) {
	if hasCompleteDarwinSystemVersion(f) {
		return
	}

	for _, source := range darwinSystemVersionSources {
		source.Parse(f)
		if hasCompleteDarwinSystemVersion(f) {
			return
		}
	}
}

func hasCompleteDarwinSystemVersion(f *types.FactDB) bool {
	for _, name := range []string{"os.name", "os.version", "os.build"} {
		value, ok := f.Get(name)
		if !ok || value == "" || value == types.ParseFailMsg {
			return false
		}
	}
	return true
}

// systemVersionPlist parses macOS SystemVersion.plist for OS info.
type systemVersionPlist struct{}

func (s *systemVersionPlist) Parse(f *types.FactDB) {
	for _, fact := range darwinSystemVersionFacts {
		fact.Value = types.ParseFailMsg
	}
	defer f.AddAll(darwinSystemVersionFacts)

	data, err := input.FileReader(systemVersionPlistPath)
	if err != nil || data == "" {
		return
	}

	for _, fact := range darwinSystemVersionFacts {
		if value := extractPlistString(data, fact.Origin); value != "" {
			fact.Value = value
		}
	}
}

// swVers collects Darwin version info via the `sw_vers` program.
type swVers struct{}

func (s *swVers) Parse(f *types.FactDB) {
	for _, name := range []string{"os.name", "os.version", "os.build"} {
		darwinSystemVersionFacts[name].Value = types.ParseFailMsg
	}
	defer func() {
		for _, name := range []string{"os.name", "os.version", "os.build"} {
			f.AddFact(*darwinSystemVersionFacts[name])
		}
	}()

	data, rc, _ := input.CommandRunner(swVersCommand)
	if data == "" || rc != 0 {
		return
	}

	factsByOrigin := map[string]*types.Fact{
		"ProductName":    darwinSystemVersionFacts["os.name"],
		"ProductVersion": darwinSystemVersionFacts["os.version"],
		"BuildVersion":   darwinSystemVersionFacts["os.build"],
	}
	for _, line := range strings.Split(strings.TrimSpace(data), "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		if fact, ok := factsByOrigin[strings.TrimSpace(parts[0])]; ok {
			if value := strings.TrimSpace(parts[1]); value != "" {
				fact.Value = value
			}
		}
	}
}

func extractPlistString(data, key string) string {
	keyTag := "<key>" + key + "</key>"
	idx := strings.Index(data, keyTag)
	if idx == -1 {
		return ""
	}
	rest := data[idx+len(keyTag):]
	start := strings.Index(rest, "<string>")
	end := strings.Index(rest, "</string>")
	if start == -1 || end == -1 || end <= start+len("<string>") {
		return ""
	}
	return rest[start+len("<string>") : end]
}
