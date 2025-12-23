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
)

// FactsCmd implements the "facts" subcommand.
type FactsCmd struct {
	Tier         int      `name:"tier" short:"t" enum:"1,2,3" default:"1" help:"Tier level (1, 2, or 3)."`
	OutputFormat string   `name:"out" enum:"dots,json,shell" default:"dots" help:"Output format: dots, json, or shell."`
	FactNames    []string `arg:"" optional:"true" help:"Zero or more fact names to report on."`
}

func (f *FactsCmd) Run() error {
	if len(f.FactNames) > 0 {
		fmt.Printf("TRACE: facts: tier=%d (stub) %s\n", f.Tier, strings.Join(f.FactNames, " "))
	} else {
		fmt.Printf("TRACE: facts: tier=%d (stub)\n", f.Tier)
	}

	osys := facts.OpSysFactory()
	provides := osys.Provides()
	facts := types.NewFacts()

	if len(f.FactNames) == 0 {
		// an empty slice means all facts
		// f.FactNames = make([]string, 0, len(provides))
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
	outputFacts := filterFacts(facts, f.FactNames, int(f.Tier))
	// Output.
	doOutput(outputFacts, f.OutputFormat)

	return nil
}

// parseFactsBySpecs replicates parse_facts_by_specs: determine sources to run.
func parseFactsBySpecs(provides types.Provides, facts *types.Facts, specs []string) {
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
func filterFacts(parsed *types.Facts, specs []string, tier int) *types.Facts {
	output := types.NewFacts()
	tierFacts := parsed.GetTier(types.Tier(tier)) // TierOne=1 mapping preserved
	tierSet := map[string]struct{}{}
	for _, k := range tierFacts {
		tierSet[k] = struct{}{}
	}
	for _, spec := range specs {
		for key, val := range parsed.ToMap() {
			if strings.HasPrefix(key, spec) {
				if _, ok := tierSet[key]; ok {
					output.Set(key, val)
				}
			}
		}
	}
	return output
}

// doOutput dispatches to format-specific output functions.
func doOutput(f *types.Facts, format string) {
	switch format {
	case "json":
		outputJSON(f)
	case "shell":
		outputShell(f)
	case "dots":
		fallthrough
	default:
		outputDots(f)
	}
}

// outputDots prints key: value lines sorted by key.
func outputDots(f *types.Facts) {
	keys := f.Keys()
	sort.Strings(keys)
	for _, k := range keys {
		v, _ := f.Get(k)
		fmt.Printf("%s: %s\n", k, v)
	}
}

// outputShell prints KEY_WITH_UNDERSCORES="value" lines.
func outputShell(f *types.Facts) {
	keys := f.Keys()
	sort.Strings(keys)
	for _, k := range keys {
		v, _ := f.Get(k)
		keyVar := strings.ToUpper(strings.ReplaceAll(k, ".", "_"))
		fmt.Printf("%s=\"%s\"\n", keyVar, v)
	}
}

// outputJSON prints a JSON map of facts.
func outputJSON(f *types.Facts) {
	data, _ := json.MarshalIndent(f.ToMap(), "", "  ")
	fmt.Println(string(data))
}
