package sources

import (
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	pkg "github.com/huntermatthews/clu/pkg"
)

// Clu provides runtime-related facts about the running CLI and environment.
type Clu struct{}

func (c *Clu) Provides(p pkg.Provides) {
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

func (c *Clu) Requires(r *pkg.Requires) { /* no external requirements */ }

func (c *Clu) Parse(f *pkg.Facts) {
	// Binary name (argv0)
	argv0 := os.Args[0]
	f.Add(pkg.TierTwo, "clu.binary", argv0)

	// Version from about.go
	f.Add(pkg.TierOne, "clu.version", pkg.Version)

	// Python analogs: in Go, provide placeholders: binary is the current executable path
	// and version is runtime.Version()
	f.Add(pkg.TierTwo, "clu.python.binary", argv0)
	f.Add(pkg.TierOne, "clu.python.version", runtime.Version())

	// Command line
	f.Add(pkg.TierTwo, "clu.cmdline", strings.Join(os.Args, " "))

	// Working directory
	cwd, _ := os.Getwd()
	f.Add(pkg.TierThree, "clu.cwd", cwd)

	// User
	u, _ := user.Current()
	if u != nil && u.Username != "" {
		f.Add(pkg.TierThree, "clu.user", u.Username)
	} else {
		f.Add(pkg.TierThree, "clu.user", "")
	}

	// RFC3339 timestamp (UTC)
	f.Add(pkg.TierTwo, "clu.date", time.Now().UTC().Format(time.RFC3339))
}
