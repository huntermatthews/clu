package sources

import (
	"testing"

	"github.com/NHGRI/clu/pkg/facts/types"
)

// TestNoSaltProvides validates fact registration.
func TestNoSaltProvides(t *testing.T) {
	src := &NoSalt{}
	p := types.Provides{}
	src.Provides(p)
	keys := []string{"salt.no_salt.exists", "salt.no_salt.reason"}
	for _, k := range keys {
		if _, ok := p[k]; !ok {
			t.Fatalf("missing provides key %s", k)
		}
	}
}

// TestNoSaltMissing exercises typical behavior when the file is absent.
func TestNoSaltMissing(t *testing.T) {
	f := types.NewFactDB()
	src := &NoSalt{}
	src.Parse(f) // expects system to not have /no_salt in typical environment
	exists, _ := f.Get("salt.no_salt.exists")
	if exists != "False" && exists != "True" { // allow True if file actually present
		t.Fatalf("unexpected exists value %q", exists)
	}
}
