/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the health of specified URL(s)",
	Long: `Performs a health check by sending a request to the specified URL(s)
	and report the status`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("check called")
		for _, url := range args {
			checkURL(url, threshold, retries)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func checkURL(url string, threshold float64, retries int) {
	var resp *http.Response
	var err error
	var duration time.Duration

	for attempt := 0; attempt <= retries; attempt++ {
		start := time.Now()
		resp, err = http.Get(url)
		duration = time.Since(start)

		if err == nil {
			break
		}

		if attempt < retries {
			// fmt.Fprintf(os.Stderr, "Attempt %d failed, retrying...\n", attempt+1)
			l.Warn("faile retrying", "url", url, "attempts", attempt+1)
			time.Sleep(time.Second * 2) // Backoff strategy
		}

	}

	if err != nil {
		// fmt.Fprintf(os.Stderr, "Error fetching %s after %d retries: %v\n", url, retries, err)
		l.Error("fetching error", "url", url, "retries", retries, "err", err)
		return
	}

	defer resp.Body.Close()

	if duration.Seconds() > threshold && verbose {
		l.Error("exceeded threshold", "url", url, "duration", duration, "threshold", threshold)
		// fmt.Fprintf(os.Stderr, "Warning: %s response time (%v) exceeded threshold of %fs\n", url, duration, threshold)
	}
	l.Info("successful check", "url", url, "status code", resp.StatusCode, "duration", duration)
	// logger.Printf("Checked %s, Status: %d, Response Time: %v", url, resp.StatusCode, duration)
}
