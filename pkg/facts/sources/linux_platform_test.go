// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package sources

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/input"
)

func TestPlatformProvides(t *testing.T) {
	provides := types.NewProvides()
	(&Platform{}).Provides(provides)

	provider, ok := provides[platformFactName]
	if !ok {
		t.Fatalf("missing provided fact %q", platformFactName)
	}
	if _, ok := provider.(*Platform); !ok {
		t.Fatalf("provider = %T, want *Platform", provider)
	}
}

func TestPlatformRequiresAllCollectors(t *testing.T) {
	requires := types.NewRequires()
	(&Platform{}).Requires(requires)

	want := []string{virtWhatCommand, systemdDetectVirtCommand}
	if !reflect.DeepEqual(requires.Programs, want) {
		t.Fatalf("programs = %v, want %v", requires.Programs, want)
	}
}

func TestPlatformUsesVirtWhatFirst(t *testing.T) {
	var commands []string
	withPlatformCommandRunner(t, func(command string) (string, int, error) {
		commands = append(commands, command)
		if command != virtWhatCommand {
			t.Fatalf("unexpected fallback command %q", command)
		}
		return "kvm\nvmware", 0, nil
	})

	facts := types.NewFactDB()
	(&Platform{}).Parse(facts)

	assertPlatformFact(t, facts, "kvm, vmware")
	if want := []string{virtWhatCommand}; !reflect.DeepEqual(commands, want) {
		t.Fatalf("commands = %v, want %v", commands, want)
	}
}

func TestPlatformTreatsEmptyVirtWhatAsPhysical(t *testing.T) {
	withPlatformCommandRunner(t, func(command string) (string, int, error) {
		if command != virtWhatCommand {
			t.Fatalf("command = %q", command)
		}
		return "\n\n", 0, nil
	})

	facts := types.NewFactDB()
	(&Platform{}).Parse(facts)
	assertPlatformFact(t, facts, "physical")
}

func TestPlatformReportsFailureWithStubbedFallback(t *testing.T) {
	withPlatformCommandRunner(t, func(command string) (string, int, error) {
		if command != virtWhatCommand {
			t.Fatalf("command = %q", command)
		}
		return "", 1, fmt.Errorf("virt-what unavailable")
	})

	facts := types.NewFactDB()
	(&Platform{}).Parse(facts)
	assertPlatformFact(t, facts, types.ParseFailMsg)
}

func TestPlatformSkipsExistingFact(t *testing.T) {
	withPlatformCommandRunner(t, func(command string) (string, int, error) {
		t.Fatalf("unexpected command %q", command)
		return "", 1, nil
	})

	facts := types.NewFactDB()
	facts.AddFact(types.Fact{Name: platformFactName, Value: "physical", Tier: types.TierOne})
	(&Platform{}).Parse(facts)
	assertPlatformFact(t, facts, "physical")
}

func withPlatformCommandRunner(t *testing.T, runner input.CommandRunnerFunc) {
	t.Helper()
	original := input.CommandRunner
	input.CommandRunner = runner
	t.Cleanup(func() { input.CommandRunner = original })
}

func assertPlatformFact(t *testing.T, facts *types.FactDB, want string) {
	t.Helper()
	got, ok := facts.Get(platformFactName)
	if !ok || got != want {
		t.Fatalf("%q = %q (present: %t), want %q", platformFactName, got, ok, want)
	}
}
