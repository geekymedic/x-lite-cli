/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/geekymedic/x-lite-cli/services"
	"github.com/geekymedic/x-lite-cli/util"
	"github.com/spf13/cobra"
	"os"
)

var (
	systemName string
	dirOut     string
)

// generateSystemCmd represents the generateSystem command
var generateSystemCmd = &cobra.Command{
	Use:   "system",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := services.CreateSystem(systemName, dirOut)
		if err != nil {
			util.StdoutExit(-1, "Fail to generate system: %v", err)
		}
	},
}

func init() {
	generateCmd.AddCommand(generateSystemCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateSystemCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	dir, _ := os.Getwd()
	generateSystemCmd.Flags().StringVar(&systemName, "name", "demo", "system name")
	generateSystemCmd.Flags().StringVar(&dirOut, "out-dir", dir, "system directory")
}
