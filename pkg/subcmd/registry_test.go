// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package subcmd

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestBuiltinCommandsRegister(t *testing.T) {
	parser, err := kong.New(&struct{}{}, KongOptions()...)
	if err != nil {
		t.Fatalf("creating parser: %v", err)
	}

	context, err := parser.Parse(nil)
	if err != nil {
		t.Fatalf("parsing bare command: %v", err)
	}
	if context.Command() != "facts" {
		t.Fatalf("bare command = %q, want %q", context.Command(), "facts")
	}

	commands := [][]string{
		{"facts"},
		{"collector"},
		{"requires", "list"},
		{"tools"},
		{"version"},
		{"check"},
		{"help"},
	}
	if len(KongOptions()) != len(commands) {
		t.Fatalf("registered command count = %d, want %d", len(KongOptions()), len(commands))
	}

	for _, args := range commands {
		if _, err := parser.Parse(args); err != nil {
			t.Errorf("parsing %q: %v", args, err)
		}
	}
}
