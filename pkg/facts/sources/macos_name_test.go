package sources

import (
	"testing"

	"github.com/huntermatthews/clu/pkg/facts/types"
)

func TestMacOSNameProvides(t *testing.T) {
	m := &MacOSName{}
	p := types.NewProvides()
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
		{"10.15", types.ParseFailMsg}, // unsupported older major
		{"", types.ParseFailMsg},      // empty version
		{"   ", types.ParseFailMsg},   // whitespace version
	}
	for _, c := range cases {
		m := &MacOSName{}
		facts := types.NewFactDB()
		if c.version != "" { // whitespace still set to test trimming
			facts.AddFact(types.Fact{Name: "os.version", Value: c.version})
		}
		m.Parse(facts)
		got, _ := facts.Get("os.code_name")
		if got != c.expected {
			t.Errorf("version %q -> code_name %q; want %q", c.version, got, c.expected)
		}
	}
}
