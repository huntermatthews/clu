package pkg

// config2.go provides a simple global configuration container with NO thread safety.
// Fields:
//  - Verbose: enable extra (non-debug) informational output
//  - Debug: enable diagnostic / debug output
//  - Output: selected output format name (e.g., "dots", "json", "shell")
//  - Facts: list of fact name prefixes requested by user / caller
//
// NOTE: This version intentionally omits any synchronization primitives. It is
// not safe for concurrent mutation. Callers in multi-goroutine contexts should
// add their own coordination if needed.

// GlobalConfig holds runtime configuration flags and selections.
type GlobalConfig struct {
	Verbose bool
	Debug   bool
	Output  string
	Facts   []string
}

// config2 is the singleton instance backing global configuration access.
var config2 = &GlobalConfig{}

// Config returns the global configuration instance for direct method use.
// Config2 returns the secondary global configuration instance.
func Config2() *GlobalConfig { return config2 }

// Configure applies a set of functional options atomically.
func Configure(opts ...Option) {
	for _, opt := range opts {
		opt(config2)
	}
}

// Option represents a functional configuration mutation.
type Option func(*GlobalConfig)

// WithVerbose sets Verbose.
func WithVerbose(v bool) Option { return func(c *GlobalConfig) { c.Verbose = v } }

// WithDebug sets Debug.
func WithDebug(d bool) Option { return func(c *GlobalConfig) { c.Debug = d } }

// WithOutput sets Output.
func WithOutput(o string) Option { return func(c *GlobalConfig) { c.Output = o } }

// WithFacts replaces the Facts slice with a copy.
func WithFacts(facts []string) Option {
	cp := make([]string, len(facts))
	copy(cp, facts)
	return func(c *GlobalConfig) { c.Facts = cp }
}

// SetVerbose updates Verbose.
func (c *GlobalConfig) SetVerbose(v bool) { c.Verbose = v }

// IsVerbose returns current Verbose state.
func (c *GlobalConfig) IsVerbose() bool { return c.Verbose }

// SetDebug updates Debug.
func (c *GlobalConfig) SetDebug(d bool) { c.Debug = d }

// IsDebug returns current Debug state.
func (c *GlobalConfig) IsDebug() bool { return c.Debug }

// SetOutput sets Output format.
func (c *GlobalConfig) SetOutput(o string) { c.Output = o }

// GetOutput returns Output format.
func (c *GlobalConfig) GetOutput() string { return c.Output }

// SetFacts replaces Facts slice (defensive copy).
func (c *GlobalConfig) SetFacts(facts []string) {
	cp := make([]string, len(facts))
	copy(cp, facts)
	c.Facts = cp
}

// AddFact appends a fact if not already present.
func (c *GlobalConfig) AddFact(f string) {
	if f == "" {
		return
	}
	for _, existing := range c.Facts {
		if existing == f {
			return
		}
	}
	c.Facts = append(c.Facts, f)
}

// HasFact reports whether the fact is present.
func (c *GlobalConfig) HasFact(f string) bool {
	for _, existing := range c.Facts {
		if existing == f {
			return true
		}
	}
	return false
}

// GetFacts returns a copy of the Facts slice.
func (c *GlobalConfig) GetFacts() []string {
	cp := make([]string, len(c.Facts))
	copy(cp, c.Facts)
	return cp
}

// Reset clears all fields to zero values.
func (c *GlobalConfig) Reset() {
	c.Verbose = false
	c.Debug = false
	c.Output = ""
	c.Facts = nil
}
