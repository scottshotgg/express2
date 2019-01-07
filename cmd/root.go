// Copyright © 2018 Scott Gaydos, scgaydos@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// TODO: make a debug logger for every level, or just make our own logger that checks the level
	RootCmd.PersistentFlags().StringP("json-indent", "j", "\t", "output lex tokens in json format")
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "output log messges during operation")
	RootCmd.PersistentFlags().StringP("output", "o", "", "Location to write output to")
	RootCmd.PersistentFlags().BoolP("emit-all", "a", false, "output tokens from all stages and final transpiled C++ with binary")
	RootCmd.PersistentFlags().BoolP("emit-lex", "l", false, "output tokens from lex stage in json format")
	RootCmd.PersistentFlags().BoolP("emit-syn", "y", false, "output tokens from syntactic stage in json format")
	RootCmd.PersistentFlags().BoolP("emit-sem", "s", false, "output tokens from semantic stage in json format")
	RootCmd.PersistentFlags().BoolP("emit-cpp", "c", false, "output transpiled C++ program")

	// RootCmd.PersistentFlags().Visit(func(f *pflag.Flag) {
	// 	var err = viper.BindPFlag(f.Name, RootCmd.PersistentFlags().Lookup(f.Name))
	// 	if err != nil {
	// 		fmt.Printf("err %+v", err)
	// 		os.Exit(9)
	// 	}
	// })

	_ = viper.BindPFlag("json-indent", RootCmd.PersistentFlags().Lookup("json-indent"))
	_ = viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("output", RootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("emit-all", RootCmd.PersistentFlags().Lookup("emit-all"))
	_ = viper.BindPFlag("emit-lex", RootCmd.PersistentFlags().Lookup("emit-lex"))
	_ = viper.BindPFlag("emit-syn", RootCmd.PersistentFlags().Lookup("emit-syn"))
	_ = viper.BindPFlag("emit-sem", RootCmd.PersistentFlags().Lookup("emit-sem"))
	_ = viper.BindPFlag("emit-cpp", RootCmd.PersistentFlags().Lookup("emit-cpp"))

	RootCmd.AddCommand(buildCmd)
	RootCmd.AddCommand(runCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var err = RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "express",
	// TODO: fix this
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("you gotta type `build` or `run` after this")
		fmt.Println("this will print out a help screen soon™")
	},
}
