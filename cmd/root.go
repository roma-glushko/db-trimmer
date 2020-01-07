/*
Copyright © 2020 Roman Glushko

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
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

import "db-trimmer/internal/poc"

var cfgFile string

var dbUser string
var dbPassword string
var dbHost string
var dbPort string
var dbName string

var chunkSize int
var plannerThreads int
var trimmerThreads int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "db-trimmer",
	Short: "Reduce database size for your development environments in an intelligent way",
	Long:  `db-trimmer - a tool to reduce database size for your development environments in an intelligent way`,
	Run: func(cmd *cobra.Command, args []string) {
		// nonblockingPoc := poc.NewNonBlockingPoc(
		// 	"mysql",
		// 	fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s`, dbUser, dbPassword, dbHost, dbPort, dbName), // todo: maybe not only tcp connection but also socket
		// 	chunkSize,
		// 	plannerThreads,
		// 	trimmerThreads,
		// )
		// nonblockingPoc.Execute()

		copySchemaPoc := poc.NewCopySchemaPoc(
			"mysql",
			dbUser,
			dbPassword,
			dbHost,
			dbPort,
			dbName,
		)
		copySchemaPoc.Execute()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.db-trimmer.yaml)")

	// db flags
	rootCmd.PersistentFlags().StringVar(&dbUser, "db-user", "root", "Database User")
	rootCmd.PersistentFlags().StringVar(&dbPassword, "db-pass", "", "Database Password")
	rootCmd.MarkPersistentFlagRequired("db-pass")
	rootCmd.PersistentFlags().StringVar(&dbHost, "db-host", "127.0.0.1", "Database Host")
	rootCmd.PersistentFlags().StringVar(&dbPort, "db-port", "3306", "Database Port")
	rootCmd.PersistentFlags().StringVar(&dbName, "db-name", "", "Database Name")
	rootCmd.MarkPersistentFlagRequired("db-name")

	// trimming flags
	rootCmd.PersistentFlags().IntVar(&chunkSize, "chunk-size", 1000, "Chunk Size")
	rootCmd.PersistentFlags().IntVar(&plannerThreads, "planner-threads", 1, "Planner Threads")
	rootCmd.PersistentFlags().IntVar(&trimmerThreads, "trimmer-threads", 1, "Trimmer Threads")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".db-trimmer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".db-trimmer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
