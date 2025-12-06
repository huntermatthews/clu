package types

import (
	"errors"
	"fmt"
)

// Tier represents the priority tier for facts.
// Mirrors the Python Enum values: "one", "two", "three".
type Tier int

const (
	TierInvalid Tier = iota
	TierOne
	TierTwo
	TierThree
)

// String returns the string form equivalent to the Python Enum values.
func (t Tier) String() string {
	switch t {
	case TierOne:
		return "one"
	case TierTwo:
		return "two"
	case TierThree:
		return "three"
	default:
		return "INVALID"
	}
}

// TierFromInt replicates the Python get_by_int() semantics.
// Accepted indexes: 1->TierOne, 2->TierTwo, 3->TierThree.
func TierFromInt(i int) (Tier, error) {
	switch i {
	case 1:
		return TierOne, nil
	case 2:
		return TierTwo, nil
	case 3:
		return TierThree, nil
	default:
		return TierInvalid, errors.New("invalid tier index")
	}
}

// Facts stores key/value facts and tracks which keys were added at which tier.
// Behavior mirrors the Python Facts class.
type Facts struct {
	facts map[string]string
	tier  map[Tier][]string
}

// NewFacts constructs a new Facts instance.
func NewFacts() *Facts {
	return &Facts{
		facts: make(map[string]string),
		tier: map[Tier][]string{
			TierOne:   {},
			TierTwo:   {},
			TierThree: {},
		},
	}
}

// Add adds a key/value at the specified priority tier.
func (f *Facts) Add(priority Tier, key, value string) {
	f.tier[priority] = append(f.tier[priority], key)
	f.facts[key] = value
}

// Set emulates Python __setitem__: assigns value and records key as TierOne.
func (f *Facts) Set(key, value string) {
	f.facts[key] = value
	f.tier[TierOne] = append(f.tier[TierOne], key)
}

// Get returns the value and whether it was present.
func (f *Facts) Get(key string) (string, bool) {
	v, ok := f.facts[key]
	return v, ok
}

// Contains reports whether the key exists.
func (f *Facts) Contains(key string) bool {
	_, ok := f.facts[key]
	return ok
}

// Update merges another map into the facts, recording each key as TierOne (Python behavior).
func (f *Facts) Update(other map[string]string) {
	for k, v := range other {
		f.Set(k, v)
	}
}

// GetTier returns keys visible at the requested tier.
// TierOne -> TierOne keys
// TierTwo -> TierOne + TierTwo keys
// TierThree -> TierOne + TierTwo + TierThree keys
func (f *Facts) GetTier(t Tier) []string {
	var result []string
	switch t {
	case TierOne:
		result = append(result, f.tier[TierOne]...)
	case TierTwo:
		result = append(result, f.tier[TierOne]...)
		result = append(result, f.tier[TierTwo]...)
	case TierThree:
		result = append(result, f.tier[TierOne]...)
		result = append(result, f.tier[TierTwo]...)
		result = append(result, f.tier[TierThree]...)
	default:
		// invalid -> empty
	}
	return result
}

// ToMap returns the underlying map (read-only by convention). Modifying it will affect Facts.
func (f *Facts) ToMap() map[string]string {
	return f.facts
}

// Keys returns all fact keys (unordered).
func (f *Facts) Keys() []string {
	keys := make([]string, 0, len(f.facts))
	for k := range f.facts {
		keys = append(keys, k)
	}
	return keys
}

// Equals compares with another Facts or a plain map (Python __eq__ semantics).
func (f *Facts) Equals(other interface{}) bool {
	switch o := other.(type) {
	case *Facts:
		if len(f.facts) != len(o.facts) {
			return false
		}
		for k, v := range f.facts {
			if ov, ok := o.facts[k]; !ok || ov != v {
				return false
			}
		}
		// Tier comparison omitted (mirrors Python FIXME comment)
		return true
	case map[string]string:
		if len(f.facts) != len(o) {
			return false
		}
		for k, v := range f.facts {
			if ov, ok := o[k]; !ok || ov != v {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// String returns a representation similar to Python __str__/__repr__.
func (f *Facts) String() string {
	return fmt.Sprintf("Facts: %v, %v", f.facts, f.tier)
}
