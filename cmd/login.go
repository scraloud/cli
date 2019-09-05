package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jdxcode/netrc"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login",
	Long:  `Login`,
	Run: func(cmd *cobra.Command, args []string) {

		// Get Credentials from user
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		fmt.Print("Password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		fmt.Println("Logging in...")

		// Login
		resp, err := http.PostForm(apiURL+"/users/login/", url.Values{
			"email":    {email},
			"password": {password},
			"cli":      {"true"},
		})
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// Read Body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Check Status
		if resp.StatusCode != http.StatusOK {
			log.Fatal(string(body))
		}

		// Extract Token
		token := struct {
			Token string
		}{}
		if err := json.Unmarshal(body, &token); err != nil {
			log.Fatal(string(body), err)
		}
		if token.Token == "" {
			log.Fatal(string(body))
		}

		if err := SaveLogin(email, token.Token); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Login Successful")
	},
}

func CheckLogin(cmd *cobra.Command, args []string) string {
	n, _ := ReadNetrc()

	parsedApiURL, _ := url.Parse(apiURL)

	if n.Machine(parsedApiURL.Host) == nil {
		fmt.Println("Please Login First")
		loginCmd.Run(cmd, args)
		return CheckLogin(cmd, args)
	}

	return n.Machine(parsedApiURL.Host).Get("password")
}

func SaveLogin(email string, password string) error {
	n, _ := ReadNetrc()

	parsedApiURL, err := url.Parse(apiURL)
	if err != nil {
		return err
	}
	parsedGitURL, err := url.Parse(gitURL)
	if err != nil {
		return err
	}

	n.AddMachine(parsedApiURL.Host, email, password)
	n.AddMachine(parsedGitURL.Host, email, password)

	// Save .netrc file
	if err := n.Save(); err != nil {
		return err
	}

	return nil
}

func ReadNetrc() (*netrc.Netrc, error) {
	// Get Current User
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read .netrc file
	n, err := netrc.Parse(filepath.Join(usr.HomeDir, ".netrc"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return n, nil
}
