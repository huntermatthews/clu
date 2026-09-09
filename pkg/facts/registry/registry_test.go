// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package registry

import (
	"testing"

	"github.com/huntermatthews/clu/pkg/facts/types"
)

type testSource struct{}

func (source *testSource) Provides(types.Provides) {}

func (source *testSource) Requires(*types.Requires) {}

func (source *testSource) Parse(*types.FactDB) {}

func TestGetSources(t *testing.T) {
	Register(func() types.Sources { return &testSource{} }, "test")
	Register(func() types.Sources { return &testSource{} }, "test")

	first := GetSources("test")
	if len(first) != 2 {
		t.Fatalf("source count = %d, want 2", len(first))
	}

	second := GetSources("test")
	if first[0] == second[0] {
		t.Error("GetSources returned a reused source instance")
	}

	if sources := GetSources("unknown"); len(sources) != 0 {
		t.Errorf("unknown OS source count = %d, want 0", len(sources))
	}
}
