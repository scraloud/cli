package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func init() {
	rootCmd.AddCommand(logsCmd)
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print logs",
	Long:  "Print recent logs",
	Run: func(cmd *cobra.Command, args []string) {
		token := CheckLogin(cmd, args)

		var old []byte
		for {
			time.Sleep(time.Second)
			resp, err := http.Get(apiURL + "/scrapers/logs/?token=" + token)
			if err != nil {
				log.Fatal(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Fatal(string(body))
			}

			if bytes.Equal(old, body) {
				continue
			}
			if len(old) == 0 {
				old = body
				fmt.Print(string(body))
			} else {
				diff := bytes.TrimPrefix(body, old)
				old = body
				fmt.Print(string(diff))
			}
			if bytes.HasSuffix(body, []byte("Scraper Running Finished\n")) {
				break
			}
		}

	},
}

func GetProjectName() (string, error) {
	output, err := exec.Command("git", "remote", "get-url", "scraloud").CombinedOutput()
	if err != nil {
		return "", err
	}

	cleanOutput := strings.TrimSpace(string(output))
	parts := strings.Split(cleanOutput, "/")
	lastPart := parts[len(parts)-1]
	projectName := strings.TrimSuffix(lastPart, ".git")

	return projectName, nil
}
