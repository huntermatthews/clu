package global

// NOTE: I'm not super happy with this global mutable state design, but it's
// simple and effective for now. We can refactor later if needed.
// Specifically, I want a _simple_ way to access Net and Debug at least in the Source.Parse() methods that
// doesn't involve passing config objects/values through multiple layers of calls.
type config struct {
	Debug      bool
	Verbose    bool
	NetEnabled bool
	MockDir    string
}

// Config is the singleton instance backing global configuration access.
var Config = &config{}
