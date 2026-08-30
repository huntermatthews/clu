// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package sources

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/huntermatthews/clu/internal/facts/types"
	"github.com/huntermatthews/clu/internal/input"
)

const sampleLsmem = `RANGE   SIZE   STATE
0-3     4G     online
Total online memory:     4294967296`

func TestMemoryProvides(t *testing.T) {
	src := &Memory{}
	provides := types.NewProvides()
	src.Provides(provides)

	provider, ok := provides[memoryFactName]
	if !ok {
		t.Fatalf("missing provided fact %q", memoryFactName)
	}
	if _, ok := provider.(*Memory); !ok {
		t.Fatalf("provider = %T, want *Memory", provider)
	}
}

func TestMemoryRequiresAllCollectors(t *testing.T) {
	requires := types.NewRequires()
	(&Memory{}).Requires(requires)

	want := []string{lsmemCommand, udevadmCommand}
	if !reflect.DeepEqual(requires.Programs, want) {
		t.Fatalf("programs = %v, want %v", requires.Programs, want)
	}
}

func TestLsmemParse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		withCommandRunner(t, func(cmdline string) (string, int, error) {
			if cmdline != lsmemCommand {
				t.Fatalf("command = %q", cmdline)
			}
			return sampleLsmem, 0, nil
		})

		facts := types.NewFactDB()
		(&lsmem{}).Parse(facts)
		assertFactValue(t, facts, input.BytesToSI(4294967296))
	})

	for _, tc := range []struct {
		name string
		data string
		rc   int
	}{
		{name: "command failure", rc: 1},
		{name: "missing total", data: "Header\nNo total line here"},
		{name: "zero total", data: "Total online memory: 0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			withCommandRunner(t, func(string) (string, int, error) {
				return tc.data, tc.rc, fmt.Errorf("unavailable")
			})

			facts := types.NewFactDB()
			(&lsmem{}).Parse(facts)
			assertFactValue(t, facts, types.ParseFailMsg)
		})
	}
}

func TestUdevadmRamParse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		withCommandRunner(t, func(cmdline string) (string, int, error) {
			if cmdline != udevadmCommand {
				t.Fatalf("command = %q", cmdline)
			}
			return "MEMORY_DEVICE_0_SIZE=1048576\nMEMORY_DEVICE_1_SIZE=2097152\n", 0, nil
		})

		facts := types.NewFactDB()
		(&udevadmRam{}).Parse(facts)
		assertFactValue(t, facts, input.BytesToSI(3145728))
	})

	for _, tc := range []struct {
		name string
		data string
		rc   int
	}{
		{name: "command failure", rc: 1},
		{name: "no memory devices", data: "OTHER=1"},
		{name: "zero-sized memory devices", data: "MEMORY_DEVICE_0_SIZE=0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			withCommandRunner(t, func(string) (string, int, error) {
				return tc.data, tc.rc, fmt.Errorf("unavailable")
			})

			facts := types.NewFactDB()
			(&udevadmRam{}).Parse(facts)
			assertFactValue(t, facts, types.ParseFailMsg)
		})
	}
}

func TestMemoryUsesLsmemBeforeFallback(t *testing.T) {
	var calls []string
	withCommandRunner(t, func(cmdline string) (string, int, error) {
		calls = append(calls, cmdline)
		if cmdline != lsmemCommand {
			t.Fatalf("unexpected fallback command %q", cmdline)
		}
		return sampleLsmem, 0, nil
	})

	facts := types.NewFactDB()
	(&Memory{}).Parse(facts)

	assertFactValue(t, facts, input.BytesToSI(4294967296))
	if want := []string{lsmemCommand}; !reflect.DeepEqual(calls, want) {
		t.Fatalf("commands = %v, want %v", calls, want)
	}
}

func TestMemoryFallsBackToUdevadm(t *testing.T) {
	var calls []string
	withCommandRunner(t, func(cmdline string) (string, int, error) {
		calls = append(calls, cmdline)
		switch cmdline {
		case lsmemCommand:
			return "", 1, fmt.Errorf("lsmem unavailable")
		case udevadmCommand:
			return "MEMORY_DEVICE_0_SIZE=2147483648", 0, nil
		default:
			t.Fatalf("unexpected command %q", cmdline)
			return "", 1, nil
		}
	})

	facts := types.NewFactDB()
	(&Memory{}).Parse(facts)

	assertFactValue(t, facts, input.BytesToSI(2147483648))
	want := []string{lsmemCommand, udevadmCommand}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("commands = %v, want %v", calls, want)
	}

	var ramFactCount int
	for _, key := range facts.GetTier(types.TierOne) {
		if key == memoryFactName {
			ramFactCount++
		}
	}
	if ramFactCount != 1 {
		t.Fatalf("%q appears %d times in TierOne facts, want once", memoryFactName, ramFactCount)
	}
}

func TestMemoryReportsFailureWhenAllCandidatesFail(t *testing.T) {
	withCommandRunner(t, func(string) (string, int, error) {
		return "", 1, fmt.Errorf("unavailable")
	})

	facts := types.NewFactDB()
	(&Memory{}).Parse(facts)
	assertFactValue(t, facts, types.ParseFailMsg)
}

func withCommandRunner(t *testing.T, runner input.CommandRunnerFunc) {
	t.Helper()
	original := input.CommandRunner
	input.CommandRunner = runner
	t.Cleanup(func() { input.CommandRunner = original })
}

func assertFactValue(t *testing.T, facts *types.FactDB, want string) {
	t.Helper()
	got, ok := facts.Get(memoryFactName)
	if !ok || got != want {
		t.Fatalf("%q = %q (present: %t), want %q", memoryFactName, got, ok, want)
	}
}
