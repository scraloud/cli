package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"time"
)

func init() {
	rootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy current scraper",
	Long:  "Deploys current scraper to server and follows logs",
	Run: func(cmd *cobra.Command, args []string) {
		_ = GetTokenOrFail()

		if out, err := exec.Command("git", "add", ".").CombinedOutput(); err != nil {
			log.Fatal(string(out), err)
		}

		if out, err := exec.Command("git", "commit", "-m", fmt.Sprintf(`"%s"`, time.Now().String())).CombinedOutput(); err != nil {
			log.Fatal(string(out), err)
		}

		fmt.Println("Pushing Code...")
		if out, err := exec.Command("git", "push", "scraloud", "master").CombinedOutput(); err != nil {
			log.Fatal(string(out), err)
		}

		fmt.Println("Running Scraper...")
		logsCmd.Run(cmd, args)
	},
}
