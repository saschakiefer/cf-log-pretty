/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/saschakiefer/cf-log-pretty/internal/formatter"
	"github.com/saschakiefer/cf-log-pretty/internal/parser"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cf-log-pretty",
	Short: "Convert SAP BTP Cloud Foundry logs to human readable format",
	Long:  `Convert SAP BTP Cloud Foundry logs to human readable format`,
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			line := scanner.Text()

			msg, ok := parser.ParseLine(line)
			if !ok {
				continue // skip malformed lines
			}

			fmt.Println(formatter.Format(msg, formatter.LevelColorizer(msg.Level)))
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("level", "l", "DEBUG", "Help message for toggle")
}
