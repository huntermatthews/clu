package sources

import (
	"os/exec"
	"testing"

	"github.com/huntermatthews/clu/pkg/facts/types"
)

func TestUnameProvides(t *testing.T) {
	u := &Uname{}
	p := types.NewProvides()
	u.Provides(p)
	for _, k := range []string{"os.kernel.name", "os.hostname", "os.kernel.version", "phy.arch"} {
		if _, ok := p[k]; !ok {
			t.Errorf("expected provides key %s", k)
		}
	}
}

func TestUnameParse(t *testing.T) {
	if _, err := exec.LookPath("uname"); err != nil {
		t.Skip("uname not available")
	}
	u := &Uname{}
	facts := types.NewFacts()
	u.Parse(facts)
	// After parse we expect at least kernel name populated (unless failure token)
	if !facts.Contains("os.kernel.name") {
		t.Errorf("missing os.kernel.name")
	}
}

func TestUnameParseSkipIfAlreadyPresent(t *testing.T) {
	u := &Uname{}
	facts := types.NewFacts()
	facts.Set("os.kernel.name", "PRESET")
	u.Parse(facts)
	if val, _ := facts.Get("os.kernel.name"); val != "PRESET" {
		t.Errorf("expected preset value not overwritten: got %s", val)
	}
}
