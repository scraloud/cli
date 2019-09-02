package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "scraloud",
	Short: "Scraloud is scraping cloud CLI for scraloud.com",
	Long:  `Scraloud deploys your scraping projects to cloud, manages requests, exported items and provides metrics`,
}

var baseURL = "http://api.scraloud.loc/v1"
var baseURLHost string

func init() {
	// Set Base URL
	if envURL := os.Getenv("SCRALOUD_BASE_URL"); envURL != "" {
		baseURL = envURL
	}
	parsedBaseURL, _ := url.Parse(baseURL)
	baseURLHost = parsedBaseURL.Host
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
