package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "scraloud",
	Short: "Scraloud is scraping cloud CLI for scraloud.com",
	Long:  `Scraloud deploys your scraping projects to cloud, manages requests, exported items and provides metrics`,
}

var apiURL = "https://api.scraloud.com/v1"
var gitURL = "https://git.scraloud.com/"

func init() {
	// Set Base URLs
	if envURL := os.Getenv("SCRALOUD_API_URL"); envURL != "" {
		apiURL = envURL
	}
	if envURL := os.Getenv("SCRALOUD_GIT_URL"); envURL != "" {
		gitURL = envURL
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
