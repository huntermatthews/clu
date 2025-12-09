package sources

import (
	"os"
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
)

func TestSystemVersionPlistProvides(t *testing.T) {
	s := &SystemVersionPlist{}
	p := pkg.NewProvides()
	s.Provides(p)
	for _, k := range []string{"os.name", "os.version", "os.build", "id.build_id"} {
		if _, ok := p[k]; !ok {
			t.Errorf("expected provides key %s", k)
		}
	}
}

func TestSystemVersionPlistParse(t *testing.T) {
	path := "/System/Library/CoreServices/SystemVersion.plist"
	if _, err := os.Stat(path); err != nil {
		t.Skip("SystemVersion.plist not present")
	}
	s := &SystemVersionPlist{}
	facts := pkg.NewFacts()
	s.Parse(facts)
	if !facts.Contains("os.name") {
		t.Errorf("expected os.name fact populated")
	}
}
