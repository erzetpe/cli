/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/configuration"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/util"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called") //TODO debug
		config, err := configuration.NewConfig()
		if err != nil {
			panic(fmt.Sprintf("get config failed: %v\n", err)) //TODO err
		}
		name, err := util.PromptForString("Environment name")
		if err != nil {
			panic(fmt.Sprintf("prompt for new environment failed: %v\n", err)) //TODO err
		}
		fmt.Printf("name is: %s\n", name) //TODO debug
		err = config.CreateNewEnvironment(name)
		if err != nil {
			panic(fmt.Sprintf("create new environemtn failed: %v\n", err)) //TODO err
		}
	},
}

func init() {
	environmentsCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
