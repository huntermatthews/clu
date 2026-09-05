// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package sources

// Linux platform fact collection. Platform prefers virt-what and retains a
// systemd-detect-virt collector stub for the future fallback implementation.

import (
	"strings"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

const (
	platformFactName         = "phy.platform"
	virtWhatCommand          = "virt-what"
	systemdDetectVirtCommand = "systemd-detect-virt"
)

var (
	platformFact = types.Fact{
		Name: platformFactName,
		Tier: types.TierOne,
	}
	platformSources = []platformCollector{
		&virtWhat{},
		&systemdDetectVirt{},
	}
)

// platformCollector is an internal strategy for deriving phy.platform. Only
// Platform is exposed as a fact source.
type platformCollector interface {
	Parse(*types.FactDB)
}

// Platform collects virtualization-platform facts using the first successful
// Linux collector.
type Platform struct{}

// Provides registers the sole phy.platform provider.
func (p *Platform) Provides(provides types.Provides) { provides[platformFact.Name] = p }

// Requires declares every collector command so the collector subcommand can
// capture the preferred and fallback command outputs.
func (p *Platform) Requires(requires *types.Requires) {
	requires.Programs = append(requires.Programs, virtWhatCommand, systemdDetectVirtCommand)
}

// Parse preserves an existing platform fact; otherwise it tries collectors in
// preference order until one supplies a non-empty, non-error value.
func (p *Platform) Parse(facts *types.FactDB) {
	if facts.Contains(platformFactName) {
		return
	}

	for _, source := range platformSources {
		source.Parse(facts)
		if value, ok := facts.Get(platformFactName); ok && value != "" && value != types.ParseFailMsg {
			return
		}
	}
}

// virtWhat determines the virtualization platform through virt-what.
type virtWhat struct{}

// Parse maps a successful, empty virt-what response to physical hardware and
// joins multiple virtualization identifiers with ", ".
func (v *virtWhat) Parse(facts *types.FactDB) {
	platformFact.Value = types.ParseFailMsg
	defer func() { facts.AddFact(platformFact) }()

	data, rc, _ := input.CommandRunner(virtWhatCommand)
	if rc != 0 {
		return
	}

	data = strings.TrimSpace(data)
	if data == "" {
		platformFact.Value = "physical"
		return
	}
	platformFact.Value = strings.ReplaceAll(data, "\n", ", ")
}

// systemdDetectVirt is a placeholder fallback. The collector subcommand
// already captures its output through Platform.Requires; parsing that output
// will be added when this collector is implemented.
type systemdDetectVirt struct{}

func (s *systemdDetectVirt) Parse(facts *types.FactDB) {
	platformFact.Value = types.ParseFailMsg
	facts.AddFact(platformFact)
}
