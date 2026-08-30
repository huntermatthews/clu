// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package types

import (
	"fmt"
)

// Tier represents the priority tier for facts.
type Tier int

const (
	TierInvalid Tier = iota
	TierOne
	TierTwo
	TierThree
)

// A single Fact
type Fact struct {
	Name   string
	Value  string
	Origin string
	Tier   Tier
}

// FactDB stores key/value facts and tracks which keys were added at which tier.
// Behavior mirrors the Python FactDB class.
type FactDB struct {
	facts map[string]string
	tier  map[Tier][]string
}

// NewFactDB constructs a new FactDB instance.
func NewFactDB() *FactDB {
	return &FactDB{
		facts: make(map[string]string),
		tier: map[Tier][]string{
			TierOne:   {},
			TierTwo:   {},
			TierThree: {},
		},
	}
}

// AddFact adds a Fact to the FactDB using its Tier, Name, and Value.
func (f *FactDB) AddFact(fact Fact) {
	if _, exists := f.facts[fact.Name]; !exists {
		f.tier[fact.Tier] = append(f.tier[fact.Tier], fact.Name)
	}
	f.facts[fact.Name] = fact.Value
}

// AddAll adds every fact in the map to the FactDB.
func (f *FactDB) AddAll(facts map[string]*Fact) {
	for _, fact := range facts {
		f.AddFact(*fact)
	}
}

// Get returns the value and whether it was present.
func (f *FactDB) Get(key string) (string, bool) {
	v, ok := f.facts[key]
	return v, ok
}

// Contains reports whether the key exists.
func (f *FactDB) Contains(key string) bool {
	_, ok := f.facts[key]
	return ok
}

// GetTier returns keys visible at the requested tier.
// TierOne -> TierOne keys
// TierTwo -> TierOne + TierTwo keys
// TierThree -> TierOne + TierTwo + TierThree keys
func (f *FactDB) GetTier(t Tier) []string {
	var result []string
	for tier := TierOne; tier <= t; tier++ {
		result = append(result, f.tier[tier]...)
	}
	return result
}

// String returns a representation similar to Python __str__/__repr__.
func (f *FactDB) String() string {
	return fmt.Sprintf("Facts: %v, %v", f.facts, f.tier)
}
