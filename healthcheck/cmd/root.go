/*
Copyright © 2026 Lucio
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/luciormoraes/healthcheck-with-cobra/healthcheck/logger"
	"github.com/spf13/cobra"
)

var (
	logFile   string
	l         *slog.Logger
	threshold float64
	retries   int

	silent  bool
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "A tool for monitoring health status and responsiveness of web applications",
	Long: `
	The healcheck command is designed to assess the health and
	responsiveness of specified web applications. It sends HTTP requests
	to URLs provided by the user, evaluating whether the services are
	accessible and how quick they respond. This command supports both
	immediate, one-off checking and continuous monitoring, allowing users
	to specify intervals for ongoing health assessment. With additional
	flags for customization, users can tailor the command to meet various
	monitoring needs, from simple uptime checks to detailed performance
	analysis.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		l = logger.New(logFile, silent, verbose)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logFile, "logfile", "healthcheck.log", "File to log output to")
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().Float64Var(&threshold, "threshold", 0.5, "Threshold value for considering a response to be too slow (in seconds)")
	rootCmd.PersistentFlags().IntVar(&retries, "retries", 3, "Number of retries for a failed request")
	rootCmd.PersistentFlags().BoolVar(&silent, "silent", false, "Run in silent mode without stdout output")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Run in verbose mode. Override silent mode")

}
