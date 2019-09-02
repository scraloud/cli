package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop current scraper",
	Long:  "Stop current running scraper",
	Run: func(cmd *cobra.Command, args []string) {
		token := GetTokenOrFail()

		fmt.Println("Stopping Scraper...")

		resp, err := http.PostForm(baseURL+"/scrapers/commands/stops/?token="+token, nil)
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

		fmt.Println("Scraper Successfully Stopped")
	},
}
