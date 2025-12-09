package sources

import (
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/facts/types"
)

// Clu provides runtime-related facts about the running CLI and environment.
type Clu struct{}

func (c *Clu) Provides(p types.Provides) {
	for _, k := range []string{
		"clu.binary",
		"clu.version",
		"clu.python.binary",
		"clu.python.version",
		"clu.cmdline",
		"clu.cwd",
		"clu.user",
		"clu.date",
	} {
		p[k] = c
	}
}

func (c *Clu) Requires(r *types.Requires) { /* no external requirements */ }

func (c *Clu) Parse(f *types.Facts) {
	// Binary name (argv0)
	argv0 := os.Args[0]
	f.Add(types.TierTwo, "clu.binary", argv0)

	// Version from about.go
	f.Add(types.TierOne, "clu.version", pkg.Version)

	// Python analogs: in Go, provide placeholders: binary is the current executable path
	// and version is runtime.Version()
	f.Add(types.TierTwo, "clu.python.binary", argv0)
	f.Add(types.TierOne, "clu.python.version", runtime.Version())

	// Command line
	f.Add(types.TierTwo, "clu.cmdline", strings.Join(os.Args, " "))

	// Working directory
	cwd, _ := os.Getwd()
	f.Add(types.TierThree, "clu.cwd", cwd)
	// User
	u, _ := user.Current()
	if u != nil && u.Username != "" {
		f.Add(types.TierThree, "clu.user", u.Username)
	} else {
		f.Add(types.TierThree, "clu.user", "")
	}

	// RFC3339 timestamp (UTC)
	f.Add(types.TierTwo, "clu.date", time.Now().UTC().Format(time.RFC3339))
}
