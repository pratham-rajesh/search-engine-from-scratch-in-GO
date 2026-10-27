package config

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	const dsn = "postgres://search:search@localhost:5432/search"

	tests := []struct {
		name    string
		env     map[string]string
		want    Config
		wantErr string // substring the error must contain; empty means no error
	}{
		{
			name: "defaults applied when only the database url is set",
			env:  map[string]string{EnvDatabaseURL: dsn},
			want: Config{Port: 8080, DatabaseURL: dsn, LogLevel: "info"},
		},
		{
			name: "all values set",
			env: map[string]string{
				EnvPort:        "9090",
				EnvDatabaseURL: dsn,
				EnvLogLevel:    "debug",
			},
			want: Config{Port: 9090, DatabaseURL: dsn, LogLevel: "debug"},
		},
		{
			name: "lowest valid port",
			env:  map[string]string{EnvPort: "1", EnvDatabaseURL: dsn},
			want: Config{Port: 1, DatabaseURL: dsn, LogLevel: "info"},
		},
		{
			name: "highest valid port",
			env:  map[string]string{EnvPort: "65535", EnvDatabaseURL: dsn},
			want: Config{Port: 65535, DatabaseURL: dsn, LogLevel: "info"},
		},
		{
			name:    "missing database url",
			env:     map[string]string{},
			wantErr: EnvDatabaseURL,
		},
		{
			name:    "port is not a number",
			env:     map[string]string{EnvPort: "http", EnvDatabaseURL: dsn},
			wantErr: EnvPort,
		},
		{
			name:    "port zero",
			env:     map[string]string{EnvPort: "0", EnvDatabaseURL: dsn},
			wantErr: EnvPort,
		},
		{
			name:    "port negative",
			env:     map[string]string{EnvPort: "-1", EnvDatabaseURL: dsn},
			wantErr: EnvPort,
		},
		{
			name:    "port too large",
			env:     map[string]string{EnvPort: "65536", EnvDatabaseURL: dsn},
			wantErr: EnvPort,
		},
		{
			name:    "unknown log level",
			env:     map[string]string{EnvDatabaseURL: dsn, EnvLogLevel: "verbose"},
			wantErr: EnvLogLevel,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Start from a clean slate so variables set in the developer's
			// shell can't leak into the test. Load treats "" as unset.
			for _, key := range []string{EnvPort, EnvDatabaseURL, EnvLogLevel} {
				t.Setenv(key, "")
			}
			for key, value := range tc.env {
				t.Setenv(key, value)
			}

			got, err := Load()

			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("Load() = %+v, want error containing %q", got, tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("Load() error = %q, want it to mention %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load() unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("Load() = %+v, want %+v", got, tc.want)
			}
		})
	}
}
