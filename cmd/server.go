/*
Copyright © 2023 Sam Fatehmanesh sam.fatehmanesh@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the alpha personal assistant server.",
	Long:  `Start the alpha personal assistant server.`,
	Run: func(cmd *cobra.Command, args []string) {

		// launch lisening
		fmt.Println("server called")

		// Define the HTTP endpoint that will receive requests
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Set the content type header to text/plain
			w.Header().Set("Content-Type", "text/plain")

			// Get the text from the request body
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}

			receivedText := string(body)
			outputText, err := initiateMindInstance(receivedText)
			if err != nil {
				log.Fatal(err)
				return
			}

			// To client
			fmt.Fprint(w, outputText)
		})

		// Start the HTTP server and listen for incoming requests
		if err := http.ListenAndServe(":22589", nil); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initiateMindInstance(input string) (string, error) {

	//return input, nil
	return "temporary output", nil
}
