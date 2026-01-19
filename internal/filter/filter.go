/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package filter

import (
	"slices"
	"strings"

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
	Level         string
	ExcludeLogger []string
}

func New(level string, exclude []string) *Filter {
	return &Filter{
		Level:         strings.ToUpper(level),
		ExcludeLogger: exclude,
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
	if len(f.ExcludeLogger) > 0 && slices.Contains(f.ExcludeLogger, msg.Logger) {
		return false
	}

	return true
}
