/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package filter

import (
	"strings"

	"github.com/saschakiefer/cf-log-pretty/internal/config"
	"github.com/saschakiefer/cf-log-pretty/internal/parser"
)

var LevelPriority = map[string]int{
	"TRACE": 1,
	"DEBUG": 2,
	"INFO":  3,
	"WARN":  4,
	"ERROR": 5,
	"-----": 6,
}

type Filter struct {
	Level   string
	Exclude []string
}

func New(cfg *config.Config) *Filter {
	return &Filter{
		Level:   strings.ToUpper(cfg.Level),
		Exclude: cfg.Exclude,
	}
}

func (f *Filter) Matches(msg *parser.LogMessage) bool {
	// Log Level
	msgLevel := strings.ToUpper(msg.Level)

	logPrio, okLog := LevelPriority[msgLevel]
	filterPrio, okFilter := LevelPriority[f.Level]

	if !okLog {
		// Unknown log level → treat as lowest (fallback)
		logPrio = LevelPriority["-----"]
	}
	if !okFilter {
		// Unknown filter level → match nothing
		return false
	}

	if logPrio < filterPrio {
		return false
	}

	// Exclude Logger
	if len(f.Exclude) > 0 && f.matchesExcludedLogger(msg.Logger) {
		return false
	}

	return true
}

// matchesExcludedLogger checks if the logger name matches any exclude pattern.
// Patterns ending with "*" are treated as package prefixes (e.g., "com.foo.core.*" matches "com.foo.core.Service").
// Patterns without "*" must match exactly.
func (f *Filter) matchesExcludedLogger(logger string) bool {
	for _, pattern := range f.Exclude {
		if strings.HasSuffix(pattern, "*") {
			// Package prefix matching
			prefix := strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(logger, prefix) {
				return true
			}
		} else {
			// Exact match
			if logger == pattern {
				return true
			}
		}
	}
	return false
}
