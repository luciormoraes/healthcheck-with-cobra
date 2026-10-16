/*
Copyright © 2026 Lucio
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	interval time.Duration
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor [urls]",
	Short: "Monitor the health of specified URL(s) over time",
	Long: `
	Continuously monitors the health of the specified URL(s)
	at the specified interval`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		monitorURLs(ctx, args)
	},
}

func init() {
	monitorCmd.Flags().DurationVar(&interval, "interval", 2*time.Second, "Interval between healthchecks")
	rootCmd.AddCommand(monitorCmd)
}

func monitorURLs(ctx context.Context, urls []string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if output == "table" {
			// print table
			table := tablewriter.NewWriter(outputWriter)
			upColor := color.New(color.FgGreen).SprintFunc()
			downColor := color.New(color.FgRed).SprintFunc()
			table.Header([]string{"URL", "Status"})
			for _, url := range urls {
				status := downColor("Down")
				if checkURL(ctx, url, threshold, retries) {
					status = upColor("Up")
				}
				table.Append([]string{url, status})
			}
			// clear the screen
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				fmt.Println("Unable to clean the screen", err)
			}
			table.Render()
		} else {
			for _, url := range urls {
				checkURL(ctx, url, threshold, retries)
			}
		}

	}
}
