/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package formatter

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/saschakiefer/cf-log-pretty/internal/parser"
)

// ColorFunc defines a flexible formatting function (with or without ANSI colors)
type ColorFunc func(format string, a ...interface{}) string

// Format renders a log message using the provided level color function
func Format(msg *parser.LogMessage, colorizeLevel ColorFunc) string {
	levelText := colorizeLevel("[%-5s]", msg.Level)

	result := fmt.Sprintf("%-22s %s %-40s : %s",
		msg.Timestamp,
		levelText,
		shortenMiddle(msg.Logger, 40),
		msg.Message,
	)

	if len(msg.StackTrace) > 0 {
		for _, line := range msg.StackTrace {
			result += "\n    " + line
		}
	}

	return result
}

func shortenMiddle(input string, max int) string {
	if len(input) <= max {
		// Pad with spaces if shorter than max
		return fmt.Sprintf("%-*s", max, input)
	}

	// Length of prefix and suffix we want to keep
	// (e.g. 18 + 3 + 19 = 40)
	dots := "..."
	keep := (max - len(dots)) / 2
	start := input[:keep]
	end := input[len(input)-(max-len(dots)-keep):]

	return start + dots + end
}

// NoColor is used for test output or non-terminal pipes
func NoColor() ColorFunc {
	return func(format string, a ...interface{}) string {
		return fmt.Sprintf(format, a...)
	}
}

// Alternative - Not sure yet, which is better
// equivalent to `var NoColor ColorFunc = func(format string, a ...interface{}) string {`
/*
var NoColor = func(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}
*/

// LevelColorizer returns a color formatting function for the given log level
func LevelColorizer(level string) ColorFunc {
	switch level {
	case "ERROR":
		return color.New(color.FgRed).Add(color.Bold).SprintfFunc()
	case "WARN":
		return color.New(color.FgYellow).Add(color.Bold).SprintfFunc()
	case "INFO":
		return color.New(color.FgCyan).Add(color.Bold).SprintfFunc()
	case "DEBUG":
		return color.New(color.FgHiBlack).Add(color.Bold).SprintfFunc()
	default:
		return NoColor()
	}
}
