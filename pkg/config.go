package pkg

import "github.com/huntermatthews/clu/pkg/facts/types"

// config.go provides a simple global configuration container with NO thread safety.
// Fields:
//  - Verbose: enable extra (non-debug) informational output
//  - Debug: enable diagnostic / debug output
//  - Output: selected output format name (e.g., "dots", "json", "shell")
//  - Facts: list of fact name prefixes requested by user / caller

// AppConfig holds runtime configuration flags and selections.
type AppConfig struct {
	Verbose    bool
	Debug      bool
	NetEnabled bool
	Output     string
	Tier       types.Tier
	FactSpecs  []string
}

// Config is the singleton instance backing global configuration access.
var Config = &AppConfig{Tier: types.TierThree}
