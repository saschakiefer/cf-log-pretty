/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package filter

import (
	"testing"

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
			f := New(tt.filterLevel, []string{""})
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
		{"Match exact logger", []string{"auth"}, "auth", false},
		{"Match substring", []string{"auth"}, "com.foo.auth.Service", true},
		{"No match", []string{"audit"}, "com.foo.auth.Service", true},
		{"Empty logger filter matches everything", []string{}, "anything.at.all", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New("DEBUG", tt.filterLogger)
			m := msg("INFO", tt.messageLogger)

			got := f.Matches(m)
			if got != tt.expectMatch {
				t.Errorf("Expected match = %v, got %v", tt.expectMatch, got)
			}
		})
	}
}
