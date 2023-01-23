/*
Copyright © 2023 Alex <alex8d@pm.me>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ad-8/gobox/io"
	"github.com/ad-8/strava-dl-json/dl"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// dlCmd represents the dl command
var dlCmd = &cobra.Command{
	Use:   "dl",
	Short: "Download and save all activities of a user",
	Run: func(cmd *cobra.Command, args []string) {
		downloadAndSave()
	},
}

// downloadAndSave downloads all Strava activities of a user
// and exports them to a JSON file.
func downloadAndSave() {
	start := time.Now()

	clientID, clientSecret, refreshToken := loadDotEnv()
	tokenInfo, _ := dl.NewTokenInfo(clientID, clientSecret, refreshToken)
	tokenInfo.Print()

	allActivities, err := dl.AllActivities(*tokenInfo)
	if err != nil {
		return
	}
	fmt.Printf("\ndownloaded %d activities in %v\n", len(allActivities), time.Since(start))

	dataJSON, err := json.MarshalIndent(allActivities, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	err = checkDataDirExists()
	if err != nil {
		log.Fatal(err)
	}
	filepath, _ := prefixDataDir(jsonFile)

	if err := io.SimpleWrite(filepath, dataJSON); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("sucessfully written to file %q\n", filepath)
}

func init() {
	rootCmd.AddCommand(dlCmd)
}
