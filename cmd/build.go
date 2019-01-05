// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

	"github.com/scottshotgg/express2/compiler"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	run            bool
	outputFileName string
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use: "build",
	// TODO: fix this
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("ERROR: You must provide an input program")
			return
		}

		// var jsonIndent = viper.GetString("json-indent")

		// // Replace the \t and \n string literals with their non-escaped counterparts
		// jsonIndent = strings.Replace(jsonIndent, `\n`, "\n", -1)
		// jsonIndent = strings.Replace(jsonIndent, `\t`, "\t", -1)

		// TODO: need to check it for all the available characters
		var filenameArg = args[len(args)-1]
		// filenameFull, err := filepath.Abs()
		stat, err := os.Stat(filenameArg)
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}

		if stat.IsDir() {
			fmt.Println("Directory level compilation is not currently supported.")
			os.Exit(0)
		}

		// This is where we get the transpiler name from... so it needs to be passed through?
		// var filename = stat.Name()

		// var outputBase = viper.GetBool("output")
		// if outputBase != "" {
		// 	stat, err := os.Stat(filenameArg)
		// 	if err != nil {
		// 		fmt.Println("ERROR:", err)
		// 		return
		// 	}

		// 	// If its a directory then write all the files with the same name as the executable
		// 	if stat.IsDir() {
		// 		rawFilename = strings.TrimSuffix(filename, ".expr")
		// 	}
		// }

		// If they set it to a directory or there are one or more options
		// enabled then make a directory
		var path = viper.GetString("output")
		// stat, err = os.Stat(path)
		// if err != nil {
		// 	fmt.Println("ERROR:", err)
		// 	return
		// }

		var c = compiler.New(path)

		// Might just make this a config file
		if viper.GetBool("emit-lex") || viper.GetBool("emit-all") {
			c.Output = append(c.Output, "lex")
		}
		if viper.GetBool("emit-compress") || viper.GetBool("emit-all") {
			c.Output = append(c.Output, "compress")
		}

		if viper.GetBool("emit-syn") || viper.GetBool("emit-all") {
			c.Output = append(c.Output, "syn")
		}

		if viper.GetBool("emit-ast") || viper.GetBool("emit-all") {
			c.Output = append(c.Output, "ast")
		}

		if viper.GetBool("emit-flatten") || viper.GetBool("emit-all") {
			c.Output = append(c.Output, "flatten")
		}

		if viper.GetBool("emit-sem") || viper.GetBool("emit-all") {
			c.Output = append(c.Output, "sem")
		}

		if viper.GetBool("emit-cpp") || viper.GetBool("emit-all") {
			c.Output = append(c.Output, "cpp")
		}

		err = c.CompileFile(filenameArg)
		if err != nil {
			fmt.Printf("\nerror: %s\n", err)
			os.Exit(9)
		}

		// if viper.GetBool("emit-cpp") || viper.GetBool("emit-all") {
		// 	output, err = exec.Command("clang-format", "-i", cppFilename).CombinedOutput()
		// 	if err != nil {
		// 		// TODO:
		// 		fmt.Println("err compile", err, string(output))
		// 		os.Exit(9)
		// 		return
		// 	}
		// }
	},
}
