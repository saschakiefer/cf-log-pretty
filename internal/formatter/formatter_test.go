/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package formatter

import (
	"strings"
	"testing"

	"github.com/saschakiefer/cf-log-pretty/internal/parser"
)

func TestFormat_NoColor(t *testing.T) {
	msg := &parser.LogMessage{
		Timestamp: "2024-01-01T12:00:00.00",
		Level:     "ERROR",
		Logger:    "com.example.service.MyLogger",
		Message:   "Something bad happened",
		StackTrace: []string{
			"java.lang.Exception: Something bad happened",
			"\tat com.example.Class.method(Class.java:42)",
		},
	}

	//output := Format(msg, NoColor)
	output := Format(msg, LevelColorizer(msg.Level))

	// Basic assertions
	if !strings.Contains(output, "[ERROR]") {
		t.Errorf("Expected [ERROR] in output, got: %s", output)
	}

	if !strings.Contains(output, "com.example.service.MyLogger") {
		t.Errorf("Expected logger name in output, got: %s", output)
	}

	if !strings.Contains(output, "Something bad happened") {
		t.Errorf("Expected message in output, got: %s", output)
	}

	if !strings.Contains(output, "java.lang.Exception") {
		t.Errorf("Expected stacktrace line in output, got: %s", output)
	}

	// Optional: check line count
	lines := strings.Split(output, "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines (main + 2 stacktrace), got %d", len(lines))
	}
}
