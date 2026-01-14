package sources

import (
	"os/exec"
	"testing"

	"github.com/huntermatthews/clu/pkg/facts/types"
)

func TestSwVersProvides(t *testing.T) {
	s := &SwVers{}
	p := types.NewProvides()
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
	facts := types.NewFactDB()
	s.Parse(facts)
	if !facts.Contains("os.name") {
		t.Errorf("expected os.name fact")
	}
}
