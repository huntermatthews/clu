package pkg

// Port of src/clu/config.py which provides a process-wide mutable configuration
// object. Python version uses argparse.Namespace with dynamic attribute setting.
// In Go we implement a simple map-based singleton with Set and Get semantics.

import "sync"

// Config is a key/value store for string -> interface{} configuration values.
// Using interface{} mirrors Python's dynamic namespace; callers can type assert.
// If stricter typing becomes desirable, replace with a struct.
type Config struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

var globalConfig = &Config{data: make(map[string]interface{})}

// SetConfig replaces or adds all key/value pairs from the provided map into the global config.
func SetConfig(cfg map[string]interface{}) {
	globalConfig.mu.Lock()
	for k, v := range cfg {
		globalConfig.data[k] = v
	}
	globalConfig.mu.Unlock()
}

// GetConfig returns a snapshot copy of current configuration map.
func GetConfig() map[string]interface{} {
	globalConfig.mu.RLock()
	copy := make(map[string]interface{}, len(globalConfig.data))
	for k, v := range globalConfig.data {
		copy[k] = v
	}
	globalConfig.mu.RUnlock()
	return copy
}

// Set sets a single key/value.
func (c *Config) Set(key string, value interface{}) {
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
}

// Get retrieves a value and presence flag.
func (c *Config) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	v, ok := c.data[key]
	c.mu.RUnlock()
	return v, ok
}

// Global returns the singleton Config pointer (for advanced usage).
func Global() *Config { return globalConfig }
