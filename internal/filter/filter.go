/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package filter

import (
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
	Level  string
	Logger string
}

func New(level, logger string) *Filter {
	return &Filter{
		Level:  strings.ToUpper(level),
		Logger: logger,
	}
}

func (f *Filter) Matches(msg *parser.LogMessage) bool {
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

	if f.Logger != "" && !strings.Contains(msg.Logger, f.Logger) {
		return false
	}

	return true
}
