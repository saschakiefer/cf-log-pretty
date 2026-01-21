/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package filter

import (
	"testing"

	"github.com/saschakiefer/cf-log-pretty/internal/config"
	"github.com/saschakiefer/cf-log-pretty/internal/parser"
)

// helper to build a LogMessage
func msg(level string, logger string) *parser.LogMessage {
	return &parser.LogMessage{
		Level:  level,
		Logger: logger,
	}
}

func TestFilter_Matches_LevelPriority(t *testing.T) {
	tests := []struct {
		name         string
		filterLevel  string
		messageLevel string
		expectMatch  bool
	}{
		// Standard level filtering
		{"DEBUG includes DEBUG", "DEBUG", "DEBUG", true},
		{"DEBUG includes INFO", "DEBUG", "INFO", true},
		{"INFO includes ERROR", "INFO", "ERROR", true},
		{"INFO excludes DEBUG", "INFO", "DEBUG", false},
		{"WARN includes ERROR", "WARN", "ERROR", true},
		{"ERROR includes ERROR", "ERROR", "ERROR", true},
		{"ERROR excludes INFO", "ERROR", "INFO", false},

		// Missing log level ("-----") is treated as lowest priority
		{"DEBUG includes missing level (-----)", "DEBUG", "-----", true},
		{"INFO includes missing level (-----)", "INFO", "-----", true},
		{"ERROR includes missing level (-----)", "ERROR", "-----", true},

		// Unknown message level is excluded
		{"Unknown log level gets included", "INFO", "FOO", true},

		// Unknown filter level behaves as "match nothing"
		{"Unknown filter level", "FOO", "ERROR", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New(&config.Config{Level: tt.filterLevel, Exclude: []string{""}})
			m := msg(tt.messageLevel, "com.sap.test")

			got := f.Matches(m)
			if got != tt.expectMatch {
				t.Errorf("Expected match = %v, got %v", tt.expectMatch, got)
			}
		})
	}
}

func TestFilter_Matches_LoggerFilter(t *testing.T) {
	tests := []struct {
		name          string
		filterLogger  []string
		messageLogger string
		expectMatch   bool
	}{
		// Exact matching (no wildcard)
		{"Exact match excludes logger", []string{"com.foo.auth.Service"}, "com.foo.auth.Service", false},
		{"Exact match - no match includes logger", []string{"com.foo.audit"}, "com.foo.auth.Service", true},
		{"Exact match - partial name no match", []string{"auth"}, "com.foo.auth.Service", true},

		// Package prefix matching (with wildcard)
		{"Package wildcard matches exact package", []string{"com.foo.core.*"}, "com.foo.core.Service", false},
		{"Package wildcard matches sub-package", []string{"com.foo.core.*"}, "com.foo.core.sub.Service", false},
		{"Package wildcard no match different package", []string{"com.foo.core.*"}, "com.foo.auth.Service", true},
		{"Package wildcard matches package root", []string{"com.foo.*"}, "com.foo.Service", false},
		{"Package wildcard matches nested", []string{"com.foo.*"}, "com.foo.core.sub.Service", false},
		{"Package wildcard with trailing dot", []string{"com.foo.core.*"}, "com.foo.core.Service", false},

		// Multiple filters
		{"Multiple filters - first matches", []string{"com.foo.auth.*", "com.bar.*"}, "com.foo.auth.Service", false},
		{"Multiple filters - second matches", []string{"com.foo.auth.*", "com.bar.*"}, "com.bar.Service", false},
		{"Multiple filters - no match", []string{"com.foo.auth.*", "com.bar.*"}, "com.baz.Service", true},
		{"Multiple filters - exact and wildcard", []string{"com.foo.auth.Service", "com.bar.*"}, "com.foo.auth.Service", false},
		{"Multiple filters - exact and wildcard 2", []string{"com.foo.auth.Service", "com.bar.*"}, "com.bar.test.Service", false},

		// Empty and edge cases
		{"Empty logger filter matches everything", []string{}, "anything.at.all", true},
		{"Empty string filter", []string{""}, "com.foo.Service", true},
		{"Wildcard only", []string{"*"}, "com.foo.Service", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New(&config.Config{Level: "DEBUG", Exclude: tt.filterLogger})
			m := msg("INFO", tt.messageLogger)

			got := f.Matches(m)
			if got != tt.expectMatch {
				t.Errorf("Expected match = %v, got %v", tt.expectMatch, got)
			}
		})
	}
}
