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
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Version is used by the build system.
var Version string

// The log flag value.
var l string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tfe",
	Short: "Manage TFE from the command line.",
	Long:  `Manage TFE from the command line.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true
		if err := setUpLogs(l); err != nil {
			return err
		}
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
	rootCmd.PersistentFlags().StringVarP(&l, "log", "l", "", "log level (debug, info, warn, error, fatal, panic)")

}

// SetUpLogs sets the log level.
func setUpLogs(level string) error {
	// Read the log level
	//  1. from the CLI first
	//  2. then the ENV vars
	//  3. then use the default value.
	if level == "" {
		level = os.Getenv("TFE_LOG_LEVEL")
		if level == "" {
			level = logrus.WarnLevel.String()
		}
	}

	// Parse the log level.
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	// Set the log level.
	logrus.SetLevel(lvl)
	return nil
}
