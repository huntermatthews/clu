package cmd

// Go port of src/clu/cmd/report.py (excluding parse_args). Implements the
// core reporting workflow: obtain OpSys, parse early facts, parse requested
// facts, filter by tier and output format (dots, shell, json).

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	pkg "github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/opsys"
	"github.com/huntermatthews/clu/pkg/sources"
)

// ReportConfig holds runtime configuration for the report command. In Python
// this derived from global cfg (set_config + argparse). Here we expose a struct
// expecting the caller (CLI layer) to populate it; defaults applied via
// SetReportDefaults.
type ReportConfig struct {
	Output string   // dots | shell | json
	Net    bool     // enable network access (placeholder, not yet used)
	Tier   int      // 1,2,3 controlling fact visibility
	Facts  []string // requested fact spec prefixes; empty means all defaults
}

// SetReportDefaults sets missing values analogous to set_report_defaults().
func SetReportDefaults(cfg *ReportConfig, allFactKeys []string) {
	if cfg.Output == "" {
		cfg.Output = "dots"
	}
	if cfg.Tier == 0 { // treat 0 as unset
		cfg.Tier = 1
	}
	if len(cfg.Facts) == 0 {
		cfg.Facts = append([]string{}, allFactKeys...)
	}
}

// parseFactsBySpecs replicates parse_facts_by_specs: determine sources to run.
func parseFactsBySpecs(provides pkg.Provides, facts *pkg.Facts, specs []string) {
	sourcesToParse := map[sources.Source]struct{}{}
	addSource := func(src interface{}) {
		if s, ok := src.(sources.Source); ok && s != nil {
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
func filterFacts(parsed *pkg.Facts, specs []string, tier int) *pkg.Facts {
	output := pkg.NewFacts()
	tierFacts := parsed.GetTier(pkg.Tier(tier)) // TierOne=1 mapping preserved
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
func doOutput(f *pkg.Facts, format string) {
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

// ReportFacts executes the report workflow, returning an exit code (0 success).
// Caller is responsible for printing errors (none currently) or using result.
func ReportFacts(cfg *ReportConfig) int {
	// Build OS abstraction (simple factory: detect GOOS). Simplified vs Python opsys_factory.
	var osys *opsys.OpSys
	// Basic platform detection using runtime.GOOS; only darwin/linux expected now.
	if osVal := runtimeGOOS(); osVal == "darwin" {
		osys = opsys.NewDarwin()
	} else {
		osys = opsys.NewLinux()
	}
	providesMap := osys.Provides()
	parsed := pkg.NewFacts()

	// Apply defaults if caller did not set certain fields.
	allKeys := make([]string, 0, len(providesMap))
	for k := range providesMap {
		allKeys = append(allKeys, k)
	}
	SetReportDefaults(cfg, allKeys)

	// Parse early facts first.
	parseFactsBySpecs(providesMap, parsed, osys.GetEarlyFacts())
	// Parse requested facts (may include duplicates; source parsing idempotent for now).
	parseFactsBySpecs(providesMap, parsed, cfg.Facts)
	// Filter down to requested + tier selection.
	outputFacts := filterFacts(parsed, cfg.Facts, cfg.Tier)
	// Output.
	doOutput(outputFacts, cfg.Output)
	return 0
}

// runtimeGOOS isolated for easier testing (DI if needed later).
var runtimeGOOS = func() string {
	if v := os.Getenv("GOOS_OVERRIDE"); v != "" {
		return v
	}
	return runtime.GOOS
}

// outputDots prints key: value lines sorted by key.
func outputDots(f *pkg.Facts) {
	keys := f.Keys()
	sort.Strings(keys)
	for _, k := range keys {
		v, _ := f.Get(k)
		fmt.Printf("%s: %s\n", k, v)
	}
}

// outputShell prints KEY_WITH_UNDERSCORES="value" lines.
func outputShell(f *pkg.Facts) {
	keys := f.Keys()
	sort.Strings(keys)
	for _, k := range keys {
		v, _ := f.Get(k)
		keyVar := strings.ToUpper(strings.ReplaceAll(k, ".", "_"))
		fmt.Printf("%s=\"%s\"\n", keyVar, v)
	}
}

// outputJSON prints a JSON map of facts.
func outputJSON(f *pkg.Facts) {
	data, _ := json.MarshalIndent(f.ToMap(), "", "  ")
	fmt.Println(string(data))
}
