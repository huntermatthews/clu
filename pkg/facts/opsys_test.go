// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package facts

import (
	"slices"
	"testing"

	"github.com/huntermatthews/clu/pkg/input"
)

func TestOpSysFactoryUsesMockedUname(t *testing.T) {
	originalChecker := input.ProgramChecker
	originalRunner := input.CommandRunner
	defer func() {
		input.ProgramChecker = originalChecker
		input.CommandRunner = originalRunner
	}()

	testCases := []struct {
		name       string
		kernel     string
		earlyFacts []string
	}{
		{
			name:       "linux",
			kernel:     "Linux\n",
			earlyFacts: []string{"phy.arch", "phy.platform"},
		},
		{
			name:       "darwin",
			kernel:     "Darwin\n",
			earlyFacts: []string{"os.version"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			input.ProgramChecker = func(string) string { return "" }
			input.CommandRunner = func(command string) (string, int, error) {
				if command != "uname" {
					t.Errorf("command = %q, want %q", command, "uname")
				}
				return testCase.kernel, 0, nil
			}

			result := OpSysFactory()
			if !slices.Equal(result.GetEarlyFacts(), testCase.earlyFacts) {
				t.Errorf("early facts = %v, want %v", result.GetEarlyFacts(), testCase.earlyFacts)
			}
		})
	}
}

func TestBuiltInSourcesRegister(t *testing.T) {
	testCases := []struct {
		name   string
		result *OpSys
		count  int
	}{
		{name: "linux", result: NewLinux(), count: 14},
		{name: "darwin", result: NewDarwin(), count: 5},
		{name: "windows", result: NewWindows(), count: 2},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if len(testCase.result.Sources) != testCase.count {
				t.Errorf("source count = %d, want %d", len(testCase.result.Sources), testCase.count)
			}
		})
	}
}
