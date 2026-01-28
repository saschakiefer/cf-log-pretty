/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/saschakiefer/cf-log-pretty/internal/config"
	"github.com/saschakiefer/cf-log-pretty/internal/filter"
	"github.com/saschakiefer/cf-log-pretty/internal/formatter"
	"github.com/saschakiefer/cf-log-pretty/internal/parser"
	"github.com/spf13/cobra"
)

var (
	Version = "1.1.0"
	cfg     = &config.Config{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cf-log-pretty",
	Version: Version,
	Short:   "Convert SAP BTP Cloud Foundry logs to human readable format",
	Long: `cf-log-pretty is a command-line tool designed to format and colorize log output 
from SAP BTP Cloud Foundry. It parses the standard Cloud Foundry log format, 
including structured JSON logs, making them easier to read in a terminal.

It reads from standard input (stdin), allowing you to pipe the output 
of 'cf logs' directly into it:

    cf logs <app-name> | cf-log-pretty`,
	PreRunE: validateFlags,
	Run:     run,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&cfg.Level, "level", "l", "TRACE", "minimum log level to include (TRACE, DEBUG, INFO, WARN, ERROR)")
	rootCmd.Flags().StringVarP(&cfg.RemovePrefix, "remove-logger-prefix", "r", "", "remove given prefix from logger names (e.g. \"com.foo.prod.\")")
	rootCmd.Flags().StringSliceVarP(&cfg.Exclude, "exclude-logger", "e", []string{}, "exclude logs from given loggers. Supports exact match (e.g. \"com.foo.Service\") or package wildcard (e.g. \"com.foo.core.*\" for packages and sub-packages)")
	rootCmd.Flags().BoolVarP(&cfg.TruncateRaw, "truncate-raw", "t", false, "truncate raw log messages to terminal width (if message is not in JSON format, e.g. platform logs)")

}

func validateFlags(_ *cobra.Command, _ []string) error {
	level := strings.ToUpper(cfg.Level)

	allowed := map[string]bool{
		"TRACE": true,
		"DEBUG": true,
		"INFO":  true,
		"WARN":  true,
		"ERROR": true,
	}

	if level != "" && !allowed[level] {
		return fmt.Errorf("invalid log level: %s (allowed: TRACE, DEBUG, INFO, WARN, ERROR)", cfg.Level)
	}
	return nil
}

func run(_ *cobra.Command, _ []string) {
	f := filter.New(cfg)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		msg, ok := parser.ParseLine(line)
		if !ok {
			continue // skip malformed lines
		}

		if !f.Matches(msg) {
			continue
		}

		fmt.Println(formatter.Format(msg, formatter.LevelColorizer(msg.Level), cfg))
	}
}
