package source

import (
	"os/exec"
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
)

func TestSwVersProvides(t *testing.T) {
	s := &SwVers{}
	p := pkg.NewProvides()
	s.Provides(p)
	for _, k := range []string{"os.name", "os.version", "os.build"} {
		if _, ok := p[k]; !ok {
			t.Errorf("expected provides key %s", k)
		}
	}
}

func TestSwVersParse(t *testing.T) {
	if _, err := exec.LookPath("sw_vers"); err != nil {
		t.Skip("sw_vers not available")
	}
	s := &SwVers{}
	facts := pkg.NewFacts()
	s.Parse(facts)
	// On macOS typical keys should be present; tolerate ParseFailMsg on non-mac hosts if command returns empty.
	if !facts.Contains("os.name") {
		t.Errorf("expected os.name fact")
	}
}
