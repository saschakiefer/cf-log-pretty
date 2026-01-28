/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package cmd

import (
	"testing"

	"github.com/saschakiefer/cf-log-pretty/internal/config"
)

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid log level TRACE",
			config: &config.Config{
				Level: "TRACE",
			},
			expectError: false,
		},
		{
			name: "valid log level DEBUG",
			config: &config.Config{
				Level: "DEBUG",
			},
			expectError: false,
		},
		{
			name: "valid log level INFO",
			config: &config.Config{
				Level: "INFO",
			},
			expectError: false,
		},
		{
			name: "valid log level WARN",
			config: &config.Config{
				Level: "WARN",
			},
			expectError: false,
		},
		{
			name: "valid log level ERROR",
			config: &config.Config{
				Level: "ERROR",
			},
			expectError: false,
		},
		{
			name: "valid log level lowercase",
			config: &config.Config{
				Level: "info",
			},
			expectError: false,
		},
		{
			name: "valid log level mixed case",
			config: &config.Config{
				Level: "WaRn",
			},
			expectError: false,
		},
		{
			name: "invalid log level",
			config: &config.Config{
				Level: "INVALID",
			},
			expectError: true,
			errorMsg:    "invalid log level: INVALID (allowed: TRACE, DEBUG, INFO, WARN, ERROR)",
		},
		{
			name: "empty log level (allowed)",
			config: &config.Config{
				Level: "",
			},
			expectError: false,
		},
		{
			name: "valid with remove prefix",
			config: &config.Config{
				Level:        "INFO",
				RemovePrefix: "com.foo.prod.",
			},
			expectError: false,
		},
		{
			name: "valid with logger name only",
			config: &config.Config{
				Level:          "INFO",
				LoggerNameOnly: true,
			},
			expectError: false,
		},
		{
			name: "invalid: both logger name only and remove prefix",
			config: &config.Config{
				Level:          "INFO",
				LoggerNameOnly: true,
				RemovePrefix:   "com.foo.",
			},
			expectError: true,
			errorMsg:    "cannot use --show-logger-name-only and --remove-logger-prefix together",
		},
		{
			name: "valid with exclude logger",
			config: &config.Config{
				Level:   "INFO",
				Exclude: []string{"com.foo.Service", "com.bar.*"},
			},
			expectError: false,
		},
		{
			name: "valid with truncate raw",
			config: &config.Config{
				Level:       "INFO",
				TruncateRaw: true,
			},
			expectError: false,
		},
		{
			name: "valid with all compatible flags",
			config: &config.Config{
				Level:        "DEBUG",
				RemovePrefix: "com.example.",
				Exclude:      []string{"com.example.noise.*"},
				TruncateRaw:  true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original config and restore after test
			origCfg := cfg
			defer func() { cfg = origCfg }()

			// Set the global config to the test config
			cfg = tt.config

			// Call validateFlags
			err := validateFlags(nil, nil)

			// Check if error matches expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q but got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateFlagsLogLevelCaseInsensitive(t *testing.T) {
	levels := []string{"trace", "TRACE", "TrAcE", "debug", "DEBUG", "info", "INFO", "warn", "WARN", "error", "ERROR"}

	for _, level := range levels {
		t.Run("level_"+level, func(t *testing.T) {
			// Store original config and restore after test
			origCfg := cfg
			defer func() { cfg = origCfg }()

			cfg = &config.Config{
				Level: level,
			}

			err := validateFlags(nil, nil)
			if err != nil {
				t.Errorf("expected no error for level %q but got: %v", level, err)
			}
		})
	}
}
