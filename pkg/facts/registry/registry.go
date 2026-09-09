// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package registry

import (
	"sync"

	"github.com/huntermatthews/clu/pkg/facts/types"
)

// SourceFactory creates a fact source.
type SourceFactory func() types.Sources

var (
	mutex   sync.RWMutex
	sources = make(map[string][]SourceFactory)
)

// Register adds a source factory for each supported operating system.
func Register(factory SourceFactory, systems ...string) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, system := range systems {
		sources[system] = append(sources[system], factory)
	}
}

// GetSources creates all sources registered for an operating system.
func GetSources(system string) []types.Sources {
	mutex.RLock()
	factories := append([]SourceFactory(nil), sources[system]...)
	mutex.RUnlock()

	result := make([]types.Sources, 0, len(factories))
	for _, factory := range factories {
		result = append(result, factory())
	}

	return result
}
