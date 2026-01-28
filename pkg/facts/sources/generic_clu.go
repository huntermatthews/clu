package sources

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/global"
)

// Clu provides runtime-related facts about the running CLI and environment.
type Clu struct{}

var cluFacts = map[string]*types.Fact{
	"clu.binary":         {Name: "clu.binary", Tier: types.TierTwo},
	"clu.version":        {Name: "clu.version", Tier: types.TierOne},
	"clu.golang.version": {Name: "clu.golang.version", Tier: types.TierThree},
	"clu.cmdline":        {Name: "clu.cmdline", Tier: types.TierTwo},
	"clu.cwd":            {Name: "clu.cwd", Tier: types.TierThree},
	"clu.user":           {Name: "clu.user", Tier: types.TierThree},
	"clu.date":           {Name: "clu.date", Tier: types.TierTwo},
}

func (c *Clu) Provides(p types.Provides) {
	for name := range cluFacts {
		p[name] = c
	}
}

func (c *Clu) Requires(r *types.Requires) { /* no external requirements */ }

func (c *Clu) Parse(f *types.FactDB) {
	// Binary name (argv0)
	argv0 := os.Args[0]
	argv0, _ = filepath.Abs(argv0)
	cluFacts["clu.binary"].Value = argv0

	// Version from about.go
	cluFacts["clu.version"].Value = global.Version

	// Go runtime version
	cluFacts["clu.golang.version"].Value = runtime.Version()

	// Command line
	cluFacts["clu.cmdline"].Value = strings.Join(os.Args, " ")

	// Working directory
	cwd, _ := os.Getwd()
	cluFacts["clu.cwd"].Value = cwd

	// User
	u, _ := user.Current()
	if u != nil {
		cluFacts["clu.user"].Value = u.Username
	} else {
		cluFacts["clu.user"].Value = types.ParseFailMsg
	}

	// RFC3339 timestamp (UTC)
	cluFacts["clu.date"].Value = time.Now().UTC().Format(time.RFC3339)

	// Add all facts to the FactDB
	for _, fact := range cluFacts {
		f.AddFact(*fact)
	}
}
