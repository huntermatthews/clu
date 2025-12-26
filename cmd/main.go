package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/subcmd"
)

// CLI defines the root command and global flags.
type CLI struct {
	Debug     bool                `help:"Enable debug logging."`
	Net       bool                `name:"net" help:"Enable network access."`
	MockDir   string              `help:"Enable mock mode for testing." hidden:""`
	Version   kong.VersionFlag    `help:"Print version information and quit."`
	Facts     subcmd.FactsCmd     `cmd:"" help:"Show facts." default:"withargs"`
	Collector subcmd.CollectorCmd `cmd:"" help:"Run collector."`
	Requires  subcmd.RequiresCmd  `cmd:"" help:"Requires actions: list or check."`
}

func main() {
	cli := &CLI{}
	k, err := kong.New(cli,
		kong.Name("clu"),
		kong.Description("Kong example with facts/collector/requires subcommands."),
		kong.UsageOnError(),
		kong.Vars{"version": "clu " + pkg.Version},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to init CLI:", err)
		os.Exit(2)
	}

	ctx, err := k.Parse(os.Args[1:])
	if err != nil {
		k.FatalIfErrorf(err)
	}

	if cli.Debug {
		fmt.Fprintln(os.Stderr, "debug: enabled")
		// see comment in pkg/config.go about this
		pkg.CluConfig.Debug = true
	}

	if cli.Net {
		fmt.Fprintln(os.Stderr, "net: enabled")
		pkg.CluConfig.NetEnabled = true
	}

	if cli.MockDir != "" {
		fmt.Fprintln(os.Stderr, "mock mode: enabled, dir =", cli.MockDir)
		pkg.CluConfig.MockDir = cli.MockDir
	}

	err = ctx.Run()
	ctx.FatalIfErrorf(err)
}
