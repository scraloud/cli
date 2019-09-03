package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jdxcode/netrc"
	"github.com/spf13/cobra"
	"io/ioutil"
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
			fmt.Println(err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		// Read Body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		token := struct {
			Token string
		}{}
		if err := json.Unmarshal(body, &token); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if token.Token == "" {
			fmt.Println(string(body))
			os.Exit(1)
		}

		SaveLogin(email, token.Token)

		fmt.Println("Login Successful")
	},
}

func GetTokenOrFail() string {
	n, _ := ReadNetrc()

	parsedApiURL, _ := url.Parse(apiURL)

	if n.Machine(parsedApiURL.Host) == nil {
		fmt.Println("Please Login")
		os.Exit(1)
	}

	token := n.Machine(parsedApiURL.Host).Get("password")
	if token == "" {
		fmt.Println("Please Login")
		os.Exit(1)
	}

	return token
}

func SaveLogin(email string, password string) {
	n, _ := ReadNetrc()

	parsedApiURL, _ := url.Parse(apiURL)

	n.AddMachine(parsedApiURL.Host, email, password)

	// Save .netrc file
	if err := n.Save(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
