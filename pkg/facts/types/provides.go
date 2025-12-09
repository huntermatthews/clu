package types


// Provides is a minimal container analogous to the Python subclass of dict.
// In Python: class Provides(dict): pass
// In Go we just use a map with helper constructor for parity.
// Keys are strings; values are generic (interface{}) to allow flexibility
// similar to Python's dynamic nature.

type Provides map[string]Sources

// NewProvides returns an empty Provides map.
func NewProvides() Provides {
	return make(Provides)
}

// Add inserts or replaces a key/value pair.
func (p Provides) Add(key string, value Sources) {
	p[key] = value
}

// Get retrieves a value and a bool indicating presence.
func (p Provides) Get(key string) (Sources, bool) {
	v, ok := p[key]
	return v, ok
}

// Keys returns all keys in the Provides container.
func (p Provides) Keys() []string {
	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	return keys
}
