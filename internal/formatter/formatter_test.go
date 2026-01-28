/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package formatter

import (
	"strings"
	"testing"

	"github.com/saschakiefer/cf-log-pretty/internal/config"
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
	output := Format(msg, LevelColorizer(msg.Level), &config.Config{TruncateRaw: false})

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

func TestFormatRawTruncation_NoColor(t *testing.T) {
	msg := &parser.LogMessage{
		Timestamp:     "2024-01-01T12:00:00.00",
		Level:         "ERROR",
		Logger:        "com.example.service.MyLogger",
		Message:       "Eum nesciunt consequatur.",
		StackTrace:    []string{},
		HasParseError: true,
	}

	//output := Format(msg, NoColor)
	output := Format(msg, LevelColorizer(msg.Level), &config.Config{TruncateRaw: true})

	if len(output) != 80 {
		t.Errorf("Expected output length of 80 characters, got: %d: %s", len(output), output)
	}

	if !strings.HasSuffix(output, "...") {
		t.Errorf("Expected output to end with '...', got: %s", output)
	}
}

func TestFormatRawNoTruncation_NoColor(t *testing.T) {
	msg := &parser.LogMessage{
		Timestamp:     "2024-01-01T12:00:00.00",
		Level:         "ERROR",
		Logger:        "com.example.service.MyLogger",
		Message:       "Eum ne",
		StackTrace:    []string{},
		HasParseError: true,
	}

	//output := Format(msg, NoColor)
	output := Format(msg, LevelColorizer(msg.Level), &config.Config{TruncateRaw: true})

	if len(output) != 80 {
		t.Errorf("Expected output length of 80 characters, got: %d: %s", len(output), output)
	}

	if !strings.HasSuffix(output, "Eum ne") {
		t.Errorf("Expected output to end with 'Eum ne' (no truncation), got: %s", output)
	}
}

func TestFormat_LoggerNameOnly(t *testing.T) {
	msg := &parser.LogMessage{
		Timestamp:  "2024-01-01T12:00:00.00",
		Level:      "INFO",
		Logger:     "com.example.service.MyLogger",
		Message:    "Test message",
		StackTrace: []string{},
	}

	output := Format(msg, LevelColorizer(msg.Level), &config.Config{
		TruncateRaw:    false,
		LoggerNameOnly: true,
	})

	// Should contain only "MyLogger" not the full package path
	if !strings.Contains(output, "MyLogger") {
		t.Errorf("Expected 'MyLogger' in output, got: %s", output)
	}

	// Should NOT contain the package prefix
	if strings.Contains(output, "com.example.service.MyLogger") {
		t.Errorf("Expected package prefix to be removed, got: %s", output)
	}

	// Should still contain the message
	if !strings.Contains(output, "Test message") {
		t.Errorf("Expected message in output, got: %s", output)
	}
}

func TestFormat_LoggerNameOnlyWithSingleName(t *testing.T) {
	msg := &parser.LogMessage{
		Timestamp:  "2024-01-01T12:00:00.00",
		Level:      "DEBUG",
		Logger:     "SimpleLogger",
		Message:    "Debug message",
		StackTrace: []string{},
	}

	output := Format(msg, LevelColorizer(msg.Level), &config.Config{
		TruncateRaw:    false,
		LoggerNameOnly: true,
	})

	// Should still work with a logger name that has no package
	if !strings.Contains(output, "SimpleLogger") {
		t.Errorf("Expected 'SimpleLogger' in output, got: %s", output)
	}
}

func TestFormat_LoggerNameOnlyDisabled(t *testing.T) {
	msg := &parser.LogMessage{
		Timestamp:  "2024-01-01T12:00:00.00",
		Level:      "WARN",
		Logger:     "com.example.service.MyLogger",
		Message:    "Warning message",
		StackTrace: []string{},
	}

	output := Format(msg, LevelColorizer(msg.Level), &config.Config{
		TruncateRaw:    false,
		LoggerNameOnly: false,
	})

	// Should contain the full logger name when LoggerNameOnly is false
	if !strings.Contains(output, "com.example.service.MyLogger") {
		t.Errorf("Expected full logger name in output, got: %s", output)
	}
}
