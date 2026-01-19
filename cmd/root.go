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

	"github.com/saschakiefer/cf-log-pretty/internal/filter"
	"github.com/saschakiefer/cf-log-pretty/internal/formatter"
	"github.com/saschakiefer/cf-log-pretty/internal/parser"
	"github.com/spf13/cobra"
)

var levelFlag string
var excludeLogger []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cf-log-pretty",
	Short:   "Convert SAP BTP Cloud Foundry logs to human readable format",
	Long:    `Convert SAP BTP Cloud Foundry logs to human readable format`,
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
	rootCmd.Flags().StringVarP(&levelFlag, "level", "l", "DEBUG", "Minimum log level to include (TRACE, DEBUG, INFO, WARN, ERROR).")
	rootCmd.Flags().StringSliceVarP(&excludeLogger, "exclude-logger", "e", []string{}, "Exclude logs from given loggers (example: -e \"com.foo.l1,com.foo.l2\"")
}

func validateFlags(cmd *cobra.Command, args []string) error {
	level := strings.ToUpper(levelFlag)

	allowed := map[string]bool{
		"TRACE": true,
		"DEBUG": true,
		"INFO":  true,
		"WARN":  true,
		"ERROR": true,
	}

	if level != "" && !allowed[level] {
		return fmt.Errorf("invalid log level: %s (allowed: TRACE, DEBUG, INFO, WARN, ERROR)", levelFlag)
	}
	return nil
}

func run(cmd *cobra.Command, args []string) {
	f := filter.New(levelFlag, excludeLogger)

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

		fmt.Println(formatter.Format(msg, formatter.LevelColorizer(msg.Level)))
	}
}
