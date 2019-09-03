package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new project for scraloud and adds git remote.",
	Long:  `Create a new project and add git remote for deployment.`,
	Run: func(cmd *cobra.Command, args []string) {
		token := CheckLogin(cmd, args)

		resp, err := http.PostForm(apiURL+"/scrapers/?token="+token, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// Read Body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		name := struct {
			Name string
		}{}
		if err := json.Unmarshal(body, &name); err != nil {
			log.Fatal(err)
		}
		if name.Name == "" {
			log.Fatal(string(body))
		}

		if output, err := exec.Command("git", "init").CombinedOutput(); err != nil {
			log.Fatal(string(output))
		}

		if out, err := exec.Command("git", "remote", "add", "scraloud", gitURL+name.Name+".git").CombinedOutput(); err != nil {
			log.Fatal(string(out), err)
		}

		//if out, err := exec.Command("git", "branch", "--set-upstream-to", "scraloud/master").CombinedOutput(); err != nil {
		//	log.Fatal(string(out), err)
		//}

		fmt.Println("Scraper Created: ", name.Name)
	},
}
