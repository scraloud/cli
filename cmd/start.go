package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start scraper",
	Long:  "Start existing scraper",
	Run: func(cmd *cobra.Command, args []string) {
		token := GetTokenOrFail()

		fmt.Println("Starting Scraper...")

		resp, err := http.PostForm(baseURL+"/scrapers/commands/starts/?token="+token, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(string(body), err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatal("Failed: ", string(body))
		}

		logsCmd.Run(cmd, args)
	},
}
