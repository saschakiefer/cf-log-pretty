/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package parser

import (
	"encoding/json"
	"regexp"
	"strings"
)

type LogMessage struct {
	Timestamp  string
	Source     string
	Direction  string
	Level      string
	Logger     string
	Message    string
	StackTrace []string
	Raw        string
}

var cfPrefixRegex = regexp.MustCompile(`^\s*(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{2,})\+\d{4}\s+\[([^]]+)]\s+(OUT|ERR)\s+`)

type structuredLog struct {
	WrittenAt  string   `json:"written_at"`
	Level      string   `json:"level"`
	Logger     string   `json:"logger"`
	Msg        string   `json:"msg"`
	StackTrace []string `json:"stacktrace"`
}

// ParseLine parses one line of CF log
func ParseLine(line string) (*LogMessage, bool) {
	matches := cfPrefixRegex.FindStringSubmatch(line)
	loc := cfPrefixRegex.FindStringIndex(line)

	if matches == nil || loc == nil || len(matches) < 4 {
		return parseFallbackLine(line)
	}

	timestamp := matches[1]
	source := matches[2]
	direction := matches[3]
	rest := strings.TrimSpace(line[loc[1]:])

	msg := &LogMessage{
		Timestamp: timestamp,
		Source:    source,
		Level:     "-----", // Default will be overwritten if available
		Direction: direction,
		Raw:       line,
	}

	// Try to parse the remaining content as JSON
	rest = strings.ReplaceAll(rest, "\t", "\\t")
	rest = strings.ReplaceAll(rest, "\n", "\\n")

	var structured structuredLog
	decoder := json.NewDecoder(strings.NewReader(rest))

	err := decoder.Decode(&structured)
	if err != nil {
		// Not a valid JSON log, fallback to plain message
		msg.Message = rest
		return msg, true
	}

	msg.Level = structured.Level

	msg.Logger = structured.Logger
	msg.Message = structured.Msg
	msg.StackTrace = structured.StackTrace

	return msg, true
}

// parseFallbackLine handles lines that don't match expected format
func parseFallbackLine(line string) (*LogMessage, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return nil, false
	}
	return &LogMessage{
		Message: trimmed,
		Level:   "-----",
		Raw:     line,
	}, true
}
