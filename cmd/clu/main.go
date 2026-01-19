package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"

	"github.com/NHGRI/clu/pkg/global"
	"github.com/NHGRI/clu/pkg/input"
	"github.com/NHGRI/clu/pkg/subcmd"
)

// CLI defines the root command and global flags.
type CLI struct {
	Debug   bool             `help:"Enable debug logging."`
	Net     bool             `name:"net" help:"Enable network access."`
	MockDir string           `help:"Enable mock mode for testing." hidden:""`
	Version kong.VersionFlag `help:"Print version information and quit."`

	Facts     subcmd.FactsCmd     `cmd:"" help:"Show facts." default:"withargs"`
	Collector subcmd.CollectorCmd `cmd:"" help:"Run collector."`
	Requires  subcmd.RequiresCmd  `cmd:"" help:"Requires actions: list or check."`
	Tools     subcmd.ToolsCmd     `cmd:"" help:"Check SSA admins workstations tools."`
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr input.Stdout) int {
	cli := &CLI{}
	k, err := kong.New(cli,
		kong.Name("clu"),
		kong.Description("Kong example with facts/collector/requires subcommands."),
		kong.UsageOnError(),
		kong.Vars{"version": "clu " + global.Version},
		kong.Writers(stdout, stderr),
		kong.BindTo(stdout, (*input.Stdout)(nil)),
		kong.BindTo(stderr, (*input.Stderr)(nil)),
	)
	if err != nil {
		fmt.Fprintln(stderr, "failed to init CLI:", err)
		return 2
	}

	ctx, err := k.Parse(args)
	if err != nil {
		k.Errorf("%v", err)
		return 1
	}

	// setup logging and global config based on debug bool flag
	setupLogging(cli.Debug)
	global.Config.Debug = cli.Debug

	if cli.Net {
		slog.Debug("network access enabled")
		global.Config.NetEnabled = true
	}

	if cli.MockDir != "" {
		slog.Debug("mock mode: enabled, dir =" + cli.MockDir)

		if err := EnableMockMode(cli.MockDir); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}

	err = ctx.Run()
	if err != nil {
		k.Errorf("%v", err)
		return 1
	}
	return 0
}

func EnableMockMode(dir string) error {
	path, err := resolveMockPath(dir)
	if err != nil {
		return err
	}

	global.Config.MockDir = path

	// Change the input functions to the mock versions.
	input.CommandRunner = input.MockTextProgram
	input.FileReader = input.MockTextFile
	input.ProgramChecker = input.MockProgramCheck
	return nil
}

func resolveMockPath(dir string) (string, error) {
	if filepath.IsAbs(dir) {
		return dir, nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// Search for testdata/<dir> starting from CWD and walking up
	searchDir := pwd
	for {
		path := filepath.Join(searchDir, "testdata", dir)
		info, err := os.Stat(path)
		if err == nil {
			if info.IsDir() {
				return path, nil
			}
			return "", fmt.Errorf("mock path is not a directory: %s", path)
		}

		parent := filepath.Dir(searchDir)
		if parent == searchDir {
			break
		}
		searchDir = parent
	}
	return "", fmt.Errorf("mock directory not found: %s", dir)
}
