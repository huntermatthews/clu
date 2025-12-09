package sources

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/huntermatthews/clu/pkg/facts"
)

// TestCluProvides ensures all expected keys are registered.
func TestCluProvides(t *testing.T) {
	c := &Clu{}
	p := facts.NewProvides()
	c.Provides(p)
	keys := []string{"clu.binary", "clu.version", "clu.python.binary", "clu.python.version", "clu.cmdline", "clu.cwd", "clu.user", "clu.date"}
	for _, k := range keys {
		if _, ok := p[k]; !ok {
			t.Errorf("expected key %s in provides map", k)
		}
	}
}

// TestCluParse basic parsing; skips assertions if runtime changes unexpectedly.
func TestCluParse(t *testing.T) {
	c := &Clu{}
	facts := facts.NewFacts()
	c.Parse(facts)
	// Presence
	for _, k := range []string{"clu.binary", "clu.version", "clu.python.version", "clu.cmdline", "clu.cwd", "clu.user", "clu.date"} {
		if !facts.Contains(k) {
			t.Errorf("missing fact %s", k)
		}
	}
	if pyver, _ := facts.Get("clu.python.version"); pyver != runtime.Version() {
		if !strings.Contains(runtime.Version(), pyver) {
			t.Errorf("python version mismatch: got %v want %v", pyver, runtime.Version())
		}
	}
	if dateStr, _ := facts.Get("clu.date"); dateStr != "" {
		if _, err := time.Parse(time.RFC3339, dateStr); err != nil {
			t.Errorf("clu.date not RFC3339: %v", dateStr)
		}
	}
	if cmd, _ := facts.Get("clu.cmdline"); !strings.Contains(cmd, os.Args[0]) {
		t.Errorf("clu.cmdline does not contain argv0: %v", cmd)
	}
}
