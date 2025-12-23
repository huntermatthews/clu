package pkg

type config struct {
	Debug      bool
	Verbose    bool
	NetEnabled bool
}

// Config is the singleton instance backing global configuration access.
var CluConfig = &config{}
