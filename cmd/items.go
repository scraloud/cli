package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var exportFileName string

func init() {
	rootCmd.AddCommand(itemsCmd)
	itemsCmd.Flags().StringVarP(&exportFileName, "export", "e", "", "Export items to file")
}

var itemsCmd = &cobra.Command{
	Use:   "items",
	Short: "Print items",
	Long:  "Print recent items collected",
	Run: func(cmd *cobra.Command, args []string) {
		token := CheckLogin(cmd, args)

		resp, err := http.Get(apiURL + "/scrapers/items/?token=" + token)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if exportFileName != "" {
			if err := ExportToFile(resp); err != nil {
				log.Fatal(err)
			}
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatal(string(body))
		}

		// Colored Output
		var items []map[string]interface{}
		if err := json.Unmarshal(body, &items); err != nil {
			log.Fatal(err)
		}
		f := &prettyjson.Formatter{
			KeyColor:        color.New(color.FgHiWhite),
			StringColor:     color.New(color.FgHiYellow),
			BoolColor:       color.New(color.FgHiGreen),
			NumberColor:     color.New(color.FgHiCyan),
			NullColor:       color.New(color.FgHiRed),
			StringMaxLength: 0,
			DisabledColor:   false,
			Indent:          2,
			Newline:         "\n",
		}
		s, _ := f.Marshal(items)
		fmt.Println(string(s))
	},
}

func ExportToFile(resp *http.Response) error {

	// Create the file
	out, err := os.Create(exportFileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	return nil
}
