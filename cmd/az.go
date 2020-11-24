/*

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
)

// azCmd represents the az command
var azCmd = &cobra.Command{
	Use:   "az",
	Short: "Enable access to set of commands used to work with Azure cloud",
	Long: `Enable access to set of commands used to work with Azure cloud:
	- authentication - let you access authentication options - e.g. create Service Principal`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("az called")
	},
}

func init() {
	rootCmd.AddCommand(azCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// azCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// azCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}