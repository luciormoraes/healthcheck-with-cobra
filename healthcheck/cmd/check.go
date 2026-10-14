/*
Copyright © 2026 Lucio
*/
package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the health of specified URL(s)",
	Long: `
	Performs a health check by sending a request to the specified URL(s)
	and report the status`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, url := range args {
			checkURL(url, threshold, retries)
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		for _, url := range args {
			if !isValidURL(url) {
				return fmt.Errorf("invalid url: %s", url)
			}
		}
		return nil
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
			l.Warn("faile retrying", "url", url, "attempts", attempt+1)
			time.Sleep(time.Second * 2) // Backoff strategy
		}

	}

	if err != nil {
		l.Error("fetching error", "url", url, "retries", retries, "err", err)
		return
	}

	defer resp.Body.Close()

	if duration.Seconds() > threshold && verbose {
		l.Error("exceeded threshold", "url", url, "duration", duration, "threshold", threshold)
	}
	l.Info("successful check", "url", url, "status code", resp.StatusCode, "duration", duration)
}

func isValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}
