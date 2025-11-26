package sources

import (
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
)

func TestMacOSNameProvides(t *testing.T) {
	m := &MacOSName{}
	p := pkg.NewProvides()
	m.Provides(p)
	if _, ok := p["os.code_name"]; !ok {
		t.Errorf("expected os.code_name to be registered in Provides")
	}
}

func TestMacOSNameParseMappings(t *testing.T) {
	cases := []struct {
		version  string
		expected string
	}{
		{"26.0", "Tahoe"},
		{"15.1.2", "Sequoia"},
		{"14.0", "Sonoma"},
		{"13.3", "Ventura"},
		{"12.6", "Monterey"},
		{"11.7", "Big Sur"},
		{"10.15", ParseFailMsg}, // unsupported older major
		{"", ParseFailMsg},      // empty version
		{"   ", ParseFailMsg},   // whitespace version
	}
	for _, c := range cases {
		m := &MacOSName{}
		facts := pkg.NewFacts()
		if c.version != "" { // Set version if non-empty (whitespace still set to test trimming)
			facts.Set("os.version", c.version)
		} else {
			// Explicitly leave it absent when empty string case
		}
		m.Parse(facts)
		got, _ := facts.Get("os.code_name")
		if got != c.expected {
			t.Errorf("version %q -> code_name %q; want %q", c.version, got, c.expected)
		}
	}
}
