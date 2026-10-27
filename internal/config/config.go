// Package config loads and validates process configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"slices"
	"strconv"
)

// Environment variable names read by Load.
const (
	EnvPort        = "SEARCH_PORT"
	EnvDatabaseURL = "SEARCH_DATABASE_URL"
	EnvLogLevel    = "SEARCH_LOG_LEVEL"
)

// Defaults applied when the corresponding variable is unset or empty.
const (
	DefaultPort     = 8080
	DefaultLogLevel = "info"
)

const (
	minPort = 1
	maxPort = 65535
)

var validLogLevels = []string{"debug", "info", "warn", "error"}

// Config holds validated settings. Construct it with Load; a hand-built
// Config has not been validated.
type Config struct {
	// Port is the TCP port the HTTP server listens on (1-65535).
	Port int
	// DatabaseURL is the Postgres connection string.
	DatabaseURL string
	// LogLevel is one of "debug", "info", "warn" or "error".
	LogLevel string
}

// Load reads configuration from the environment, applies defaults and
// validates the result. It returns an error naming the offending variable
// if any setting is invalid.
func Load() (Config, error) {
	cfg := Config{Port: DefaultPort, LogLevel: DefaultLogLevel}

	if v := os.Getenv(EnvPort); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("%s must be an integer: %w", EnvPort, err)
		}
		if port < minPort || port > maxPort {
			return Config{}, fmt.Errorf("%s must be between %d and %d, got %d", EnvPort, minPort, maxPort, port)
		}
		cfg.Port = port
	}

	// The value is deliberately not echoed in errors: connection strings
	// usually embed a password, and errors end up in logs.
	cfg.DatabaseURL = os.Getenv(EnvDatabaseURL)
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("%s is required", EnvDatabaseURL)
	}

	if v := os.Getenv(EnvLogLevel); v != "" {
		if !slices.Contains(validLogLevels, v) {
			return Config{}, fmt.Errorf("%s must be one of %v, got %q", EnvLogLevel, validLogLevels, v)
		}
		cfg.LogLevel = v
	}

	return cfg, nil
}
