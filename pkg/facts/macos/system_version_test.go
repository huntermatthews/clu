// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package macos

import (
	"errors"
	"reflect"
	"testing"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

const sampleSystemVersionPlist = `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0"><dict>
<key>ProductName</key><string>macOS</string>
<key>ProductVersion</key><string>15.4</string>
<key>ProductBuildVersion</key><string>24E248</string>
<key>BuildID</key><string>ABC123</string>
</dict></plist>`

const sampleSwVers = `ProductName: macOS
ProductVersion: 15.4
BuildVersion: 24E248`

func TestSystemVersionProvides(t *testing.T) {
	provides := types.NewProvides()
	(&SystemVersion{}).Provides(provides)

	for _, name := range []string{"os.name", "os.version", "os.build", "id.build_id"} {
		provider, ok := provides[name]
		if !ok {
			t.Errorf("missing provided fact %q", name)
			continue
		}
		if _, ok := provider.(*SystemVersion); !ok {
			t.Errorf("provider for %q = %T, want *SystemVersion", name, provider)
		}
	}
}

func TestSystemVersionRequiresAllCollectors(t *testing.T) {
	requires := types.NewRequires()
	(&SystemVersion{}).Requires(requires)

	if want := []string{systemVersionPlistPath}; !reflect.DeepEqual(requires.Files, want) {
		t.Fatalf("files = %v, want %v", requires.Files, want)
	}
	if want := []string{swVersCommand}; !reflect.DeepEqual(requires.Programs, want) {
		t.Fatalf("programs = %v, want %v", requires.Programs, want)
	}
}

func TestSystemVersionUsesPlistBeforeSwVers(t *testing.T) {
	withFileReader(t, func(path string) (string, error) {
		if path != systemVersionPlistPath {
			t.Fatalf("file path = %q", path)
		}
		return sampleSystemVersionPlist, nil
	})
	withCommandRunner(t, func(command string) (string, int, error) {
		t.Fatalf("unexpected fallback command %q", command)
		return "", 1, nil
	})

	facts := types.NewFactDB()
	(&SystemVersion{}).Parse(facts)

	assertDarwinSystemVersionFacts(t, facts, map[string]string{
		"os.name":     "macOS",
		"os.version":  "15.4",
		"os.build":    "24E248",
		"id.build_id": "ABC123",
	})
}

func TestSystemVersionDoesNotRepeatACompleteParse(t *testing.T) {
	var reads int
	withFileReader(t, func(string) (string, error) {
		reads++
		return sampleSystemVersionPlist, nil
	})
	withCommandRunner(t, func(command string) (string, int, error) {
		t.Fatalf("unexpected fallback command %q", command)
		return "", 1, nil
	})

	facts := types.NewFactDB()
	source := &SystemVersion{}
	source.Parse(facts)
	source.Parse(facts)

	if reads != 1 {
		t.Fatalf("plist reads = %d, want 1", reads)
	}
}

func TestSystemVersionFallsBackToSwVers(t *testing.T) {
	withFileReader(t, func(string) (string, error) {
		return "", errors.New("plist unavailable")
	})
	var commands []string
	withCommandRunner(t, func(command string) (string, int, error) {
		commands = append(commands, command)
		if command != swVersCommand {
			t.Fatalf("command = %q", command)
		}
		return sampleSwVers, 0, nil
	})

	facts := types.NewFactDB()
	(&SystemVersion{}).Parse(facts)

	assertDarwinSystemVersionFacts(t, facts, map[string]string{
		"os.name":     "macOS",
		"os.version":  "15.4",
		"os.build":    "24E248",
		"id.build_id": types.ParseFailMsg,
	})
	if want := []string{swVersCommand}; !reflect.DeepEqual(commands, want) {
		t.Fatalf("commands = %v, want %v", commands, want)
	}
}

func TestSystemVersionFallsBackWhenPlistIsIncomplete(t *testing.T) {
	withFileReader(t, func(string) (string, error) {
		return "<key>ProductName</key><string>macOS</string>", nil
	})
	withCommandRunner(t, func(string) (string, int, error) {
		return sampleSwVers, 0, nil
	})

	facts := types.NewFactDB()
	(&SystemVersion{}).Parse(facts)

	assertDarwinSystemVersionFacts(t, facts, map[string]string{
		"os.name":     "macOS",
		"os.version":  "15.4",
		"os.build":    "24E248",
		"id.build_id": types.ParseFailMsg,
	})
}

func TestSystemVersionReportsFailureWhenAllCollectorsFail(t *testing.T) {
	withFileReader(t, func(string) (string, error) {
		return "", errors.New("plist unavailable")
	})
	withCommandRunner(t, func(string) (string, int, error) {
		return "", 1, errors.New("sw_vers unavailable")
	})

	facts := types.NewFactDB()
	(&SystemVersion{}).Parse(facts)

	assertDarwinSystemVersionFacts(t, facts, map[string]string{
		"os.name":     types.ParseFailMsg,
		"os.version":  types.ParseFailMsg,
		"os.build":    types.ParseFailMsg,
		"id.build_id": types.ParseFailMsg,
	})
}

func withFileReader(t *testing.T, reader input.FileReaderFunc) {
	t.Helper()
	original := input.FileReader
	input.FileReader = reader
	t.Cleanup(func() { input.FileReader = original })
}

func withCommandRunner(test *testing.T, runner input.CommandRunnerFunc) {
	test.Helper()
	original := input.CommandRunner
	input.CommandRunner = runner
	test.Cleanup(func() { input.CommandRunner = original })
}

func assertDarwinSystemVersionFacts(t *testing.T, facts *types.FactDB, want map[string]string) {
	t.Helper()
	for name, wantValue := range want {
		got, ok := facts.Get(name)
		if !ok || got != wantValue {
			t.Errorf("%q = %q (present: %t), want %q", name, got, ok, wantValue)
		}
	}
}
