// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package subcmd

// Go port of src/clu/cmd/report.py (excluding parse_args). Implements the
// core reporting workflow: obtain OpSys, parse early facts, parse requested
// facts, filter by tier and output format (dots, shell, json).

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/huntermatthews/clu/pkg/facts"
	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

// FactsCmd implements the "facts" subcommand.
type FactsCmd struct {
	Tier         types.Tier `name:"tier" short:"t" enum:"1,2,3" default:"1" help:"Tier level (1, 2, or 3)."`
	OutputFormat string     `name:"out" enum:"dots,json,shell" default:"dots" help:"Output format: dots, json, or shell."`
	FactNames    []string   `arg:"" optional:"true" help:"Zero or more fact names to report on."`
}

func (f *FactsCmd) Run(stdout input.Stdout, stderr input.Stderr) error {

	osys := facts.OpSysFactory()
	provides := osys.Provides()
	facts := types.NewFactDB()

	if len(f.FactNames) == 0 {
		for k := range provides {
			f.FactNames = append(f.FactNames, k)
		}
	}

	// Parse early facts first. Early facts are our primitive way of handling inter-source dependancies.
	// They are always parsed, even if not actually needed later.
	parseFactsBySpecs(provides, facts, osys.GetEarlyFacts())

	// Parse requested facts (may include duplicates; source parsing is idempotent).
	parseFactsBySpecs(provides, facts, f.FactNames)

	// Filter down to requested + tier selection.
	outputFacts := filterFacts(facts, f.FactNames, f.Tier)

	doOutput(stdout, stderr, outputFacts, f.OutputFormat)

	return nil
}

// parseFactsBySpecs replicates parse_facts_by_specs: determine sources to run.
func parseFactsBySpecs(provides types.Provides, facts *types.FactDB, specs []string) {
	sourcesToParse := map[types.Sources]struct{}{}
	addSource := func(src interface{}) {
		if s, ok := src.(types.Sources); ok && s != nil {
			sourcesToParse[s] = struct{}{}
		}
	}
	if len(specs) > 0 {
		for _, spec := range specs {
			for key, src := range provides {
				if strings.HasPrefix(key, spec) {
					addSource(src)
				}
			}
		}
	} else {
		for _, src := range provides {
			addSource(src)
		}
	}
	for src := range sourcesToParse {
		src.Parse(facts)
	}
}

// filterFacts mirrors filter_facts: retain facts matching specs and tier.
func filterFacts(parsed *types.FactDB, specs []string, tier types.Tier) map[string]string {
	output := map[string]string{}
	for _, key := range parsed.GetTier(tier) {
		for _, spec := range specs {
			if strings.HasPrefix(key, spec) {
				val, _ := parsed.Get(key)
				output[key] = val
				break
			}
		}
	}
	return output
}

// doOutput dispatches to format-specific output functions.
func doOutput(stdout input.Stdout, stderr input.Stderr, facts map[string]string, format string) {
	switch format {
	case "json":
		outputJSON(stdout, stderr, facts)
	case "shell":
		outputShell(stdout, stderr, facts)
	case "dots":
		fallthrough
	default:
		outputDots(stdout, stderr, facts)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// outputDots prints key: value lines sorted by key.
func outputDots(stdout input.Stdout, _ input.Stderr, facts map[string]string) {
	for _, k := range sortedKeys(facts) {
		fmt.Fprintf(stdout, "%s: %s\n", k, facts[k])
	}
}

// outputShell prints KEY_WITH_UNDERSCORES="value" lines.
func outputShell(stdout input.Stdout, _ input.Stderr, facts map[string]string) {
	for _, k := range sortedKeys(facts) {
		keyVar := strings.ToUpper(strings.ReplaceAll(k, ".", "_"))
		fmt.Fprintf(stdout, "%s=\"%s\"\n", keyVar, facts[k])
	}
}

// outputJSON prints a JSON map of facts.
func outputJSON(stdout input.Stdout, _ input.Stderr, facts map[string]string) {
	data, _ := json.MarshalIndent(facts, "", "  ")
	fmt.Fprintln(stdout, string(data))
}
